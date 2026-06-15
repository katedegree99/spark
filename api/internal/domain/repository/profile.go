package repository

import "context"

type ProfileRepository interface {
	ExistsByUserID(ctx context.Context, userID uint) (bool, error)
	Create(ctx context.Context, p *ProfileRecord) error
	FindByUserID(ctx context.Context, userID uint) (*ProfileRecord, error)
	Update(ctx context.Context, p *ProfileRecord) error
	SetDoings(ctx context.Context, userID uint, thingIDs []uint) error
	SetWants(ctx context.Context, userID uint, thingIDs []uint) error
}

type ProfileRecord struct {
	UserID      uint
	Name        string
	Bio         *string
	IconImageID *uint
}
