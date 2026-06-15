package repository

import (
	"context"

	domainrepo "github.com/katedegree/spark/api/internal/domain/repository"
	"gorm.io/gorm"
)

type imageRepository struct {
	db *gorm.DB
}

func NewImageRepository(db *gorm.DB) domainrepo.ImageRepository {
	return &imageRepository{db: db}
}

func (r *imageRepository) Create(ctx context.Context, directory, url string) (*domainrepo.ImageRecord, error) {
	panic("not implemented")
}

func (r *imageRepository) FindByID(ctx context.Context, id uint) (*domainrepo.ImageRecord, error) {
	panic("not implemented")
}
