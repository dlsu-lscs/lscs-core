package helpers

import (
	"context"
	"database/sql"

	"github.com/rs/zerolog/log"

	"github.com/dlsu-lscs/lscs-core-api/internal/database"
	"github.com/dlsu-lscs/lscs-core-api/internal/repository"
)

// AuthorizeIfRNDAndAVP checks if member is:
// 1. In RND committee (any position), OR
// 2. AVP position or higher (any committee)
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

	// RND members are always authorized (any position)
	if authenticatedRequesterInfo.CommitteeID.String == "RND" {
		return true
	}

	// non-RND members need AVP+ position
	if allowedPositions[authenticatedRequesterInfo.PositionID.String] {
		return true
	}

	log.Warn().
		Str("email", email).
		Str("committee", authenticatedRequesterInfo.CommitteeID.String).
		Str("position", authenticatedRequesterInfo.PositionID.String).
		Msg("authorization denied: not RND and not AVP+")
	return false
}
