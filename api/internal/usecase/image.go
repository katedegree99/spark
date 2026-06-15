package usecase

import (
	"context"
	"io"

	"github.com/katedegree/spark/api/internal/domain/repository"
	"github.com/katedegree/spark/api/internal/domain/storage"
)

type ImageUsecase interface {
	Upload(ctx context.Context, directory, filename string, r io.Reader) (*repository.ImageRecord, error)
}

type imageUsecase struct {
	imageRepo      repository.ImageRepository
	storageService storage.StorageService
}

func NewImageUsecase(imageRepo repository.ImageRepository, storageSvc storage.StorageService) ImageUsecase {
	return &imageUsecase{imageRepo: imageRepo, storageService: storageSvc}
}

func (u *imageUsecase) Upload(ctx context.Context, directory, filename string, r io.Reader) (*repository.ImageRecord, error) {
	panic("not implemented")
}
