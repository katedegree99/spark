package repository

import (
	"context"

	domainrepo "github.com/katedegree/spark/api/internal/domain/repository"
	"gorm.io/gorm"
)

type authRepository struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) domainrepo.AuthRepository {
	return &authRepository{db: db}
}

func (r *authRepository) FindUserByEmail(ctx context.Context, email string) (*domainrepo.AuthUser, error) {
	panic("not implemented")
}

func (r *authRepository) CreateUser(ctx context.Context, email, hashedPassword string) (*domainrepo.AuthUser, error) {
	panic("not implemented")
}

func (r *authRepository) SaveOTP(ctx context.Context, email, code string) error {
	panic("not implemented")
}

func (r *authRepository) VerifyOTP(ctx context.Context, email, code string) (bool, error) {
	panic("not implemented")
}
