package repository

import (
	"context"
	"time"
)

type PickupCandidateRecord struct {
	UserID      uint
	Name        string
	Bio         *string
	IconImageID *uint
	ThingIDs    []uint
}

type PickupRepository interface {
	FindCandidates(ctx context.Context, excludeUserID uint) ([]*PickupCandidateRecord, error)
	FindCache(ctx context.Context, userID uint, date time.Time) ([]uint, error)
	SaveCache(ctx context.Context, userID uint, date time.Time, pickedUserIDs []uint) error
}
