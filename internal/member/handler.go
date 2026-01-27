package member

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"

	"github.com/dlsu-lscs/lscs-core-api/internal/helpers"
	"github.com/dlsu-lscs/lscs-core-api/internal/repository"
)

// GetMemberInfo retrieves detailed member information by email
// @Summary Get member info by email
// @Description Get complete member information using their email address
// @Tags members
// @Accept json
// @Produce json
// @Param request body EmailRequest true "Email Request"
// @Success 200 {object} FullInfoMemberResponse "Member information"
// @Failure 400 {object} helpers.ErrorResponse "Invalid request"
// @Failure 404 {object} helpers.ErrorResponse "Member not found"
// @Failure 500 {object} helpers.ErrorResponse "Internal server error"
// @Security BearerAuth
// @Router /member [post]
func (h *Handler) GetMemberInfo(c echo.Context) error {
	ctx := c.Request().Context()
	dbconn := h.dbService.GetConnection()
	q := repository.New(dbconn)

	req := new(EmailRequest)

	if err := helpers.BindAndValidate(c, req); err != nil {
		return err
	}

	memberInfo, err := q.GetMemberInfo(ctx, req.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			response := map[string]string{
				"error": "Not an LSCS member",
				"state": "absent",
				"email": req.Email,
			}
			return c.JSON(http.StatusNotFound, response)
		}
		log.Error().Err(err).Msg("error checking email")
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
	}

	response := toFullInfoMemberResponse(memberInfo)

	return c.JSON(http.StatusOK, response)
}

// GetMemberInfoByID retrieves detailed member information by ID
// @Summary Get member info by ID
// @Description Get complete member information using their student ID
// @Tags members
// @Accept json
// @Produce json
// @Param request body IdRequest true "ID Request"
// @Success 200 {object} FullInfoMemberResponse "Member information"
// @Failure 400 {object} helpers.ErrorResponse "Invalid request"
// @Failure 404 {object} helpers.ErrorResponse "Member not found"
// @Failure 500 {object} helpers.ErrorResponse "Internal server error"
// @Security BearerAuth
// @Router /member-id [post]
func (h *Handler) GetMemberInfoByID(c echo.Context) error {
	ctx := c.Request().Context()
	dbconn := h.dbService.GetConnection()
	q := repository.New(dbconn)

	req := new(IdRequest)

	if err := helpers.BindAndValidate(c, req); err != nil {
		return err
	}

	memberInfo, err := q.GetMemberInfoById(ctx, int32(req.Id))
	if err != nil {
		if err == sql.ErrNoRows {
			log.Error().Err(err).Int("id", req.Id).Msg("id is not an LSCS member")
			response := map[string]string{
				"error": "Not an LSCS member",
				"state": "absent",
				"id":    strconv.Itoa(req.Id),
			}
			return c.JSON(http.StatusNotFound, response)
		}
		log.Error().Err(err).Msg("error checking id")
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
	}

	response := toFullInfoMemberResponse(repository.GetMemberInfoRow(memberInfo))

	return c.JSON(http.StatusOK, response)
}

// GetAllMembersHandler lists all members
// @Summary List all members
// @Description Get a list of all LSCS members with basic information
// @Tags members
// @Produce json
// @Success 200 {array} MemberResponse "List of members"
// @Failure 500 {object} helpers.ErrorResponse "Internal server error"
// @Security BearerAuth
// @Router /members [get]
func (h *Handler) GetAllMembersHandler(c echo.Context) error {
	ctx := c.Request().Context()
	dbconn := h.dbService.GetConnection()
	queries := repository.New(dbconn)

	members, err := queries.ListMembers(ctx)
	if err != nil {
		log.Error().Err(err).Msg("failed to list members")
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to list members"})
	}

	response := make([]MemberResponse, 0, len(members))
	for _, m := range members {
		response = append(response, toMemberResponse(m))
	}

	return c.JSON(http.StatusOK, response)
}

// CheckEmailHandler checks if an email belongs to an LSCS member
// @Summary Check email membership
// @Description Verify if an email address belongs to an LSCS member
// @Tags members
// @Accept json
// @Produce json
// @Param request body EmailRequest true "Email Request"
// @Success 200 {object} map[string]interface{} "Member exists"
// @Failure 400 {object} helpers.ErrorResponse "Invalid request"
// @Failure 404 {object} helpers.ErrorResponse "Member not found"
// @Failure 500 {object} helpers.ErrorResponse "Internal server error"
// @Security BearerAuth
// @Router /check-email [post]
func (h *Handler) CheckEmailHandler(c echo.Context) error {
	var req EmailRequest

	if err := helpers.BindAndValidate(c, &req); err != nil {
		return err
	}

	ctx := c.Request().Context()
	dbconn := h.dbService.GetConnection()
	queries := repository.New(dbconn)
	memberEmail, err := queries.CheckEmailIfMember(ctx, req.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			response := map[string]string{
				"error": "Not an LSCS member",
				"state": "absent",
				"email": req.Email,
			}
			return c.JSON(http.StatusNotFound, response)
		}
		log.Error().Err(err).Msg("error checking email")
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
	}

	response := map[string]any{
		"success": "Email is an LSCS member",
		"state":   "present",
		"email":   memberEmail,
	}
	return c.JSON(http.StatusOK, response)
}

// CheckIDIfMember checks if an ID belongs to an LSCS member
// @Summary Check ID membership
// @Description Verify if a student ID belongs to an LSCS member
// @Tags members
// @Accept json
// @Produce json
// @Param request body IdRequest true "ID Request"
// @Success 200 {object} map[string]interface{} "Member exists"
// @Failure 400 {object} helpers.ErrorResponse "Invalid request"
// @Failure 404 {object} helpers.ErrorResponse "Member not found"
// @Failure 500 {object} helpers.ErrorResponse "Internal server error"
// @Security BearerAuth
// @Router /check-id [post]
func (h *Handler) CheckIDIfMember(c echo.Context) error {
	var req IdRequest

	if err := helpers.BindAndValidate(c, &req); err != nil {
		return err
	}

	dbconn := h.dbService.GetConnection()
	q := repository.New(dbconn)
	id, err := q.CheckIdIfMember(c.Request().Context(), int32(req.Id))
	if err != nil {
		if err == sql.ErrNoRows {
			response := map[string]any{
				"error": "Not an LSCS member",
				"state": "absent",
				"id":    req.Id,
			}
			return c.JSON(http.StatusNotFound, response)
		}
		log.Error().Err(err).Msg("invalid ID")
		return c.JSON(http.StatusNotFound, map[string]string{"error": "invalid ID"})
	}

	response := map[string]any{
		"success": "ID is an LSCS member",
		"state":   "present",
		"id":      id,
	}
	return c.JSON(http.StatusOK, response)
}
