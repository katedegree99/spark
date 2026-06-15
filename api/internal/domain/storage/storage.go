package storage

import (
	"context"
	"io"
)

type StorageService interface {
	Upload(ctx context.Context, directory, filename string, r io.Reader) (url string, err error)
}
