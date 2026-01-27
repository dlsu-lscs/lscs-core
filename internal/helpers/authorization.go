package helpers

import (
	"context"
	"database/sql"

	"github.com/rs/zerolog/log"

	"github.com/dlsu-lscs/lscs-core-api/internal/database"
	"github.com/dlsu-lscs/lscs-core-api/internal/repository"
)

// AuthorizeIfRNDAndAVP assumes a member, then checks if member is in RND and AVP or higher
func AuthorizeIfRNDAndAVP(ctx context.Context, dbService database.Service, email string) bool {
	dbconn := dbService.GetConnection()
	q := repository.New(dbconn)

	allowedPositions := map[string]bool{
		"PRES": true,
		"EVP":  true,
		"VP":   true,
		"AVP":  true,
		"CT":   false,
		"JO":   false,
		"MEM":  false,
	}

	authenticatedRequesterInfo, err := q.GetMemberInfo(ctx, email)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Error().Str("email", email).Err(err).Msg("not an LSCS member")
			return false
		}
		log.Error().Str("email", email).Err(err).Msg("error checking email")
		return false
	}

	// only allow RND members and those that are AVPs and above
	if authenticatedRequesterInfo.CommitteeID.String != "RND" {
		log.Error().Str("email", email).Msg("not a member of Research and Development")
		if !allowedPositions[authenticatedRequesterInfo.PositionID.String] {
			log.Error().Str("email", email).Msg("not AVP or higher")
			return false
		}
		return false
	}

	return true
}
