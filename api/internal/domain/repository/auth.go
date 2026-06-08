package repository

import (
	"context"
	"time"
)

type AuthRepository interface {
	FindUserByEmail(ctx context.Context, email string) (*AuthUser, error)
	CreateUser(ctx context.Context, email, hashedPassword string) (*AuthUser, error)
	DeleteUserByEmail(ctx context.Context, email string) error
	FindOrCreateUserByEmail(ctx context.Context, email string) (*AuthUser, error)
	SaveOTP(ctx context.Context, email, code string) error
	VerifyOTP(ctx context.Context, email, code string) (bool, error)
	SaveRefreshToken(ctx context.Context, userID uint, token string, expiresAt time.Time) error
	FindRefreshToken(ctx context.Context, token string) (*RefreshTokenRecord, error)
	DeleteRefreshToken(ctx context.Context, token string) error
}

type AuthUser struct {
	ID             uint
	Email          string
	HashedPassword string
}

type RefreshTokenRecord struct {
	UserID    uint
	Token     string
	ExpiresAt time.Time
}
