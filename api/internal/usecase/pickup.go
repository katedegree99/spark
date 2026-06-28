package usecase

import (
	"context"
	"sort"

	"github.com/katedegree/spark/api/internal/domain/repository"
)

const pickupLimit = 5

type PickupUserResult struct {
	UserID          uint
	Name            string
	Bio             *string
	IconURL         *string
	MatchedThings   []*repository.ThingRecord
	UnmatchedThings []*repository.ThingRecord
}

type PickupUsecase interface {
	ListPickup(ctx context.Context, userID uint) ([]*PickupUserResult, error)
}

type pickupUsecase struct {
	pickupRepo  repository.PickupRepository
	profileRepo repository.ProfileRepository
	thingRepo   repository.ThingRepository
	imageRepo   repository.ImageRepository
}

func NewPickupUsecase(
	pickupRepo repository.PickupRepository,
	profileRepo repository.ProfileRepository,
	thingRepo repository.ThingRepository,
	imageRepo repository.ImageRepository,
) PickupUsecase {
	return &pickupUsecase{
		pickupRepo:  pickupRepo,
		profileRepo: profileRepo,
		thingRepo:   thingRepo,
		imageRepo:   imageRepo,
	}
}

func (u *pickupUsecase) ListPickup(ctx context.Context, userID uint) ([]*PickupUserResult, error) {
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

	candidates, err := u.pickupRepo.FindCandidates(ctx, userID)
	if err != nil {
		return nil, err
	}

	type scoredCandidate struct {
		candidate *repository.PickupCandidateRecord
		matched   []uint
		unmatched []uint
	}

	scored := make([]scoredCandidate, 0, len(candidates))
	for _, c := range candidates {
		var matched, unmatched []uint
		for _, id := range c.ThingIDs {
			if _, ok := myThingSet[id]; ok {
				matched = append(matched, id)
			} else {
				unmatched = append(unmatched, id)
			}
		}
		scored = append(scored, scoredCandidate{c, matched, unmatched})
	}

	sort.Slice(scored, func(i, j int) bool {
		return len(scored[i].matched) > len(scored[j].matched)
	})
	if len(scored) > pickupLimit {
		scored = scored[:pickupLimit]
	}

	results := make([]*PickupUserResult, 0, len(scored))
	for _, s := range scored {
		var matchedThings []*repository.ThingRecord
		if len(s.matched) > 0 {
			matchedThings, err = u.thingRepo.FindByIDs(ctx, s.matched)
			if err != nil {
				return nil, err
			}
		}

		var unmatchedThings []*repository.ThingRecord
		if len(s.unmatched) > 0 {
			unmatchedThings, err = u.thingRepo.FindByIDs(ctx, s.unmatched)
			if err != nil {
				return nil, err
			}
		}

		var iconURL *string
		if s.candidate.IconImageID != nil {
			img, err := u.imageRepo.FindByID(ctx, *s.candidate.IconImageID)
			if err != nil {
				return nil, err
			}
			if img != nil {
				iconURL = &img.URL
			}
		}

		results = append(results, &PickupUserResult{
			UserID:          s.candidate.UserID,
			Name:            s.candidate.Name,
			Bio:             s.candidate.Bio,
			IconURL:         iconURL,
			MatchedThings:   matchedThings,
			UnmatchedThings: unmatchedThings,
		})
	}

	return results, nil
}
