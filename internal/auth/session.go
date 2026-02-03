package auth

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/dlsu-lscs/lscs-core-api/internal/config"
	"github.com/dlsu-lscs/lscs-core-api/internal/repository"
)

// SessionService handles session management for web UI authentication
type SessionService interface {
	CreateSession(ctx context.Context, memberID int32, rememberMe bool, userAgent, ipAddress string) (*Session, error)
	GetSession(ctx context.Context, sessionID string) (*SessionWithMember, error)
	UpdateActivity(ctx context.Context, sessionID string) error
	ExtendSession(ctx context.Context, sessionID string, duration time.Duration) error
	DeleteSession(ctx context.Context, sessionID string) error
	DeleteAllSessionsForMember(ctx context.Context, memberID int32) error
	CleanupExpiredSessions(ctx context.Context) error
	ShouldExtendSession(session *SessionWithMember, duration time.Duration) bool
}

// Session represents a user session
type Session struct {
	ID           string
	MemberID     int32
	CreatedAt    time.Time
	ExpiresAt    time.Time
	LastActivity time.Time
	UserAgent    string
	IPAddress    string
}

// SessionWithMember includes member info for context population
type SessionWithMember struct {
	Session
	Email    string
	FullName string
}

type sessionService struct {
	db                     *sql.DB
	defaultDuration        time.Duration
	rememberMeDuration     time.Duration
	slidingExtendThreshold float64 // extend if remaining time is less than this fraction of duration
}

// NewSessionService creates a new session service
func NewSessionService(db *sql.DB, cfg *config.Config) SessionService {
	defaultDuration := 24 * time.Hour
	rememberMeDuration := 30 * 24 * time.Hour

	if cfg != nil {
		if cfg.SessionDuration > 0 {
			defaultDuration = time.Duration(cfg.SessionDuration) * time.Second
		}
		if cfg.SessionRememberDuration > 0 {
			rememberMeDuration = time.Duration(cfg.SessionRememberDuration) * time.Second
		}
	}

	return &sessionService{
		db:                     db,
		defaultDuration:        defaultDuration,
		rememberMeDuration:     rememberMeDuration,
		slidingExtendThreshold: 0.5, // extend when less than 50% time remaining
	}
}

// generateSessionID creates a cryptographically secure session ID
func generateSessionID() (string, error) {
	bytes := make([]byte, 32) // 256 bits
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// CreateSession creates a new session for a member
func (s *sessionService) CreateSession(ctx context.Context, memberID int32, rememberMe bool, userAgent, ipAddress string) (*Session, error) {
	sessionID, err := generateSessionID()
	if err != nil {
		return nil, err
	}

	duration := s.defaultDuration
	if rememberMe {
		duration = s.rememberMeDuration
	}

	expiresAt := time.Now().Add(duration)

	q := repository.New(s.db)
	err = q.CreateSession(ctx, repository.CreateSessionParams{
		ID:        sessionID,
		MemberID:  memberID,
		ExpiresAt: expiresAt,
		UserAgent: sql.NullString{String: userAgent, Valid: userAgent != ""},
		IpAddress: sql.NullString{String: ipAddress, Valid: ipAddress != ""},
	})
	if err != nil {
		return nil, err
	}

	return &Session{
		ID:           sessionID,
		MemberID:     memberID,
		CreatedAt:    time.Now(),
		ExpiresAt:    expiresAt,
		LastActivity: time.Now(),
		UserAgent:    userAgent,
		IPAddress:    ipAddress,
	}, nil
}

// GetSession retrieves a session by ID (only if not expired)
func (s *sessionService) GetSession(ctx context.Context, sessionID string) (*SessionWithMember, error) {
	q := repository.New(s.db)
	row, err := q.GetSessionWithMember(ctx, sessionID)
	if err != nil {
		return nil, err
	}

	session := &SessionWithMember{
		Session: Session{
			ID:        row.ID,
			MemberID:  row.MemberID,
			ExpiresAt: row.ExpiresAt,
		},
		Email:    row.Email,
		FullName: row.FullName,
	}

	if row.CreatedAt.Valid {
		session.CreatedAt = row.CreatedAt.Time
	}
	if row.LastActivity.Valid {
		session.LastActivity = row.LastActivity.Time
	}
	if row.UserAgent.Valid {
		session.UserAgent = row.UserAgent.String
	}
	if row.IpAddress.Valid {
		session.IPAddress = row.IpAddress.String
	}

	return session, nil
}

// UpdateActivity updates the last activity timestamp
func (s *sessionService) UpdateActivity(ctx context.Context, sessionID string) error {
	q := repository.New(s.db)
	return q.UpdateSessionActivity(ctx, sessionID)
}

// ExtendSession extends the session expiration
func (s *sessionService) ExtendSession(ctx context.Context, sessionID string, duration time.Duration) error {
	q := repository.New(s.db)
	newExpiry := time.Now().Add(duration)
	return q.ExtendSession(ctx, repository.ExtendSessionParams{
		ExpiresAt: newExpiry,
		ID:        sessionID,
	})
}

// DeleteSession removes a session
func (s *sessionService) DeleteSession(ctx context.Context, sessionID string) error {
	q := repository.New(s.db)
	return q.DeleteSession(ctx, sessionID)
}

// DeleteAllSessionsForMember removes all sessions for a member (logout everywhere)
func (s *sessionService) DeleteAllSessionsForMember(ctx context.Context, memberID int32) error {
	q := repository.New(s.db)
	return q.DeleteAllSessionsForMember(ctx, memberID)
}

// CleanupExpiredSessions removes all expired sessions from the database
func (s *sessionService) CleanupExpiredSessions(ctx context.Context) error {
	q := repository.New(s.db)
	return q.CleanupExpiredSessions(ctx)
}

// ShouldExtendSession determines if a session should be extended based on remaining time
func (s *sessionService) ShouldExtendSession(session *SessionWithMember, duration time.Duration) bool {
	remaining := time.Until(session.ExpiresAt)
	threshold := time.Duration(float64(duration) * s.slidingExtendThreshold)
	return remaining < threshold
}

// GetDefaultDuration returns the default session duration
func (s *sessionService) GetDefaultDuration() time.Duration {
	return s.defaultDuration
}

// GetRememberMeDuration returns the remember me session duration
func (s *sessionService) GetRememberMeDuration() time.Duration {
	return s.rememberMeDuration
}

// StartCleanupJob starts a background goroutine that periodically cleans up expired sessions
func StartCleanupJob(ctx context.Context, service SessionService, interval time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				log.Info().Msg("session cleanup job stopped")
				return
			case <-ticker.C:
				if err := service.CleanupExpiredSessions(ctx); err != nil {
					log.Error().Err(err).Msg("failed to cleanup expired sessions")
				} else {
					log.Debug().Msg("expired sessions cleanup completed")
				}
			}
		}
	}()
}
