package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/dlsu-lscs/lscs-core-api/internal/config"
)

// KeyType represents the type of API key being generated
type KeyType int

const (
	KeyTypeDev   KeyType = iota // development key - 30 days default
	KeyTypeProd                 // production key - 365 days default
	KeyTypeAdmin                // admin key - no expiration
)

// default expiration durations
const (
	defaultDevExpiryDays  = 30
	defaultProdExpiryDays = 365
)

// JwtCustomClaims are custom claims extending default ones.
type JwtCustomClaims struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}

// Service is the interface for the auth service.
// It can be used for mocking.
type Service interface {
	GenerateJWT(email string, keyType KeyType) (string, *time.Time, error)
}

type service struct {
	jwtSecret      []byte
	devExpiryDays  int
	prodExpiryDays int
}

// NewService creates a new auth service.
func NewService(secret string, cfg *config.Config) Service {
	devExpiry := defaultDevExpiryDays
	prodExpiry := defaultProdExpiryDays

	// use config values if available
	if cfg != nil {
		if cfg.JWTDevExpiryDays > 0 {
			devExpiry = cfg.JWTDevExpiryDays
		}
		if cfg.JWTProdExpiryDays > 0 {
			prodExpiry = cfg.JWTProdExpiryDays
		}
	}

	return &service{
		jwtSecret:      []byte(secret),
		devExpiryDays:  devExpiry,
		prodExpiryDays: prodExpiry,
	}
}

// GenerateJWT generates a new JWT token with appropriate expiration based on key type.
// returns the token string, expiration time (nil for admin keys), and error if any.
func (s *service) GenerateJWT(email string, keyType KeyType) (string, *time.Time, error) {
	now := time.Now()
	claims := &JwtCustomClaims{
		Email:            email,
		RegisteredClaims: jwt.RegisteredClaims{},
	}

	var expiresAt *time.Time

	switch keyType {
	case KeyTypeDev:
		exp := now.Add(time.Duration(s.devExpiryDays) * 24 * time.Hour)
		expiresAt = &exp
		claims.RegisteredClaims.ExpiresAt = jwt.NewNumericDate(exp)
		claims.RegisteredClaims.IssuedAt = jwt.NewNumericDate(now)
	case KeyTypeProd:
		exp := now.Add(time.Duration(s.prodExpiryDays) * 24 * time.Hour)
		expiresAt = &exp
		claims.RegisteredClaims.ExpiresAt = jwt.NewNumericDate(exp)
		claims.RegisteredClaims.IssuedAt = jwt.NewNumericDate(now)
	case KeyTypeAdmin:
		// admin keys do not expire
		claims.RegisteredClaims.IssuedAt = jwt.NewNumericDate(now)
		expiresAt = nil
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(s.jwtSecret)
	if err != nil {
		return "", nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return tokenString, expiresAt, nil
}
