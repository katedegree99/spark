package handler_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/katedegree99/spark/api/internal/adapter/handler"
	"github.com/katedegree99/spark/api/internal/adapter/router"
	"github.com/katedegree99/spark/api/internal/usecase"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mockAuthUsecase implements usecase.AuthUsecase for testing.
type mockAuthUsecase struct {
	register        func(ctx context.Context, email, password string) error
	login           func(ctx context.Context, email, password string) error
	verifyOTP       func(ctx context.Context, email, code string) (*usecase.TokenPair, error)
	loginWithGoogle func(ctx context.Context, idToken string) (*usecase.TokenPair, error)
	refreshToken    func(ctx context.Context, refreshToken string) (*usecase.TokenPair, error)
	logout          func(ctx context.Context, refreshToken string) error
}

func (m *mockAuthUsecase) Register(ctx context.Context, email, password string) error {
	return m.register(ctx, email, password)
}
func (m *mockAuthUsecase) Login(ctx context.Context, email, password string) error {
	return m.login(ctx, email, password)
}
func (m *mockAuthUsecase) VerifyOTP(ctx context.Context, email, code string) (*usecase.TokenPair, error) {
	return m.verifyOTP(ctx, email, code)
}
func (m *mockAuthUsecase) LoginWithGoogle(ctx context.Context, idToken string) (*usecase.TokenPair, error) {
	return m.loginWithGoogle(ctx, idToken)
}
func (m *mockAuthUsecase) RefreshToken(ctx context.Context, refreshToken string) (*usecase.TokenPair, error) {
	return m.refreshToken(ctx, refreshToken)
}
func (m *mockAuthUsecase) Logout(ctx context.Context, refreshToken string) error {
	return m.logout(ctx, refreshToken)
}

func newTestEcho(uc usecase.AuthUsecase) *echo.Echo {
	h := &handler.Handler{
		AuthHandler: handler.NewAuthHandler(uc),
	}
	return router.NewRouter(h, handler.NewWsHub())
}

const testJWTSecret = "test-secret"

func testBearerToken(userID uint) string {
	os.Setenv("JWT_SECRET", testJWTSecret)
	claims := jwt.MapClaims{
		"sub": float64(userID),
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(15 * time.Minute).Unix(),
	}
	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, _ := tok.SignedString([]byte(testJWTSecret))
	return "Bearer " + signed
}

func do(e *echo.Echo, method, path, body string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	return rec
}

func doWithAuth(e *echo.Echo, method, path, body string, userID uint) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set("Authorization", testBearerToken(userID))
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	return rec
}

// --- POST /auth/register ---

func TestRegisterWithEmail_Success(t *testing.T) {
	uc := &mockAuthUsecase{
		register: func(_ context.Context, _, _ string) error { return nil },
	}
	rec := do(newTestEcho(uc), http.MethodPost, "/auth/register",
		`{"email":"test@example.com","password":"password123"}`)

	assert.Equal(t, http.StatusOK, rec.Code)
	require.Contains(t, rec.Body.String(), "OTP sent")
}

func TestRegisterWithEmail_EmailAlreadyExists(t *testing.T) {
	uc := &mockAuthUsecase{
		register: func(_ context.Context, _, _ string) error { return usecase.ErrEmailAlreadyExists },
	}
	rec := do(newTestEcho(uc), http.MethodPost, "/auth/register",
		`{"email":"existing@example.com","password":"password123"}`)

	assert.Equal(t, http.StatusConflict, rec.Code)
	require.Contains(t, rec.Body.String(), "EMAIL_ALREADY_EXISTS")
}

func TestRegisterWithEmail_InternalError(t *testing.T) {
	uc := &mockAuthUsecase{
		register: func(_ context.Context, _, _ string) error { return errors.New("db error") },
	}
	rec := do(newTestEcho(uc), http.MethodPost, "/auth/register",
		`{"email":"test@example.com","password":"password123"}`)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

// --- POST /auth/login ---

func TestLoginWithEmail_Success(t *testing.T) {
	uc := &mockAuthUsecase{
		login: func(_ context.Context, _, _ string) error { return nil },
	}
	rec := do(newTestEcho(uc), http.MethodPost, "/auth/login",
		`{"email":"test@example.com","password":"password123"}`)

	assert.Equal(t, http.StatusOK, rec.Code)
	require.Contains(t, rec.Body.String(), "OTP sent")
}

func TestLoginWithEmail_InvalidCredentials(t *testing.T) {
	uc := &mockAuthUsecase{
		login: func(_ context.Context, _, _ string) error { return usecase.ErrInvalidCredentials },
	}
	rec := do(newTestEcho(uc), http.MethodPost, "/auth/login",
		`{"email":"test@example.com","password":"wrong"}`)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
	require.Contains(t, rec.Body.String(), "INVALID_CREDENTIALS")
}

func TestLoginWithEmail_InternalError(t *testing.T) {
	uc := &mockAuthUsecase{
		login: func(_ context.Context, _, _ string) error { return errors.New("db error") },
	}
	rec := do(newTestEcho(uc), http.MethodPost, "/auth/login",
		`{"email":"test@example.com","password":"password123"}`)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

// --- POST /auth/otp/verify ---

func TestVerifyOtp_Success(t *testing.T) {
	uc := &mockAuthUsecase{
		verifyOTP: func(_ context.Context, _, _ string) (*usecase.TokenPair, error) {
			return &usecase.TokenPair{AccessToken: "access-tok", RefreshToken: "refresh-tok"}, nil
		},
	}
	rec := do(newTestEcho(uc), http.MethodPost, "/auth/otp/verify",
		`{"email":"test@example.com","code":"123456"}`)

	assert.Equal(t, http.StatusOK, rec.Code)
	require.Contains(t, rec.Body.String(), "access-tok")
	require.Contains(t, rec.Body.String(), "refresh-tok")
}

func TestVerifyOtp_InvalidOTP(t *testing.T) {
	uc := &mockAuthUsecase{
		verifyOTP: func(_ context.Context, _, _ string) (*usecase.TokenPair, error) {
			return nil, usecase.ErrInvalidOTP
		},
	}
	rec := do(newTestEcho(uc), http.MethodPost, "/auth/otp/verify",
		`{"email":"test@example.com","code":"000000"}`)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	require.Contains(t, rec.Body.String(), "INVALID_OTP")
}

func TestVerifyOtp_InternalError(t *testing.T) {
	uc := &mockAuthUsecase{
		verifyOTP: func(_ context.Context, _, _ string) (*usecase.TokenPair, error) {
			return nil, errors.New("db error")
		},
	}
	rec := do(newTestEcho(uc), http.MethodPost, "/auth/otp/verify",
		`{"email":"test@example.com","code":"123456"}`)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

// --- POST /auth/google ---

func TestLoginWithGoogle_Success(t *testing.T) {
	uc := &mockAuthUsecase{
		loginWithGoogle: func(_ context.Context, _ string) (*usecase.TokenPair, error) {
			return &usecase.TokenPair{AccessToken: "access-tok", RefreshToken: "refresh-tok"}, nil
		},
	}
	rec := do(newTestEcho(uc), http.MethodPost, "/auth/google",
		`{"id_token":"valid-google-token"}`)

	assert.Equal(t, http.StatusOK, rec.Code)
	require.Contains(t, rec.Body.String(), "access-tok")
}

func TestLoginWithGoogle_InvalidToken(t *testing.T) {
	uc := &mockAuthUsecase{
		loginWithGoogle: func(_ context.Context, _ string) (*usecase.TokenPair, error) {
			return nil, usecase.ErrInvalidGoogleToken
		},
	}
	rec := do(newTestEcho(uc), http.MethodPost, "/auth/google",
		`{"id_token":"bad-token"}`)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
	require.Contains(t, rec.Body.String(), "INVALID_GOOGLE_TOKEN")
}

// --- POST /auth/refresh ---

func TestRefreshToken_Success(t *testing.T) {
	uc := &mockAuthUsecase{
		refreshToken: func(_ context.Context, _ string) (*usecase.TokenPair, error) {
			return &usecase.TokenPair{AccessToken: "new-access", RefreshToken: "new-refresh"}, nil
		},
	}
	rec := do(newTestEcho(uc), http.MethodPost, "/auth/refresh",
		`{"refresh_token":"valid-refresh-token"}`)

	assert.Equal(t, http.StatusOK, rec.Code)
	require.Contains(t, rec.Body.String(), "new-access")
}

func TestRefreshToken_InvalidToken(t *testing.T) {
	uc := &mockAuthUsecase{
		refreshToken: func(_ context.Context, _ string) (*usecase.TokenPair, error) {
			return nil, usecase.ErrInvalidRefreshToken
		},
	}
	rec := do(newTestEcho(uc), http.MethodPost, "/auth/refresh",
		`{"refresh_token":"expired-token"}`)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
	require.Contains(t, rec.Body.String(), "INVALID_REFRESH_TOKEN")
}

// --- POST /auth/logout ---

func TestLogout_Success(t *testing.T) {
	uc := &mockAuthUsecase{
		logout: func(_ context.Context, _ string) error { return nil },
	}
	rec := doWithAuth(newTestEcho(uc), http.MethodPost, "/auth/logout",
		`{"refresh_token":"valid-refresh-token"}`, 1)

	assert.Equal(t, http.StatusNoContent, rec.Code)
}

func TestLogout_InternalError(t *testing.T) {
	uc := &mockAuthUsecase{
		logout: func(_ context.Context, _ string) error { return errors.New("db error") },
	}
	rec := doWithAuth(newTestEcho(uc), http.MethodPost, "/auth/logout",
		`{"refresh_token":"some-token"}`, 1)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}
