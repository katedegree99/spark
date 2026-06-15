package usecase

import (
	"context"
	"errors"

	"github.com/katedegree/spark/api/internal/domain/repository"
)

var ErrThingAlreadyExists = errors.New("thing already exists")

type ThingUsecase interface {
	Search(ctx context.Context, q string) ([]*repository.ThingRecord, error)
	Create(ctx context.Context, name string) (*repository.ThingRecord, error)
}

type thingUsecase struct {
	thingRepo repository.ThingRepository
}

func NewThingUsecase(thingRepo repository.ThingRepository) ThingUsecase {
	return &thingUsecase{thingRepo: thingRepo}
}

func (u *thingUsecase) Search(ctx context.Context, q string) ([]*repository.ThingRecord, error) {
	return u.thingRepo.Search(ctx, q)
}

func (u *thingUsecase) Create(ctx context.Context, name string) (*repository.ThingRecord, error) {
	existing, err := u.thingRepo.FindByName(ctx, name)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, ErrThingAlreadyExists
	}
	return u.thingRepo.Create(ctx, name)
}
