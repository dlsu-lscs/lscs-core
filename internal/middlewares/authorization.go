package middlewares

import (
	"log/slog"
	"net/http"

	"github.com/dlsu-lscs/lscs-core-api/internal/database"
	"github.com/dlsu-lscs/lscs-core-api/internal/helpers"
	"github.com/labstack/echo/v4"
)

// AuthorizationMiddleware enforces that the user identified by "user_email" context key
// is an active R&D member with AVP position or higher.
func AuthorizationMiddleware(dbService database.Service) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			email, ok := c.Get("user_email").(string)
			if !ok || email == "" {
				slog.Error("AuthorizationMiddleware: user_email not found in context")
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
