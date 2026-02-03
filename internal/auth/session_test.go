package auth

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/dlsu-lscs/lscs-core-api/internal/config"
)

func TestSessionService_CreateSession(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	cfg := &config.Config{
		SessionDuration:         86400,   // 24h
		SessionRememberDuration: 2592000, // 30d
	}
	service := NewSessionService(db, cfg)

	t.Run("create session without remember me", func(t *testing.T) {
		mock.ExpectExec("INSERT INTO sessions").
			WithArgs(sqlmock.AnyArg(), int32(123), sqlmock.AnyArg(), "Mozilla/5.0", "192.168.1.1").
			WillReturnResult(sqlmock.NewResult(1, 1))

		session, err := service.CreateSession(context.Background(), 123, false, "Mozilla/5.0", "192.168.1.1")

		assert.NoError(t, err)
		assert.NotEmpty(t, session.ID)
		assert.Equal(t, int32(123), session.MemberID)
		assert.Equal(t, "Mozilla/5.0", session.UserAgent)
		assert.Equal(t, "192.168.1.1", session.IPAddress)
		// expiry should be around 24h from now
		assert.WithinDuration(t, time.Now().Add(24*time.Hour), session.ExpiresAt, 5*time.Second)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("create session with remember me", func(t *testing.T) {
		mock.ExpectExec("INSERT INTO sessions").
			WithArgs(sqlmock.AnyArg(), int32(456), sqlmock.AnyArg(), "Chrome/100", "10.0.0.1").
			WillReturnResult(sqlmock.NewResult(1, 1))

		session, err := service.CreateSession(context.Background(), 456, true, "Chrome/100", "10.0.0.1")

		assert.NoError(t, err)
		assert.NotEmpty(t, session.ID)
		assert.Equal(t, int32(456), session.MemberID)
		// expiry should be around 30d from now
		assert.WithinDuration(t, time.Now().Add(30*24*time.Hour), session.ExpiresAt, 5*time.Second)

		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestSessionService_GetSession(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	service := NewSessionService(db, nil)

	t.Run("get valid session", func(t *testing.T) {
		sessionID := "abc123def456"
		now := time.Now()
		expiresAt := now.Add(12 * time.Hour)

		rows := sqlmock.NewRows([]string{
			"id", "member_id", "created_at", "expires_at", "last_activity", "user_agent", "ip_address", "email", "full_name",
		}).AddRow(
			sessionID, int32(123), now, expiresAt, now, "Mozilla/5.0", "192.168.1.1", "test@dlsu.edu.ph", "Test User",
		)

		mock.ExpectQuery("SELECT").
			WithArgs(sessionID).
			WillReturnRows(rows)

		session, err := service.GetSession(context.Background(), sessionID)

		assert.NoError(t, err)
		assert.Equal(t, sessionID, session.ID)
		assert.Equal(t, int32(123), session.MemberID)
		assert.Equal(t, "test@dlsu.edu.ph", session.Email)
		assert.Equal(t, "Test User", session.FullName)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("session not found", func(t *testing.T) {
		sessionID := "nonexistent"

		mock.ExpectQuery("SELECT").
			WithArgs(sessionID).
			WillReturnError(sql.ErrNoRows)

		session, err := service.GetSession(context.Background(), sessionID)

		assert.Error(t, err)
		assert.Nil(t, session)
		assert.Equal(t, sql.ErrNoRows, err)

		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestSessionService_DeleteSession(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	service := NewSessionService(db, nil)

	t.Run("delete existing session", func(t *testing.T) {
		sessionID := "abc123"

		mock.ExpectExec("DELETE FROM sessions").
			WithArgs(sessionID).
			WillReturnResult(sqlmock.NewResult(0, 1))

		err := service.DeleteSession(context.Background(), sessionID)

		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestSessionService_ShouldExtendSession(t *testing.T) {
	service := &sessionService{
		slidingExtendThreshold: 0.5,
	}

	t.Run("should extend - less than 50% time remaining", func(t *testing.T) {
		session := &SessionWithMember{
			Session: Session{
				ExpiresAt: time.Now().Add(11 * time.Hour), // less than 12h remaining for 24h session
			},
		}
		duration := 24 * time.Hour

		result := service.ShouldExtendSession(session, duration)

		assert.True(t, result)
	})

	t.Run("should not extend - more than 50% time remaining", func(t *testing.T) {
		session := &SessionWithMember{
			Session: Session{
				ExpiresAt: time.Now().Add(13 * time.Hour), // more than 12h remaining for 24h session
			},
		}
		duration := 24 * time.Hour

		result := service.ShouldExtendSession(session, duration)

		assert.False(t, result)
	})
}

func TestGenerateSessionID(t *testing.T) {
	t.Run("generates unique IDs", func(t *testing.T) {
		id1, err1 := generateSessionID()
		id2, err2 := generateSessionID()

		assert.NoError(t, err1)
		assert.NoError(t, err2)
		assert.NotEqual(t, id1, id2)
		assert.Len(t, id1, 64) // 32 bytes = 64 hex chars
		assert.Len(t, id2, 64)
	})
}
