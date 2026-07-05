package repository

import "context"

type InterestRepository interface {
	Exists(ctx context.Context, fromUserID, toUserID uint) (bool, error)
	Create(ctx context.Context, fromUserID, toUserID uint) error
}
