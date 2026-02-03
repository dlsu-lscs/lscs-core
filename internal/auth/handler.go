package auth

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"

	"github.com/dlsu-lscs/lscs-core-api/internal/database"
	"github.com/dlsu-lscs/lscs-core-api/internal/helpers"
	"github.com/dlsu-lscs/lscs-core-api/internal/repository"
)

// RequestKeyRequest represents the request body for requesting an API key
type RequestKeyRequest struct {
	Project       string `json:"project" validate:"omitempty,max=255" example:"My LSCS Project"`
	AllowedOrigin string `json:"allowed_origin" validate:"omitempty,url" example:"https://example.com"`
	IsDev         bool   `json:"is_dev" example:"false"`
	IsAdmin       bool   `json:"is_admin" example:"false"`
}

// RequestKeyResponse represents the response for a successful API key request
type RequestKeyResponse struct {
	Email     string `json:"email" example:"user@dlsu.edu.ph"`
	APIKey    string `json:"api_key" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	ExpiresAt string `json:"expires_at,omitempty" example:"2027-01-27T15:04:05Z"`
}

type Handler struct {
	authService Service
	dbService   database.Service
	rbacService *RBACService
}

func NewHandler(authService Service, dbService database.Service, rbacService *RBACService) *Handler {
	return &Handler{
		authService: authService,
		dbService:   dbService,
		rbacService: rbacService,
	}
}

// RequestKeyHandler generates a new API key for authorized RND members
// @Summary Request API Key
// @Description Generate a new API key for external projects. Only RND members with AVP position or higher can request keys.
// @Tags auth
// @Accept json
// @Produce json
// @Param request body RequestKeyRequest true "API Key Request"
// @Success 200 {object} RequestKeyResponse "API key generated successfully"
// @Failure 400 {object} helpers.ErrorResponse "Invalid request"
// @Failure 401 {object} helpers.ErrorResponse "Unauthorized"
// @Failure 404 {object} helpers.ErrorResponse "Member not found"
// @Failure 409 {object} helpers.ErrorResponse "Origin already exists"
// @Failure 500 {object} helpers.ErrorResponse "Internal server error"
// @Security GoogleAuth
// @Router /request-key [post]
func (h *Handler) RequestKeyHandler(c echo.Context) error {
	dbconn := h.dbService.GetConnection()
	q := repository.New(dbconn)
	ctx := c.Request().Context()

	// NOTE: this is set via google middleware
	emailRequestor, ok := c.Get("user_email").(string)
	if !ok || emailRequestor == "" {
		log.Error().Msg("user_email not found in context")
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
	}

	isAuthorized := h.rbacService.CanAccessAPIByEmail(c.Request().Context(), emailRequestor)
	if !isAuthorized {
		// forbidden
		// only expose the reason why its unauthorized to the server logs (not on client)
		log.Error().Str("email", emailRequestor).Msg("user has unauthorized position or committee")
		return c.JSON(http.StatusForbidden, map[string]string{"error": "User has unauthorized position or committee"})
	}

	var req RequestKeyRequest
	if err := helpers.BindAndValidate(c, &req); err != nil {
		return err
	}

	memberInfo, err := q.GetMemberInfo(ctx, emailRequestor)
	if err != nil {
		if err == sql.ErrNoRows {
			response := map[string]string{
				"error": "Not an LSCS member",
				"state": "absent",
				"email": emailRequestor,
			}
			return c.JSON(http.StatusNotFound, response)
		}
		log.Error().Err(err).Msg("error checking email")
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
	}

	var allowedOriginForDB sql.NullString
	var keyType KeyType

	if req.IsAdmin {
		allowedOriginForDB = sql.NullString{Valid: false}
		keyType = KeyTypeAdmin
	} else if req.IsDev {
		// TODO: might also want to include the LSCS dev server here instead of just localhost
		if !strings.HasPrefix(req.AllowedOrigin, "http://localhost") {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "For dev keys, allowed_origin must start with http://localhost"})
		}
		allowedOriginForDB = sql.NullString{Valid: false}
		keyType = KeyTypeDev
	} else {
		// production key
		keyType = KeyTypeProd
		if req.AllowedOrigin == "" {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "allowed_origin is required for production keys"})
		}
		_, err := url.ParseRequestURI(req.AllowedOrigin)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid URL for allowed_origin"})
		}
		if strings.HasPrefix(req.AllowedOrigin, "http://localhost") {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "localhost is not a valid origin for production keys"})
		}

		exists, err := q.CheckAllowedOriginExists(ctx, sql.NullString{String: req.AllowedOrigin, Valid: true})
		if err != nil {
			log.Error().Err(err).Msg("failed to check allowed origin")
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error checking origin"})
		}
		if exists {
			return c.JSON(http.StatusConflict, map[string]string{"error": fmt.Sprintf("API key for origin %s already exists", req.AllowedOrigin)})
		}

		allowedOriginForDB = sql.NullString{String: req.AllowedOrigin, Valid: true}
	}

	tokenString, expiresAt, err := h.authService.GenerateJWT(memberInfo.Email, keyType)
	if err != nil {
		log.Error().Err(err).Msg("failed to generate token")
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error generating token"})
	}

	hash := sha256.Sum256([]byte(tokenString))
	hashedTokenString := hex.EncodeToString(hash[:])

	var projectForDB sql.NullString
	if req.Project != "" {
		projectForDB = sql.NullString{String: req.Project, Valid: true}
	} else {
		projectForDB = sql.NullString{Valid: false}
	}

	var expiresAtForDB sql.NullTime
	if expiresAt != nil {
		expiresAtForDB = sql.NullTime{Time: *expiresAt, Valid: true}
	} else {
		expiresAtForDB = sql.NullTime{Valid: false}
	}

	params := repository.StoreAPIKeyParams{
		MemberEmail:   memberInfo.Email,
		ApiKeyHash:    hashedTokenString,
		Project:       projectForDB,
		AllowedOrigin: allowedOriginForDB,
		IsDev:         req.IsDev,
		IsAdmin:       req.IsAdmin,
		ExpiresAt:     expiresAtForDB,
	}

	err = q.StoreAPIKey(ctx, params)
	if err != nil {
		log.Error().Err(err).Msg("failed to store api key")
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error storing API key"})
	}

	response := map[string]interface{}{
		"email":   memberInfo.Email,
		"api_key": tokenString,
	}

	// include expiration time in response if applicable
	if expiresAt != nil {
		response["expires_at"] = expiresAt.Format("2006-01-02T15:04:05Z07:00")
	}

	return c.JSON(http.StatusOK, response)
}

// ListAPIKeysResponse represents an API key in the list response
type ListAPIKeysResponse struct {
	APIKeyID      int32  `json:"api_key_id"`
	MemberEmail   string `json:"member_email"`
	Project       string `json:"project,omitempty"`
	AllowedOrigin string `json:"allowed_origin,omitempty"`
	IsDev         bool   `json:"is_dev"`
	IsAdmin       bool   `json:"is_admin"`
	CreatedAt     string `json:"created_at"`
	ExpiresAt     string `json:"expires_at,omitempty"`
}

// ListAPIKeys returns all API keys for the authenticated user
// @Summary List API Keys
// @Description Get all API keys belonging to the authenticated RND member
// @Tags auth
// @Accept json
// @Produce json
// @Success 200 {array} ListAPIKeysResponse "List of API keys"
// @Failure 401 {object} helpers.ErrorResponse "Unauthorized"
// @Failure 403 {object} helpers.ErrorResponse "Forbidden - RND committee only"
// @Failure 500 {object} helpers.ErrorResponse "Internal server error"
// @Security SessionAuth
// @Router /api-keys [get]
func (h *Handler) ListAPIKeys(c echo.Context) error {
	ctx := c.Request().Context()
	dbconn := h.dbService.GetConnection()
	q := repository.New(dbconn)

	// Get user email from context (set by session middleware)
	email, ok := c.Get("user_email").(string)
	if !ok || email == "" {
		return helpers.ErrUnauthorized(c, "")
	}

	// Check RND AVP+ authorization
	isAuthorized := h.rbacService.CanAccessAPIByEmail(ctx, email)
	if !isAuthorized {
		log.Error().Str("email", email).Msg("user has unauthorized position or committee for API key management")
		return helpers.ErrForbidden(c, "RND AVP+ position required for API key management")
	}

	apiKeys, err := q.ListAPIKeysByEmail(ctx, email)
	if err != nil {
		log.Error().Err(err).Str("email", email).Msg("failed to list API keys")
		return helpers.ErrInternal(c, "Failed to retrieve API keys")
	}

	// Transform to response format
	response := make([]ListAPIKeysResponse, len(apiKeys))
	for i, key := range apiKeys {
		resp := ListAPIKeysResponse{
			APIKeyID:    key.ApiKeyID,
			MemberEmail: key.MemberEmail,
			IsDev:       key.IsDev,
			IsAdmin:     key.IsAdmin,
		}

		if key.CreatedAt.Valid {
			resp.CreatedAt = key.CreatedAt.Time.Format("2006-01-02T15:04:05Z07:00")
		}
		if key.Project.Valid {
			resp.Project = key.Project.String
		}
		if key.AllowedOrigin.Valid {
			resp.AllowedOrigin = key.AllowedOrigin.String
		}
		if key.ExpiresAt.Valid {
			resp.ExpiresAt = key.ExpiresAt.Time.Format("2006-01-02T15:04:05Z07:00")
		}

		response[i] = resp
	}

	return c.JSON(http.StatusOK, response)
}

// RevokeAPIKey deletes an API key by ID
// @Summary Revoke API Key
// @Description Delete/revoke an API key by ID (must belong to authenticated user)
// @Tags auth
// @Accept json
// @Produce json
// @Param id path int true "API Key ID"
// @Success 200 {object} map[string]string "Success message"
// @Failure 401 {object} helpers.ErrorResponse "Unauthorized"
// @Failure 403 {object} helpers.ErrorResponse "Forbidden - not owner or RND"
// @Failure 404 {object} helpers.ErrorResponse "Not found"
// @Failure 500 {object} helpers.ErrorResponse "Internal server error"
// @Security SessionAuth
// @Router /api-keys/{id} [delete]
func (h *Handler) RevokeAPIKey(c echo.Context) error {
	ctx := c.Request().Context()
	dbconn := h.dbService.GetConnection()
	q := repository.New(dbconn)

	// Get user email from context
	email, ok := c.Get("user_email").(string)
	if !ok || email == "" {
		return helpers.ErrUnauthorized(c, "")
	}

	// Check RND AVP+ authorization
	isAuthorized := h.rbacService.CanAccessAPIByEmail(ctx, email)
	if !isAuthorized {
		log.Error().Str("email", email).Msg("user has unauthorized position or committee for API key management")
		return helpers.ErrForbidden(c, "RND AVP+ position required for API key management")
	}

	// Parse API key ID from URL
	apiKeyID, err := strconv.ParseInt(c.Param("id"), 10, 32)
	if err != nil {
		return helpers.ErrBadRequest(c, "Invalid API key ID")
	}

	// Delete the key (will only succeed if member_email matches)
	err = q.DeleteAPIKeyById(ctx, repository.DeleteAPIKeyByIdParams{
		ApiKeyID:    int32(apiKeyID),
		MemberEmail: email,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			return helpers.ErrNotFound(c, "API key not found or you don't have permission to revoke it")
		}
		log.Error().Err(err).Int64("api_key_id", apiKeyID).Str("email", email).Msg("failed to revoke API key")
		return helpers.ErrInternal(c, "Failed to revoke API key")
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "API key revoked successfully",
	})
}
