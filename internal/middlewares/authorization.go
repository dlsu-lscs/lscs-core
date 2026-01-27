package middlewares

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"

	"github.com/dlsu-lscs/lscs-core-api/internal/database"
	"github.com/dlsu-lscs/lscs-core-api/internal/helpers"
)

// AuthorizationMiddleware enforces that the user identified by "user_email" context key
// is an active R&D member with AVP position or higher.
func AuthorizationMiddleware(dbService database.Service) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			email, ok := c.Get("user_email").(string)
			if !ok || email == "" {
				log.Error().Msg("AuthorizationMiddleware: user_email not found in context")
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
			}

			// Perform strict DB check
			isAuthorized := helpers.AuthorizeIfRNDAndAVP(c.Request().Context(), dbService, email)
			if !isAuthorized {
				// The helper already logs the specific reason
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized: Insufficient privileges"})
			}

			return next(c)
		}
	}
}
