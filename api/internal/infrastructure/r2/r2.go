package r2

import (
	"context"
	"io"

	"github.com/katedegree/spark/api/internal/domain/storage"
)

type r2Service struct{}

func NewR2Service() storage.StorageService {
	return &r2Service{}
}

func (s *r2Service) Upload(ctx context.Context, directory, filename string, r io.Reader) (string, error) {
	panic("not implemented")
}
