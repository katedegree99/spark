package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/katedegree99/spark/api/internal/domain/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

type mockEmailService struct {
	sendOTP func(ctx context.Context, to, code string) error
}

func (m *mockEmailService) SendOTP(ctx context.Context, to, code string) error {
	if m.sendOTP != nil {
		return m.sendOTP(ctx, to, code)
	}
	return nil
}

// mockAuthRepository implements repository.AuthRepository for testing.
type mockAuthRepository struct {
	findUserByEmail         func(ctx context.Context, email string) (*repository.AuthUser, error)
	createUser              func(ctx context.Context, email, hashedPassword string) (*repository.AuthUser, error)
	deleteUserByEmail       func(ctx context.Context, email string) error
	findOrCreateUserByEmail func(ctx context.Context, email string) (*repository.AuthUser, error)
	saveOTP                 func(ctx context.Context, email, code string) error
	verifyOTP               func(ctx context.Context, email, code string) (bool, error)
	saveRefreshToken        func(ctx context.Context, userID uint, token string, expiresAt time.Time) error
	findRefreshToken        func(ctx context.Context, token string) (*repository.RefreshTokenRecord, error)
	deleteRefreshToken      func(ctx context.Context, token string) error
}

func (m *mockAuthRepository) FindUserByEmail(ctx context.Context, email string) (*repository.AuthUser, error) {
	return m.findUserByEmail(ctx, email)
}
func (m *mockAuthRepository) CreateUser(ctx context.Context, email, hashedPassword string) (*repository.AuthUser, error) {
	return m.createUser(ctx, email, hashedPassword)
}
func (m *mockAuthRepository) DeleteUserByEmail(ctx context.Context, email string) error {
	if m.deleteUserByEmail != nil {
		return m.deleteUserByEmail(ctx, email)
	}
	return nil
}
func (m *mockAuthRepository) FindOrCreateUserByEmail(ctx context.Context, email string) (*repository.AuthUser, error) {
	return m.findOrCreateUserByEmail(ctx, email)
}
func (m *mockAuthRepository) SaveOTP(ctx context.Context, email, code string) error {
	return m.saveOTP(ctx, email, code)
}
func (m *mockAuthRepository) VerifyOTP(ctx context.Context, email, code string) (bool, error) {
	return m.verifyOTP(ctx, email, code)
}
func (m *mockAuthRepository) SaveRefreshToken(ctx context.Context, userID uint, token string, expiresAt time.Time) error {
	return m.saveRefreshToken(ctx, userID, token, expiresAt)
}
func (m *mockAuthRepository) FindRefreshToken(ctx context.Context, token string) (*repository.RefreshTokenRecord, error) {
	return m.findRefreshToken(ctx, token)
}
func (m *mockAuthRepository) DeleteRefreshToken(ctx context.Context, token string) error {
	return m.deleteRefreshToken(ctx, token)
}

// mockProfileRepository implements repository.ProfileRepository for testing.
type mockProfileRepository struct {
	existsByUserID      func(ctx context.Context, userID uint) (bool, error)
	create              func(ctx context.Context, p *repository.ProfileRecord) error
	findByUserID        func(ctx context.Context, userID uint) (*repository.ProfileRecord, error)
	update              func(ctx context.Context, p *repository.ProfileRecord) error
	setDoings           func(ctx context.Context, userID uint, thingIDs []uint) error
	setWants            func(ctx context.Context, userID uint, thingIDs []uint) error
	findDoingIDsByUserID func(ctx context.Context, userID uint) ([]uint, error)
	findWantIDsByUserID  func(ctx context.Context, userID uint) ([]uint, error)
}

func (m *mockProfileRepository) ExistsByUserID(ctx context.Context, userID uint) (bool, error) {
	if m.existsByUserID != nil {
		return m.existsByUserID(ctx, userID)
	}
	return false, nil
}
func (m *mockProfileRepository) Create(ctx context.Context, p *repository.ProfileRecord) error {
	if m.create != nil {
		return m.create(ctx, p)
	}
	return nil
}
func (m *mockProfileRepository) FindByUserID(ctx context.Context, userID uint) (*repository.ProfileRecord, error) {
	if m.findByUserID != nil {
		return m.findByUserID(ctx, userID)
	}
	return nil, nil
}
func (m *mockProfileRepository) Update(ctx context.Context, p *repository.ProfileRecord) error {
	if m.update != nil {
		return m.update(ctx, p)
	}
	return nil
}
func (m *mockProfileRepository) SetDoings(ctx context.Context, userID uint, thingIDs []uint) error {
	if m.setDoings != nil {
		return m.setDoings(ctx, userID, thingIDs)
	}
	return nil
}
func (m *mockProfileRepository) SetWants(ctx context.Context, userID uint, thingIDs []uint) error {
	if m.setWants != nil {
		return m.setWants(ctx, userID, thingIDs)
	}
	return nil
}
func (m *mockProfileRepository) FindDoingIDsByUserID(ctx context.Context, userID uint) ([]uint, error) {
	if m.findDoingIDsByUserID != nil {
		return m.findDoingIDsByUserID(ctx, userID)
	}
	return nil, nil
}
func (m *mockProfileRepository) FindWantIDsByUserID(ctx context.Context, userID uint) ([]uint, error) {
	if m.findWantIDsByUserID != nil {
		return m.findWantIDsByUserID(ctx, userID)
	}
	return nil, nil
}

func newTestUsecase(repo *mockAuthRepository) *authUsecase {
	return &authUsecase{
		authRepo:    repo,
		profileRepo: &mockProfileRepository{},
		emailSvc:    &mockEmailService{},
		validateGoogleToken: func(ctx context.Context, idToken, audience string) (string, error) {
			return "", errors.New("default mock: not configured")
		},
	}
}

// --- Register ---

func TestRegister_Success(t *testing.T) {
	repo := &mockAuthRepository{
		findUserByEmail: func(_ context.Context, _ string) (*repository.AuthUser, error) { return nil, nil },
		createUser: func(_ context.Context, email, _ string) (*repository.AuthUser, error) {
			return &repository.AuthUser{ID: 1, Email: email}, nil
		},
		saveOTP: func(_ context.Context, _, _ string) error { return nil },
	}
	uc := newTestUsecase(repo)

	err := uc.Register(context.Background(), "test@example.com", "password123")
	require.NoError(t, err)
}

func TestRegister_EmailAlreadyExists(t *testing.T) {
	repo := &mockAuthRepository{
		findUserByEmail: func(_ context.Context, email string) (*repository.AuthUser, error) {
			return &repository.AuthUser{ID: 1, Email: email}, nil
		},
	}
	uc := newTestUsecase(repo)

	err := uc.Register(context.Background(), "existing@example.com", "password123")
	assert.ErrorIs(t, err, ErrEmailAlreadyExists)
}

func TestRegister_RepositoryError(t *testing.T) {
	dbErr := errors.New("db error")
	repo := &mockAuthRepository{
		findUserByEmail: func(_ context.Context, _ string) (*repository.AuthUser, error) {
			return nil, dbErr
		},
	}
	uc := newTestUsecase(repo)

	err := uc.Register(context.Background(), "test@example.com", "password123")
	assert.ErrorIs(t, err, dbErr)
}

// --- Login ---

func TestLogin_Success(t *testing.T) {
	hashed, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.MinCost)
	repo := &mockAuthRepository{
		findUserByEmail: func(_ context.Context, email string) (*repository.AuthUser, error) {
			return &repository.AuthUser{ID: 1, Email: email, HashedPassword: string(hashed)}, nil
		},
		saveOTP: func(_ context.Context, _, _ string) error { return nil },
	}
	uc := newTestUsecase(repo)

	err := uc.Login(context.Background(), "test@example.com", "password123")
	require.NoError(t, err)
}

func TestLogin_UserNotFound(t *testing.T) {
	repo := &mockAuthRepository{
		findUserByEmail: func(_ context.Context, _ string) (*repository.AuthUser, error) { return nil, nil },
	}
	uc := newTestUsecase(repo)

	err := uc.Login(context.Background(), "nobody@example.com", "password123")
	assert.ErrorIs(t, err, ErrInvalidCredentials)
}

func TestLogin_WrongPassword(t *testing.T) {
	hashed, _ := bcrypt.GenerateFromPassword([]byte("correct"), bcrypt.MinCost)
	repo := &mockAuthRepository{
		findUserByEmail: func(_ context.Context, email string) (*repository.AuthUser, error) {
			return &repository.AuthUser{ID: 1, Email: email, HashedPassword: string(hashed)}, nil
		},
	}
	uc := newTestUsecase(repo)

	err := uc.Login(context.Background(), "test@example.com", "wrong")
	assert.ErrorIs(t, err, ErrInvalidCredentials)
}

// --- VerifyOTP ---

func TestVerifyOTP_Success(t *testing.T) {
	t.Setenv("JWT_SECRET", "test-secret")

	repo := &mockAuthRepository{
		verifyOTP: func(_ context.Context, _, _ string) (bool, error) { return true, nil },
		findUserByEmail: func(_ context.Context, email string) (*repository.AuthUser, error) {
			return &repository.AuthUser{ID: 1, Email: email}, nil
		},
		saveRefreshToken: func(_ context.Context, _ uint, _ string, _ time.Time) error { return nil },
	}
	uc := newTestUsecase(repo)

	pair, err := uc.VerifyOTP(context.Background(), "test@example.com", "123456")
	require.NoError(t, err)
	assert.NotEmpty(t, pair.AccessToken)
	assert.NotEmpty(t, pair.RefreshToken)
}

func TestVerifyOTP_InvalidCode(t *testing.T) {
	repo := &mockAuthRepository{
		verifyOTP: func(_ context.Context, _, _ string) (bool, error) { return false, nil },
	}
	uc := newTestUsecase(repo)

	_, err := uc.VerifyOTP(context.Background(), "test@example.com", "000000")
	assert.ErrorIs(t, err, ErrInvalidOTP)
}

func TestVerifyOTP_UserNotFound(t *testing.T) {
	repo := &mockAuthRepository{
		verifyOTP:       func(_ context.Context, _, _ string) (bool, error) { return true, nil },
		findUserByEmail: func(_ context.Context, _ string) (*repository.AuthUser, error) { return nil, nil },
	}
	uc := newTestUsecase(repo)

	_, err := uc.VerifyOTP(context.Background(), "gone@example.com", "123456")
	assert.ErrorIs(t, err, ErrInvalidOTP)
}

// --- LoginWithGoogle ---

func TestLoginWithGoogle_Success(t *testing.T) {
	t.Setenv("JWT_SECRET", "test-secret")

	repo := &mockAuthRepository{
		findOrCreateUserByEmail: func(_ context.Context, email string) (*repository.AuthUser, error) {
			return &repository.AuthUser{ID: 2, Email: email}, nil
		},
		saveRefreshToken: func(_ context.Context, _ uint, _ string, _ time.Time) error { return nil },
	}
	uc := newTestUsecase(repo)
	uc.validateGoogleToken = func(_ context.Context, _, _ string) (string, error) {
		return "google@example.com", nil
	}

	pair, err := uc.LoginWithGoogle(context.Background(), "valid-google-token")
	require.NoError(t, err)
	assert.NotEmpty(t, pair.AccessToken)
	assert.NotEmpty(t, pair.RefreshToken)
}

func TestLoginWithGoogle_InvalidToken(t *testing.T) {
	repo := &mockAuthRepository{}
	uc := newTestUsecase(repo)
	uc.validateGoogleToken = func(_ context.Context, _, _ string) (string, error) {
		return "", errors.New("invalid token")
	}

	_, err := uc.LoginWithGoogle(context.Background(), "bad-token")
	assert.ErrorIs(t, err, ErrInvalidGoogleToken)
}

// --- RefreshToken ---

func TestRefreshToken_Success(t *testing.T) {
	t.Setenv("JWT_SECRET", "test-secret")

	repo := &mockAuthRepository{
		findRefreshToken: func(_ context.Context, token string) (*repository.RefreshTokenRecord, error) {
			return &repository.RefreshTokenRecord{
				UserID:    1,
				Token:     token,
				ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
			}, nil
		},
		saveRefreshToken: func(_ context.Context, _ uint, _ string, _ time.Time) error { return nil },
	}
	uc := newTestUsecase(repo)

	pair, err := uc.RefreshToken(context.Background(), "valid-refresh-token")
	require.NoError(t, err)
	assert.NotEmpty(t, pair.AccessToken)
	assert.NotEmpty(t, pair.RefreshToken)
}

func TestRefreshToken_TokenNotFound(t *testing.T) {
	repo := &mockAuthRepository{
		findRefreshToken: func(_ context.Context, _ string) (*repository.RefreshTokenRecord, error) {
			return nil, nil
		},
	}
	uc := newTestUsecase(repo)

	_, err := uc.RefreshToken(context.Background(), "unknown-token")
	assert.ErrorIs(t, err, ErrInvalidRefreshToken)
}

func TestRefreshToken_Expired(t *testing.T) {
	repo := &mockAuthRepository{
		findRefreshToken: func(_ context.Context, token string) (*repository.RefreshTokenRecord, error) {
			return &repository.RefreshTokenRecord{
				UserID:    1,
				Token:     token,
				ExpiresAt: time.Now().Add(-1 * time.Hour),
			}, nil
		},
	}
	uc := newTestUsecase(repo)

	_, err := uc.RefreshToken(context.Background(), "expired-token")
	assert.ErrorIs(t, err, ErrInvalidRefreshToken)
}

// --- Logout ---

func TestLogout_Success(t *testing.T) {
	repo := &mockAuthRepository{
		deleteRefreshToken: func(_ context.Context, _ string) error { return nil },
	}
	uc := newTestUsecase(repo)

	err := uc.Logout(context.Background(), "some-refresh-token")
	require.NoError(t, err)
}

func TestLogout_RepositoryError(t *testing.T) {
	dbErr := errors.New("db error")
	repo := &mockAuthRepository{
		deleteRefreshToken: func(_ context.Context, _ string) error { return dbErr },
	}
	uc := newTestUsecase(repo)

	err := uc.Logout(context.Background(), "some-token")
	assert.ErrorIs(t, err, dbErr)
}
