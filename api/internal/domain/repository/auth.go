package repository

import "context"

type AuthRepository interface {
	FindUserByEmail(ctx context.Context, email string) (*AuthUser, error)
	CreateUser(ctx context.Context, email, hashedPassword string) (*AuthUser, error)
	SaveOTP(ctx context.Context, email, code string) error
	VerifyOTP(ctx context.Context, email, code string) (bool, error)
}

type AuthUser struct {
	ID             uint
	Email          string
	HashedPassword string
}
