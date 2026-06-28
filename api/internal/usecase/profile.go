package usecase

import (
	"context"
	"errors"

	"github.com/katedegree99/spark/api/internal/domain/repository"
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
	exists, err := u.profileRepo.ExistsByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrProfileAlreadyExists
	}

	p := &repository.ProfileRecord{
		UserID:      userID,
		Name:        input.Name,
		Bio:         input.Bio,
		IconImageID: input.IconImageID,
	}
	if err := u.profileRepo.Create(ctx, p); err != nil {
		return nil, err
	}

	if err := u.profileRepo.SetDoings(ctx, userID, input.DoingIDs); err != nil {
		return nil, err
	}
	if err := u.profileRepo.SetWants(ctx, userID, input.WantIDs); err != nil {
		return nil, err
	}

	return u.Get(ctx, userID)
}

func (u *profileUsecase) Get(ctx context.Context, userID uint) (*ProfileResult, error) {
	p, err := u.profileRepo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if p == nil {
		return nil, ErrProfileNotFound
	}

	doingIDs, err := u.profileRepo.FindDoingIDsByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	var doings []*repository.ThingRecord
	if len(doingIDs) > 0 {
		doings, err = u.thingRepo.FindByIDs(ctx, doingIDs)
		if err != nil {
			return nil, err
		}
	}

	wantIDs, err := u.profileRepo.FindWantIDsByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	var wants []*repository.ThingRecord
	if len(wantIDs) > 0 {
		wants, err = u.thingRepo.FindByIDs(ctx, wantIDs)
		if err != nil {
			return nil, err
		}
	}

	var image *repository.ImageRecord
	if p.IconImageID != nil {
		image, err = u.imageRepo.FindByID(ctx, *p.IconImageID)
		if err != nil {
			return nil, err
		}
	}

	return &ProfileResult{
		Profile: p,
		Doings:  doings,
		Wants:   wants,
		Image:   image,
	}, nil
}

func (u *profileUsecase) Update(ctx context.Context, userID uint, input ProfileInput) (*ProfileResult, error) {
	p, err := u.profileRepo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if p == nil {
		return nil, ErrProfileNotFound
	}

	p.Name = input.Name
	p.Bio = input.Bio
	p.IconImageID = input.IconImageID

	if err := u.profileRepo.Update(ctx, p); err != nil {
		return nil, err
	}
	if err := u.profileRepo.SetDoings(ctx, userID, input.DoingIDs); err != nil {
		return nil, err
	}
	if err := u.profileRepo.SetWants(ctx, userID, input.WantIDs); err != nil {
		return nil, err
	}

	return u.Get(ctx, userID)
}
