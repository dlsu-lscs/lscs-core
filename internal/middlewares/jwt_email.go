package middlewares

import (
	"github.com/dlsu-lscs/lscs-core-api/internal/auth"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

// JWTEmailMiddleware extracts the email from the JWT claims and sets it in the context.
// This should be used AFTER the echojwt middleware.
func JWTEmailMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		user := c.Get("user").(*jwt.Token)
		claims := user.Claims.(*auth.JwtCustomClaims)
		
		c.Set("user_email", claims.Email)
		
		return next(c)
	}
}
