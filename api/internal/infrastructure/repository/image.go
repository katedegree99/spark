package repository

import (
	"context"
	"errors"

	"github.com/katedegree/spark/api/internal/domain/model"
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
	img := model.Image{
		Directory: directory,
		URL:       url,
	}
	if err := r.db.WithContext(ctx).Create(&img).Error; err != nil {
		return nil, err
	}
	return &domainrepo.ImageRecord{
		ID:        img.ID,
		Directory: img.Directory,
		URL:       img.URL,
	}, nil
}

func (r *imageRepository) FindByID(ctx context.Context, id uint) (*domainrepo.ImageRecord, error) {
	var img model.Image
	err := r.db.WithContext(ctx).First(&img, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &domainrepo.ImageRecord{
		ID:        img.ID,
		Directory: img.Directory,
		URL:       img.URL,
	}, nil
}
