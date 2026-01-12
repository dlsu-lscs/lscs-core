package auth

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"strings"

	"github.com/dlsu-lscs/lscs-core-api/internal/database"
	"github.com/dlsu-lscs/lscs-core-api/internal/helpers"
	"github.com/dlsu-lscs/lscs-core-api/internal/repository"
	"github.com/labstack/echo/v4"
)

type RequestKeyRequest struct {
	Email         string `json:"email" validate:"required,email"`
	Project       string `json:"project"`
	AllowedOrigin string `json:"allowed_origin"`
	IsDev         bool   `json:"is_dev"`
	IsAdmin       bool   `json:"is_admin"`
}

type Handler struct {
	authService Service
	dbService   database.Service
}

func NewHandler(authService Service, dbService database.Service) *Handler {
	return &Handler{
		authService: authService,
		dbService:   dbService,
	}
}

func (h *Handler) RequestKeyHandler(c echo.Context) error {
	dbconn := h.dbService.GetConnection()
	q := repository.New(dbconn)
	ctx := c.Request().Context()

	// NOTE: this is set via google middleware
	emailRequestor := c.Get("user_email").(string)
	isAuthorized := helpers.AuthorizeIfRNDAndAVP(c.Request().Context(), h.dbService, emailRequestor)
	if !isAuthorized {
		// forbidden
		// only expose the reason why its unauthorized to the server logs (not on client)
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized request."})
	}

	var req RequestKeyRequest
	if err := c.Bind(&req); err != nil {
		slog.Error("cannot read body", "error", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "cannot read body"})
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
		slog.Error("error checking email", "error", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
	}

	var allowedOriginForDB sql.NullString
	var isDevForDB bool

	if req.IsAdmin {
		allowedOriginForDB = sql.NullString{Valid: false}
		isDevForDB = false
	} else if req.IsDev {
		// TODO: might also want to include the LSCS dev server here instead of just localhost
		if !strings.HasPrefix(req.AllowedOrigin, "http://localhost") {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "For dev keys, allowed_origin must start with http://localhost"})
		}
		allowedOriginForDB = sql.NullString{Valid: false}
		isDevForDB = true
	} else {
		// Production key
		// TODO: if "is_dev: true" (have expiry time for API_KEY token)
		// TODO: only "is_admin: true" API_KEY do not expire
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
			slog.Error("failed to check allowed origin", "error", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error checking origin"})
		}
		if exists {
			return c.JSON(http.StatusConflict, map[string]string{"error": fmt.Sprintf("API key for origin %s already exists", req.AllowedOrigin)})
		}

		allowedOriginForDB = sql.NullString{String: req.AllowedOrigin, Valid: true}
		isDevForDB = false
	}

	tokenString, err := h.authService.GenerateJWT(memberInfo.Email)
	if err != nil {
		slog.Error("failed to generate token", "error", err)
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

	params := repository.StoreAPIKeyParams{
		MemberEmail:   memberInfo.Email,
		ApiKeyHash:    hashedTokenString,
		Project:       projectForDB,
		AllowedOrigin: allowedOriginForDB,
		IsDev:         isDevForDB,
		IsAdmin:       req.IsAdmin,
		ExpiresAt:     sql.NullTime{Valid: false},
	}

	err = q.StoreAPIKey(ctx, params)
	if err != nil {
		slog.Error("failed to store api key", "error", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error storing API key"})
	}

	response := map[string]interface{}{
		"email":   memberInfo.Email,
		"api_key": tokenString,
	}

	return c.JSON(http.StatusOK, response)
}
