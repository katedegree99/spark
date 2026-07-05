package usecase

import (
	"context"

	"github.com/katedegree99/spark/api/internal/domain/repository"
)

type ListInterestsInput struct {
	Direction string // "sent" or "received"
	Offset    int
	Limit     int
}

type ListInterestsUsecase interface {
	ListInterests(ctx context.Context, userID uint, input ListInterestsInput) ([]*RecommendUserResult, error)
}

type listInterestsUsecase struct {
	interestRepo repository.InterestRepository
	profileRepo  repository.ProfileRepository
	thingRepo    repository.ThingRepository
	imageRepo    repository.ImageRepository
}

func NewListInterestsUsecase(
	interestRepo repository.InterestRepository,
	profileRepo repository.ProfileRepository,
	thingRepo repository.ThingRepository,
	imageRepo repository.ImageRepository,
) ListInterestsUsecase {
	return &listInterestsUsecase{
		interestRepo: interestRepo,
		profileRepo:  profileRepo,
		thingRepo:    thingRepo,
		imageRepo:    imageRepo,
	}
}

func (u *listInterestsUsecase) ListInterests(ctx context.Context, userID uint, input ListInterestsInput) ([]*RecommendUserResult, error) {
	// 自分のタグセットを取得
	doingIDs, err := u.profileRepo.FindDoingIDsByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	wantIDs, err := u.profileRepo.FindWantIDsByUserID(ctx, userID)
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

	var candidates []*repository.InterestRecord
	if input.Direction == "sent" {
		candidates, err = u.interestRepo.ListSent(ctx, userID, input.Offset, input.Limit)
	} else {
		candidates, err = u.interestRepo.ListReceived(ctx, userID, input.Offset, input.Limit)
	}
	if err != nil {
		return nil, err
	}

	results := make([]*RecommendUserResult, 0, len(candidates))
	for _, c := range candidates {
		var matched, unmatched []uint
		for _, id := range c.ThingIDs {
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
		if c.IconImageID != nil {
			img, err := u.imageRepo.FindByID(ctx, *c.IconImageID)
			if err != nil {
				return nil, err
			}
			if img != nil {
				iconURL = &img.URL
			}
		}

		results = append(results, &RecommendUserResult{
			UserID:          c.UserID,
			Name:            c.Name,
			Bio:             c.Bio,
			IconURL:         iconURL,
			CommonCount:     len(matched),
			MatchedThings:   matchedThings,
			UnmatchedThings: unmatchedThings,
		})
	}

	return results, nil
}
