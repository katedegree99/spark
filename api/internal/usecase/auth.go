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
	"github.com/katedegree/spark/api/internal/domain/repository"
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
	AccessToken  string
	RefreshToken string
}

type authUsecase struct {
	authRepo repository.AuthRepository
}

func NewAuthUsecase(authRepo repository.AuthRepository) AuthUsecase {
	return &authUsecase{authRepo: authRepo}
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
	return u.authRepo.SaveOTP(ctx, email, code)
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
	return u.authRepo.SaveOTP(ctx, email, code)
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
	payload, err := idtoken.Validate(ctx, rawIDToken, clientID)
	if err != nil {
		return nil, ErrInvalidGoogleToken
	}

	emailVal, ok := payload.Claims["email"]
	if !ok {
		return nil, ErrInvalidGoogleToken
	}
	email, ok := emailVal.(string)
	if !ok || email == "" {
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

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
