package usecase

import (
	"context"
	"errors"

	"github.com/katedegree/spark/api/internal/domain/repository"
)

var ErrProfileAlreadyExists = errors.New("profile already exists")
var ErrProfileNotFound = errors.New("profile not found")

type ProfileInput struct {
	Name        string
	Bio         *string
	IconImageID *uint
	DoingIDs    []uint
	WantIDs     []uint
}

type ProfileResult struct {
	Profile *repository.ProfileRecord
	Doings  []*repository.ThingRecord
	Wants   []*repository.ThingRecord
	Image   *repository.ImageRecord
}

type ProfileUsecase interface {
	Create(ctx context.Context, userID uint, input ProfileInput) (*ProfileResult, error)
	Get(ctx context.Context, userID uint) (*ProfileResult, error)
	Update(ctx context.Context, userID uint, input ProfileInput) (*ProfileResult, error)
}

type profileUsecase struct {
	profileRepo repository.ProfileRepository
	thingRepo   repository.ThingRepository
	imageRepo   repository.ImageRepository
}

func NewProfileUsecase(
	profileRepo repository.ProfileRepository,
	thingRepo repository.ThingRepository,
	imageRepo repository.ImageRepository,
) ProfileUsecase {
	return &profileUsecase{
		profileRepo: profileRepo,
		thingRepo:   thingRepo,
		imageRepo:   imageRepo,
	}
}

func (u *profileUsecase) Create(ctx context.Context, userID uint, input ProfileInput) (*ProfileResult, error) {
	panic("not implemented")
}

func (u *profileUsecase) Get(ctx context.Context, userID uint) (*ProfileResult, error) {
	panic("not implemented")
}

func (u *profileUsecase) Update(ctx context.Context, userID uint, input ProfileInput) (*ProfileResult, error) {
	panic("not implemented")
}
