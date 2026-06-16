package repository

import (
	"context"
	"time"
)

type ThingRepository interface {
	Search(ctx context.Context, q string) ([]*ThingRecord, error)
	FindByID(ctx context.Context, id uint) (*ThingRecord, error)
	FindByName(ctx context.Context, name string) (*ThingRecord, error)
	FindByIDs(ctx context.Context, ids []uint) ([]*ThingRecord, error)
	Create(ctx context.Context, name string) (*ThingRecord, error)
	FindAliasesByThingID(ctx context.Context, thingID uint) ([]string, error)
	HasAlias(ctx context.Context, thingID uint) (bool, error)
	AddAlias(ctx context.Context, thingID uint, alias string) error
	UpdateName(ctx context.Context, thingID uint, name string) error
	FindAllWithAliases(ctx context.Context) ([]*ThingRecord, error)
}

type ThingRecord struct {
	ID        uint
	Name      string
	Aliases   []string
	CreatedAt time.Time
}
