package usecase

import (
	"context"

	"github.com/katedegree/spark/api/internal/domain/repository"
)

type AuthUsecase interface {
	Register(ctx context.Context, email, password string) error
	Login(ctx context.Context, email, password string) error
	VerifyOTP(ctx context.Context, email, code string) (*TokenPair, error)
	LoginWithGoogle(ctx context.Context, idToken string) (*TokenPair, error)
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
	panic("not implemented")
}

func (u *authUsecase) Login(ctx context.Context, email, password string) error {
	panic("not implemented")
}

func (u *authUsecase) VerifyOTP(ctx context.Context, email, code string) (*TokenPair, error) {
	panic("not implemented")
}

func (u *authUsecase) LoginWithGoogle(ctx context.Context, idToken string) (*TokenPair, error) {
	panic("not implemented")
}
