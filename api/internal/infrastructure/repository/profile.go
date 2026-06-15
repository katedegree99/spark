package repository

import (
	"context"

	domainrepo "github.com/katedegree/spark/api/internal/domain/repository"
	"gorm.io/gorm"
)

type profileRepository struct {
	db *gorm.DB
}

func NewProfileRepository(db *gorm.DB) domainrepo.ProfileRepository {
	return &profileRepository{db: db}
}

func (r *profileRepository) ExistsByUserID(ctx context.Context, userID uint) (bool, error) {
	panic("not implemented")
}

func (r *profileRepository) Create(ctx context.Context, p *domainrepo.ProfileRecord) error {
	panic("not implemented")
}

func (r *profileRepository) FindByUserID(ctx context.Context, userID uint) (*domainrepo.ProfileRecord, error) {
	panic("not implemented")
}

func (r *profileRepository) Update(ctx context.Context, p *domainrepo.ProfileRecord) error {
	panic("not implemented")
}

func (r *profileRepository) SetDoings(ctx context.Context, userID uint, thingIDs []uint) error {
	panic("not implemented")
}

func (r *profileRepository) SetWants(ctx context.Context, userID uint, thingIDs []uint) error {
	panic("not implemented")
}
