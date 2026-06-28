package repository

import "context"

type PickupCandidateRecord struct {
	UserID      uint
	Name        string
	Bio         *string
	IconImageID *uint
	ThingIDs    []uint
}

type PickupRepository interface {
	FindCandidates(ctx context.Context, excludeUserID uint) ([]*PickupCandidateRecord, error)
}
