package repository

import "context"

type ImageRepository interface {
	Create(ctx context.Context, directory, url string) (*ImageRecord, error)
	FindByID(ctx context.Context, id uint) (*ImageRecord, error)
}

type ImageRecord struct {
	ID        uint
	Directory string
	URL       string
}
