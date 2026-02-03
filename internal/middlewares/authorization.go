package middlewares

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"

	"github.com/dlsu-lscs/lscs-core-api/internal/auth"
	"github.com/dlsu-lscs/lscs-core-api/internal/database"
	"github.com/dlsu-lscs/lscs-core-api/internal/helpers"
)

// AuthorizationMiddleware enforces that the user identified by "user_email" context key
// is an active R&D member with AVP position or higher.
// Deprecated: Use RBACMiddleware for more granular control
func AuthorizationMiddleware(dbService database.Service) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			email, ok := c.Get("user_email").(string)
			if !ok || email == "" {
				log.Error().Msg("AuthorizationMiddleware: user_email not found in context")
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
			}

			// perform strict DB check
			isAuthorized := helpers.AuthorizeIfRNDAndAVP(c.Request().Context(), dbService, email)
			if !isAuthorized {
				// the helper already logs the specific reason
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized: Insufficient privileges"})
			}

			return next(c)
		}
	}
}

// RequireAdmin middleware ensures the user has the ADMIN role
func RequireAdmin(rbacService *auth.RBACService) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			memberID, ok := c.Get("user_id").(int32)
			if !ok {
				log.Error().Msg("RequireAdmin: user_id not found in context")
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
			}

			if !rbacService.IsAdmin(c.Request().Context(), memberID) {
				log.Warn().Int32("member_id", memberID).Msg("admin access denied")
				return c.JSON(http.StatusForbidden, map[string]string{"error": "Admin access required"})
			}

			return next(c)
		}
	}
}

// RequireRole middleware ensures the user has a specific role
func RequireRole(rbacService *auth.RBACService, roleID string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			memberID, ok := c.Get("user_id").(int32)
			if !ok {
				log.Error().Msg("RequireRole: user_id not found in context")
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
			}

			if !rbacService.HasRole(c.Request().Context(), memberID, roleID) {
				log.Warn().Int32("member_id", memberID).Str("role", roleID).Msg("role access denied")
				return c.JSON(http.StatusForbidden, map[string]string{"error": "Insufficient permissions"})
			}

			return next(c)
		}
	}
}

// RequirePosition middleware ensures the user has a minimum position level
func RequirePosition(dbService database.Service, minPosition string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			memberID, ok := c.Get("user_id").(int32)
			if !ok {
				log.Error().Msg("RequirePosition: user_id not found in context")
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
			}

			// get member position from db
			member, err := getMemberByID(c, dbService, memberID)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
			}

			if !auth.IsHigherOrEqualPosition(member.PositionID.String, minPosition) {
				log.Warn().
					Int32("member_id", memberID).
					Str("position", member.PositionID.String).
					Str("required", minPosition).
					Msg("position access denied")
				return c.JSON(http.StatusForbidden, map[string]string{"error": "Insufficient position level"})
			}

			return next(c)
		}
	}
}

// RequireAPIKeyAccess middleware ensures the user can access API key management
func RequireAPIKeyAccess(rbacService *auth.RBACService) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			memberID, ok := c.Get("user_id").(int32)
			if !ok {
				log.Error().Msg("RequireAPIKeyAccess: user_id not found in context")
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
			}

			if !rbacService.CanAccessAPIKeyManagement(c.Request().Context(), memberID) {
				log.Warn().Int32("member_id", memberID).Msg("API key access denied")
				return c.JSON(http.StatusForbidden, map[string]string{"error": "API key management access required"})
			}

			return next(c)
		}
	}
}

// RequireCanEditMember middleware ensures the user can edit the target member
// Expects the target member ID to be in the URL path parameter "id"
func RequireCanEditMember(rbacService *auth.RBACService) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			actorID, ok := c.Get("user_id").(int32)
			if !ok {
				log.Error().Msg("RequireCanEditMember: user_id not found in context")
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
			}

			targetIDStr := c.Param("id")
			targetID, err := strconv.ParseInt(targetIDStr, 10, 32)
			if err != nil {
				return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid member ID"})
			}

			if !rbacService.CanEditMember(c.Request().Context(), actorID, int32(targetID)) {
				log.Warn().
					Int32("actor_id", actorID).
					Int64("target_id", targetID).
					Msg("edit member access denied")
				return c.JSON(http.StatusForbidden, map[string]string{"error": "Cannot edit this member"})
			}

			// store target ID for handler use
			c.Set("target_member_id", int32(targetID))

			return next(c)
		}
	}
}

// RequireAdminOrSelf middleware allows access if user is admin OR is accessing their own resource
// Expects the target member ID to be in the URL path parameter "id"
func RequireAdminOrSelf(rbacService *auth.RBACService) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			actorID, ok := c.Get("user_id").(int32)
			if !ok {
				log.Error().Msg("RequireAdminOrSelf: user_id not found in context")
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
			}

			targetIDStr := c.Param("id")
			targetID, err := strconv.ParseInt(targetIDStr, 10, 32)
			if err != nil {
				return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid member ID"})
			}

			// allow if same user
			if actorID == int32(targetID) {
				c.Set("target_member_id", int32(targetID))
				return next(c)
			}

			// allow if admin
			if rbacService.IsAdmin(c.Request().Context(), actorID) {
				c.Set("target_member_id", int32(targetID))
				return next(c)
			}

			log.Warn().
				Int32("actor_id", actorID).
				Int64("target_id", targetID).
				Msg("admin or self access denied")
			return c.JSON(http.StatusForbidden, map[string]string{"error": "Access denied"})
		}
	}
}

// helper to get member by ID using repository
func getMemberByID(c echo.Context, dbService database.Service, memberID int32) (*memberInfo, error) {
	ctx := c.Request().Context()
	db := dbService.GetConnection()

	var member memberInfo
	err := db.QueryRowContext(ctx,
		`SELECT id, position_id, committee_id FROM members WHERE id = ?`,
		memberID,
	).Scan(&member.ID, &member.PositionID, &member.CommitteeID)

	if err != nil {
		log.Error().Err(err).Int32("member_id", memberID).Msg("failed to get member for authorization")
		return nil, err
	}

	return &member, nil
}

// memberInfo is a minimal struct for authorization checks
type memberInfo struct {
	ID          int32
	PositionID  nullableString
	CommitteeID nullableString
}

// nullableString implements sql.Scanner for nullable strings
type nullableString struct {
	String string
	Valid  bool
}

func (ns *nullableString) Scan(value interface{}) error {
	if value == nil {
		ns.String, ns.Valid = "", false
		return nil
	}
	switch v := value.(type) {
	case string:
		ns.String, ns.Valid = v, true
	case []byte:
		ns.String, ns.Valid = string(v), true
	}
	return nil
}
