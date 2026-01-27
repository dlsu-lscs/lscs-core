package middlewares

import (
	"context"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"google.golang.org/api/idtoken"

	"github.com/dlsu-lscs/lscs-core-api/internal/config"
)

// GoogleAuthMiddleware validates Google ID tokens and extracts user email
func GoogleAuthMiddleware(cfg *config.Config) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Authorization header is required"})
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid Authorization header format"})
			}

			tokenString := parts[1]
			audience := cfg.GoogleClientID
			if audience == "" {
				log.Error().Msg("GOOGLE_CLIENT_ID not configured")
				return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
			}

			payload, err := idtoken.Validate(context.Background(), tokenString, audience)
			if err != nil {
				log.Error().Err(err).Msg("failed to validate google token")
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid ID token"})
			}

			email, ok := payload.Claims["email"].(string)
			if !ok || email == "" {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Email not found in token"})
			}

			c.Set("user_email", email)
			return next(c)
		}
	}
}
