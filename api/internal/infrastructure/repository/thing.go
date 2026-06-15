package repository

import (
	"context"

	domainrepo "github.com/katedegree/spark/api/internal/domain/repository"
	"gorm.io/gorm"
)

type thingRepository struct {
	db *gorm.DB
}

func NewThingRepository(db *gorm.DB) domainrepo.ThingRepository {
	return &thingRepository{db: db}
}

func (r *thingRepository) Search(ctx context.Context, q string) ([]*domainrepo.ThingRecord, error) {
	panic("not implemented")
}

func (r *thingRepository) FindByID(ctx context.Context, id uint) (*domainrepo.ThingRecord, error) {
	panic("not implemented")
}

func (r *thingRepository) FindByName(ctx context.Context, name string) (*domainrepo.ThingRecord, error) {
	panic("not implemented")
}

func (r *thingRepository) FindByIDs(ctx context.Context, ids []uint) ([]*domainrepo.ThingRecord, error) {
	panic("not implemented")
}

func (r *thingRepository) Create(ctx context.Context, name string) (*domainrepo.ThingRecord, error) {
	panic("not implemented")
}

func (r *thingRepository) FindAliasesByThingID(ctx context.Context, thingID uint) ([]string, error) {
	panic("not implemented")
}
