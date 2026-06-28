package repository

import "context"

type NewUserRecord struct {
	UserID  uint
	Name    string
	IconURL *string
}

type NewUserRepository interface {
	FindNew(ctx context.Context, excludeUserID uint, limit int) ([]*NewUserRecord, error)
}
