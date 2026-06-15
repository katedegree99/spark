package repository

import "context"

type ThingRepository interface {
	Search(ctx context.Context, q string) ([]*ThingRecord, error)
	FindByID(ctx context.Context, id uint) (*ThingRecord, error)
	FindByName(ctx context.Context, name string) (*ThingRecord, error)
	FindByIDs(ctx context.Context, ids []uint) ([]*ThingRecord, error)
	Create(ctx context.Context, name string) (*ThingRecord, error)
	FindAliasesByThingID(ctx context.Context, thingID uint) ([]string, error)
}

type ThingRecord struct {
	ID      uint
	Name    string
	Aliases []string
}
