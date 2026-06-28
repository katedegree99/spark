package usecase

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	mathrand "math/rand/v2"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/katedegree99/spark/api/internal/domain/email"
	"github.com/katedegree99/spark/api/internal/domain/repository"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/api/idtoken"
)


var (
	ErrEmailAlreadyExists  = errors.New("email already exists")
	ErrInvalidCredentials  = errors.New("invalid credentials")
	ErrInvalidOTP          = errors.New("invalid or expired OTP")
	ErrInvalidRefreshToken = errors.New("invalid or expired refresh token")
	ErrInvalidGoogleToken  = errors.New("invalid google id token")
)

type AuthUsecase interface {
	Register(ctx context.Context, email, password string) error
	Login(ctx context.Context, email, password string) error
	VerifyOTP(ctx context.Context, email, code string) (*TokenPair, error)
	LoginWithGoogle(ctx context.Context, idToken string) (*TokenPair, error)
	RefreshToken(ctx context.Context, refreshToken string) (*TokenPair, error)
	Logout(ctx context.Context, refreshToken string) error
}

type TokenPair struct {
	AccessToken   string
	RefreshToken  string
	ProfileExists bool
}

// googleTokenValidator abstracts idtoken.Validate for testability.
type googleTokenValidator func(ctx context.Context, idToken, audience string) (email string, err error)

type authUsecase struct {
	authRepo            repository.AuthRepository
	profileRepo         repository.ProfileRepository
	emailSvc            email.EmailService
	validateGoogleToken googleTokenValidator
}

func NewAuthUsecase(authRepo repository.AuthRepository, profileRepo repository.ProfileRepository, emailSvc email.EmailService) AuthUsecase {
	return &authUsecase{
		authRepo:    authRepo,
		profileRepo: profileRepo,
		emailSvc:    emailSvc,
		validateGoogleToken: func(ctx context.Context, idToken, audience string) (string, error) {
			payload, err := idtoken.Validate(ctx, idToken, audience)
			if err != nil {
				return "", err
			}
			email, ok := payload.Claims["email"].(string)
			if !ok || email == "" {
				return "", errors.New("no email claim")
			}
			return email, nil
		},
	}
}

func (u *authUsecase) Register(ctx context.Context, email, password string) error {
	existing, err := u.authRepo.FindUserByEmail(ctx, email)
	if err != nil {
		return err
	}
	if existing != nil {
		return ErrEmailAlreadyExists
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	if _, err = u.authRepo.CreateUser(ctx, email, string(hashed)); err != nil {
		return err
	}

	code := fmt.Sprintf("%06d", mathrand.IntN(1000000))
	if err := u.authRepo.SaveOTP(ctx, email, code); err != nil {
		_ = u.authRepo.DeleteUserByEmail(ctx, email)
		return err
	}
	if err := u.emailSvc.SendOTP(ctx, email, code); err != nil {
		_ = u.authRepo.DeleteUserByEmail(ctx, email)
		return err
	}
	return nil
}

func (u *authUsecase) Login(ctx context.Context, email, password string) error {
	user, err := u.authRepo.FindUserByEmail(ctx, email)
	if err != nil {
		return err
	}
	if user == nil {
		return ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(password)); err != nil {
		return ErrInvalidCredentials
	}

	code := fmt.Sprintf("%06d", mathrand.IntN(1000000))
	if err := u.authRepo.SaveOTP(ctx, email, code); err != nil {
		return err
	}
	return u.emailSvc.SendOTP(ctx, email, code)
}

func (u *authUsecase) VerifyOTP(ctx context.Context, email, code string) (*TokenPair, error) {
	valid, err := u.authRepo.VerifyOTP(ctx, email, code)
	if err != nil {
		return nil, err
	}
	if !valid {
		return nil, ErrInvalidOTP
	}

	user, err := u.authRepo.FindUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrInvalidOTP
	}

	return u.issueTokenPair(ctx, user.ID)
}

func (u *authUsecase) LoginWithGoogle(ctx context.Context, rawIDToken string) (*TokenPair, error) {
	clientID := os.Getenv("GOOGLE_CLIENT_ID")
	email, err := u.validateGoogleToken(ctx, rawIDToken, clientID)
	if err != nil {
		return nil, ErrInvalidGoogleToken
	}

	user, err := u.authRepo.FindOrCreateUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	return u.issueTokenPair(ctx, user.ID)
}

func (u *authUsecase) RefreshToken(ctx context.Context, refreshToken string) (*TokenPair, error) {
	record, err := u.authRepo.FindRefreshToken(ctx, refreshToken)
	if err != nil {
		return nil, err
	}
	if record == nil || time.Now().After(record.ExpiresAt) {
		return nil, ErrInvalidRefreshToken
	}

	return u.issueTokenPair(ctx, record.UserID)
}

func (u *authUsecase) Logout(ctx context.Context, refreshToken string) error {
	return u.authRepo.DeleteRefreshToken(ctx, refreshToken)
}

func (u *authUsecase) issueTokenPair(ctx context.Context, userID uint) (*TokenPair, error) {
	secret := os.Getenv("JWT_SECRET")
	now := time.Now()
	claims := jwt.MapClaims{
		"sub": userID,
		"iat": now.Unix(),
		"exp": now.Add(15 * time.Minute).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accessToken, err := token.SignedString([]byte(secret))
	if err != nil {
		return nil, err
	}

	rawBytes := make([]byte, 32)
	if _, err := rand.Read(rawBytes); err != nil {
		return nil, err
	}
	refreshToken := hex.EncodeToString(rawBytes)

	expiresAt := now.Add(7 * 24 * time.Hour)
	if err := u.authRepo.SaveRefreshToken(ctx, userID, refreshToken, expiresAt); err != nil {
		return nil, err
	}

	profileExists, err := u.profileRepo.ExistsByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &TokenPair{
		AccessToken:   accessToken,
		RefreshToken:  refreshToken,
		ProfileExists: profileExists,
	}, nil
}
