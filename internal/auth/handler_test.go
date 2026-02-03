package auth

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

// mockAuthService is a mock implementation of the auth.Service interface.
type mockAuthService struct{}

func (m *mockAuthService) GenerateJWT(email string, keyType KeyType) (string, *time.Time, error) {
	// for dev/prod keys, return an expiration time
	if keyType != KeyTypeAdmin {
		exp := time.Now().Add(30 * 24 * time.Hour)
		return "test_jwt_token", &exp, nil
	}
	// admin keys don't expire
	return "test_jwt_token", nil, nil
}

// mockDBService is a mock implementation of the database.Service interface.
type mockDBService struct {
	db *sql.DB
}

func (m *mockDBService) Health() map[string]string {
	return nil
}

func (m *mockDBService) Close() error {
	return nil
}

func (m *mockDBService) GetConnection() *sql.DB {
	return m.db
}

func TestRequestKeyHandler(t *testing.T) {
	t.Run("success - RND member", func(t *testing.T) {
		e := echo.New()
		// use RequestKeyRequest, not EmailRequest
		reqBody := RequestKeyRequest{
			Project:       "Test Project",
			AllowedOrigin: "",
			IsDev:         false,
			IsAdmin:       true, // admin key so no origin required
		}
		jsonBody, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/request-key", bytes.NewReader(jsonBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		testEmail := "test@dlsu.edu.ph"
		// set user_email in context (set by Google OAuth middleware)
		c.Set("user_email", testEmail)

		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		// mock GetMemberAuthInfo for RBAC check (lightweight query)
		authInfoRow := sqlmock.NewRows([]string{
			"id", "position_id", "committee_id",
		}).AddRow(1, "AVP", "RND")
		mock.ExpectQuery("SELECT id, position_id, committee_id FROM members").
			WithArgs(testEmail).
			WillReturnRows(authInfoRow)

		// GetMemberInfo is called by handler for full member info
		memberRow := sqlmock.NewRows([]string{
			"id", "email", "full_name", "nickname", "image_url",
			"committee_id", "committee_name",
			"division_id", "division_name",
			"position_id", "position_name",
			"house_name",
			"contact_number", "college", "program",
			"interests", "discord", "fb_link", "telegram",
		}).AddRow(
			1, testEmail, "Test User", nil, nil,
			"RND", "Research and Development",
			"INT", "Internal",
			"AVP", "Associate Vice President",
			"Gell-Mann",
			nil, "CCS", "CS-ST",
			nil, nil, nil, nil,
		)
		mock.ExpectQuery("SELECT (.+) FROM members m").
			WithArgs(testEmail).
			WillReturnRows(memberRow)

		mock.ExpectExec("INSERT INTO api_keys").
			WillReturnResult(sqlmock.NewResult(1, 1))

		dbService := &mockDBService{db: db}
		authService := &mockAuthService{}
		rbacService := NewRBACService(dbService)
		h := NewHandler(authService, dbService, rbacService)

		if assert.NoError(t, h.RequestKeyHandler(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
			var resp map[string]interface{}
			json.Unmarshal(rec.Body.Bytes(), &resp)
			assert.Equal(t, testEmail, resp["email"])
			assert.Equal(t, "test_jwt_token", resp["api_key"])
		}
	})

	t.Run("fail - non-RND member", func(t *testing.T) {
		e := echo.New()
		reqBody := RequestKeyRequest{
			Project: "Test Project",
			IsAdmin: true,
		}
		jsonBody, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/request-key", bytes.NewReader(jsonBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		testEmail := "test@dlsu.edu.ph"
		// set user_email in context
		c.Set("user_email", testEmail)

		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		// mock GetMemberAuthInfo for RBAC check - non-RND, non-AVP+ member
		authInfoRow := sqlmock.NewRows([]string{
			"id", "position_id", "committee_id",
		}).AddRow(1, "MEM", "EXT") // not RND, not AVP+
		mock.ExpectQuery("SELECT id, position_id, committee_id FROM members").
			WithArgs(testEmail).
			WillReturnRows(authInfoRow)

		dbService := &mockDBService{db: db}
		authService := &mockAuthService{}
		rbacService := NewRBACService(dbService)
		h := NewHandler(authService, dbService, rbacService)

		if assert.NoError(t, h.RequestKeyHandler(c)) {
			assert.Equal(t, http.StatusForbidden, rec.Code)
		}
	})

	t.Run("fail - not an LSCS member", func(t *testing.T) {
		e := echo.New()
		reqBody := RequestKeyRequest{
			Project: "Test Project",
			IsAdmin: true,
		}
		jsonBody, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/request-key", bytes.NewReader(jsonBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		testEmail := "test@dlsu.edu.ph"
		// set user_email in context
		c.Set("user_email", testEmail)

		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		// mock GetMemberAuthInfo for RBAC check - no rows (not a member)
		mock.ExpectQuery("SELECT id, position_id, committee_id FROM members").
			WithArgs(testEmail).
			WillReturnError(sql.ErrNoRows)

		dbService := &mockDBService{db: db}
		authService := &mockAuthService{}
		rbacService := NewRBACService(dbService)
		h := NewHandler(authService, dbService, rbacService)

		if assert.NoError(t, h.RequestKeyHandler(c)) {
			assert.Equal(t, http.StatusForbidden, rec.Code)
		}
	})

	t.Run("fail - user_email not in context", func(t *testing.T) {
		e := echo.New()
		reqBody := RequestKeyRequest{
			Project: "Test Project",
			IsAdmin: true,
		}
		jsonBody, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/request-key", bytes.NewReader(jsonBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		// don't set user_email in context

		db, _, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		dbService := &mockDBService{db: db}
		authService := &mockAuthService{}
		rbacService := NewRBACService(dbService)
		h := NewHandler(authService, dbService, rbacService)

		if assert.NoError(t, h.RequestKeyHandler(c)) {
			assert.Equal(t, http.StatusUnauthorized, rec.Code)
		}
	})
}
