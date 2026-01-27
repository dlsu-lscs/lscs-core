package committee

import (
	"net/http"

	"github.com/rs/zerolog/log"

	"github.com/dlsu-lscs/lscs-core-api/internal/database"
	"github.com/dlsu-lscs/lscs-core-api/internal/helpers"
	"github.com/dlsu-lscs/lscs-core-api/internal/repository"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	dbService database.Service
}

func NewHandler(dbService database.Service) *Handler {
	return &Handler{
		dbService: dbService,
	}
}

// GetAllCommitteesHandler godoc
// @Summary Get all committees
// @Description Retrieves a list of all committees in the organization
// @Tags committees
// @Accept json
// @Produce json
// @Success 200 {object} GetAllCommitteesResponse "List of committees"
// @Failure 401 {object} helpers.ErrorResponse "Unauthorized"
// @Failure 500 {object} helpers.ErrorResponse "Internal server error"
// @Security BearerAuth
// @Router /committees [get]
func (h *Handler) GetAllCommitteesHandler(c echo.Context) error {
	ctx := c.Request().Context()
	dbconn := h.dbService.GetConnection()
	q := repository.New(dbconn)

	committees, err := q.GetAllCommittees(ctx)
	if err != nil {
		log.Error().Err(err).Msg("failed to get all committees")
		return helpers.ErrInternal(c, "")
	}

	// convert to response format
	response := make([]CommitteeResponse, len(committees))
	for i, comm := range committees {
		response[i] = CommitteeResponse{
			CommitteeID:   comm.CommitteeID,
			CommitteeName: comm.CommitteeName,
		}
		if comm.CommitteeHead.Valid {
			response[i].CommitteeHead = &comm.CommitteeHead.Int32
		}
		if comm.DivisionID.Valid {
			response[i].DivisionID = &comm.DivisionID.String
		}
	}

	return c.JSON(http.StatusOK, GetAllCommitteesResponse{
		Committees: response,
	})
}
