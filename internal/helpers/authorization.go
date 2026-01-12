package helpers

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"

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
			slog.Error(fmt.Sprintf("not an LSCS member: %s", email), "error", err)
			return false
		}
		slog.Error(fmt.Sprintf("error checking email: %s", email), "error", err)
		return false
	}

	// only allow RND members and those that are AVPs and above
	if authenticatedRequesterInfo.CommitteeID.String != "RND" {
		slog.Error(fmt.Sprintf("Email: %s is not a member of Research and Development.", email))
		if !allowedPositions[authenticatedRequesterInfo.PositionID.String] {
			slog.Error(fmt.Sprintf("Email: %s is not AVP or higher.", email))
			return false
		}
		return false
	}

	return true
}
