package repository

import "context"

type InterestRecord struct {
	UserID      uint
	Name        string
	Bio         *string
	IconImageID *uint
	ThingIDs    []uint
}

type InterestRepository interface {
	Exists(ctx context.Context, fromUserID, toUserID uint) (bool, error)
	Create(ctx context.Context, fromUserID, toUserID uint) error
	ListSent(ctx context.Context, fromUserID uint, offset, limit int) ([]*InterestRecord, error)
	ListReceived(ctx context.Context, toUserID uint, offset, limit int) ([]*InterestRecord, error)
}
