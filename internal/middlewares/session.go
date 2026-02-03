package middlewares

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"

	"github.com/dlsu-lscs/lscs-core-api/internal/auth"
	"github.com/dlsu-lscs/lscs-core-api/internal/config"
)

const sessionCookieName = "session_id"

// SessionMiddleware validates session cookies and populates request context with user info.
// It also implements sliding expiration for active sessions.
func SessionMiddleware(sessionService auth.SessionService, cfg *config.Config) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cookie, err := c.Cookie(sessionCookieName)
			if err != nil || cookie.Value == "" {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
			}

			sessionID := cookie.Value
			session, err := sessionService.GetSession(c.Request().Context(), sessionID)
			if err != nil {
				log.Debug().Err(err).Str("session_id", sessionID[:8]+"...").Msg("invalid session")
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Session expired or invalid"})
			}

			// set user info in context
			c.Set("user_id", session.MemberID)
			c.Set("user_email", session.Email)
			c.Set("session_id", session.ID)

			// implement sliding expiration
			duration := time.Duration(cfg.SessionDuration) * time.Second
			if sessionService.ShouldExtendSession(session, duration) {
				if err := sessionService.ExtendSession(c.Request().Context(), sessionID, duration); err != nil {
					log.Warn().Err(err).Str("session_id", sessionID[:8]+"...").Msg("failed to extend session")
				} else {
					log.Debug().Str("session_id", sessionID[:8]+"...").Msg("session extended")
				}
			}

			// update last activity (fire and forget)
			go func() {
				if err := sessionService.UpdateActivity(c.Request().Context(), sessionID); err != nil {
					log.Debug().Err(err).Msg("failed to update session activity")
				}
			}()

			return next(c)
		}
	}
}

// OptionalSessionMiddleware attempts to validate session but allows request to proceed even if no session exists.
// Useful for endpoints that have different behavior for authenticated vs unauthenticated users.
func OptionalSessionMiddleware(sessionService auth.SessionService, cfg *config.Config) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cookie, err := c.Cookie(sessionCookieName)
			if err != nil || cookie.Value == "" {
				// no session, but allow request to proceed
				return next(c)
			}

			sessionID := cookie.Value
			session, err := sessionService.GetSession(c.Request().Context(), sessionID)
			if err != nil {
				// invalid session, but allow request to proceed
				return next(c)
			}

			// set user info in context
			c.Set("user_id", session.MemberID)
			c.Set("user_email", session.Email)
			c.Set("session_id", session.ID)

			// implement sliding expiration
			duration := time.Duration(cfg.SessionDuration) * time.Second
			if sessionService.ShouldExtendSession(session, duration) {
				if err := sessionService.ExtendSession(c.Request().Context(), sessionID, duration); err != nil {
					log.Warn().Err(err).Msg("failed to extend session")
				}
			}

			return next(c)
		}
	}
}
