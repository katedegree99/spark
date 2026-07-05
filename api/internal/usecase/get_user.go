package usecase

import (
	"context"
	"errors"

	"github.com/katedegree99/spark/api/internal/domain/repository"
)

var ErrUserNotFound = errors.New("user not found")

type GetUserResult struct {
	UserID          uint
	Name            string
	Bio             *string
	IconURL         *string
	CommonCount     int
	MatchedThings   []*repository.ThingRecord
	UnmatchedThings []*repository.ThingRecord
	IsInterested    bool
}

type GetUserUsecase interface {
	GetUser(ctx context.Context, myUserID, targetUserID uint) (*GetUserResult, error)
}

type getUserUsecase struct {
	userRepo     repository.UserRepository
	profileRepo  repository.ProfileRepository
	thingRepo    repository.ThingRepository
	imageRepo    repository.ImageRepository
	interestRepo repository.InterestRepository
}

func NewGetUserUsecase(
	userRepo repository.UserRepository,
	profileRepo repository.ProfileRepository,
	thingRepo repository.ThingRepository,
	imageRepo repository.ImageRepository,
	interestRepo repository.InterestRepository,
) GetUserUsecase {
	return &getUserUsecase{
		userRepo:     userRepo,
		profileRepo:  profileRepo,
		thingRepo:    thingRepo,
		imageRepo:    imageRepo,
		interestRepo: interestRepo,
	}
}

func (u *getUserUsecase) GetUser(ctx context.Context, myUserID, targetUserID uint) (*GetUserResult, error) {
	profile, err := u.profileRepo.FindByUserID(ctx, targetUserID)
	if err != nil {
		return nil, err
	}
	if profile == nil {
		return nil, ErrUserNotFound
	}

	// 自分のタグセットを取得
	doingIDs, err := u.profileRepo.FindDoingIDsByUserID(ctx, myUserID)
	if err != nil {
		return nil, err
	}
	wantIDs, err := u.profileRepo.FindWantIDsByUserID(ctx, myUserID)
	if err != nil {
		return nil, err
	}
	myThingSet := make(map[uint]struct{}, len(doingIDs)+len(wantIDs))
	for _, id := range doingIDs {
		myThingSet[id] = struct{}{}
	}
	for _, id := range wantIDs {
		myThingSet[id] = struct{}{}
	}

	// 対象ユーザーのタグを取得
	targetDoingIDs, err := u.profileRepo.FindDoingIDsByUserID(ctx, targetUserID)
	if err != nil {
		return nil, err
	}
	targetWantIDs, err := u.profileRepo.FindWantIDsByUserID(ctx, targetUserID)
	if err != nil {
		return nil, err
	}
	targetThingSet := make(map[uint]struct{}, len(targetDoingIDs)+len(targetWantIDs))
	for _, id := range targetDoingIDs {
		targetThingSet[id] = struct{}{}
	}
	for _, id := range targetWantIDs {
		targetThingSet[id] = struct{}{}
	}

	var matched, unmatched []uint
	for id := range targetThingSet {
		if _, ok := myThingSet[id]; ok {
			matched = append(matched, id)
		} else {
			unmatched = append(unmatched, id)
		}
	}

	var matchedThings []*repository.ThingRecord
	if len(matched) > 0 {
		matchedThings, err = u.thingRepo.FindByIDs(ctx, matched)
		if err != nil {
			return nil, err
		}
	}
	var unmatchedThings []*repository.ThingRecord
	if len(unmatched) > 0 {
		unmatchedThings, err = u.thingRepo.FindByIDs(ctx, unmatched)
		if err != nil {
			return nil, err
		}
	}

	var iconURL *string
	if profile.IconImageID != nil {
		img, err := u.imageRepo.FindByID(ctx, *profile.IconImageID)
		if err != nil {
			return nil, err
		}
		if img != nil {
			iconURL = &img.URL
		}
	}

	isInterested, err := u.interestRepo.Exists(ctx, myUserID, targetUserID)
	if err != nil {
		return nil, err
	}

	return &GetUserResult{
		UserID:          targetUserID,
		Name:            profile.Name,
		Bio:             profile.Bio,
		IconURL:         iconURL,
		CommonCount:     len(matched),
		MatchedThings:   matchedThings,
		UnmatchedThings: unmatchedThings,
		IsInterested:    isInterested,
	}, nil
}
