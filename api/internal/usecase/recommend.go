package usecase

import (
	"context"
	"math/rand/v2"
	"sort"

	"github.com/katedegree99/spark/api/internal/domain/repository"
)

const (
	recommendLimit      = 12
	recommendPoolRatio  = 3 // 上位 recommendLimit×3 件からランダムに選出
)

type RecommendUserResult struct {
	UserID          uint
	Name            string
	Bio             *string
	IconURL         *string
	CommonCount     int
	MatchedThings   []*repository.ThingRecord
	UnmatchedThings []*repository.ThingRecord
}

type RecommendUsecase interface {
	ListRecommend(ctx context.Context, userID uint) ([]*RecommendUserResult, error)
}

type recommendUsecase struct {
	pickupRepo  repository.PickupRepository
	profileRepo repository.ProfileRepository
	thingRepo   repository.ThingRepository
	imageRepo   repository.ImageRepository
}

func NewRecommendUsecase(
	pickupRepo repository.PickupRepository,
	profileRepo repository.ProfileRepository,
	thingRepo repository.ThingRepository,
	imageRepo repository.ImageRepository,
) RecommendUsecase {
	return &recommendUsecase{
		pickupRepo:  pickupRepo,
		profileRepo: profileRepo,
		thingRepo:   thingRepo,
		imageRepo:   imageRepo,
	}
}

func (u *recommendUsecase) ListRecommend(ctx context.Context, userID uint) ([]*RecommendUserResult, error) {
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

	// タグID だけでスコアリング（DB クエリなし）
	type scored struct {
		candidate *repository.PickupCandidateRecord
		matched   []uint
		unmatched []uint
	}

	all := make([]scored, 0, len(candidates))
	for _, c := range candidates {
		var matched, unmatched []uint
		for _, id := range c.ThingIDs {
			if _, ok := myThingSet[id]; ok {
				matched = append(matched, id)
			} else {
				unmatched = append(unmatched, id)
			}
		}
		all = append(all, scored{c, matched, unmatched})
	}

	// 共通タグ数で降順ソート
	sort.Slice(all, func(i, j int) bool {
		return len(all[i].matched) > len(all[j].matched)
	})

	// 上位 pool 件に絞ってシャッフル → recommendLimit 件選出
	pool := recommendLimit * recommendPoolRatio
	if len(all) > pool {
		all = all[:pool]
	}
	rand.Shuffle(len(all), func(i, j int) { all[i], all[j] = all[j], all[i] })
	if len(all) > recommendLimit {
		all = all[:recommendLimit]
	}

	// 選出された件数分だけタグ名・画像を取得
	results := make([]*RecommendUserResult, 0, len(all))
	for _, s := range all {
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

		results = append(results, &RecommendUserResult{
			UserID:          s.candidate.UserID,
			Name:            s.candidate.Name,
			Bio:             s.candidate.Bio,
			IconURL:         iconURL,
			CommonCount:     len(s.matched),
			MatchedThings:   matchedThings,
			UnmatchedThings: unmatchedThings,
		})
	}

	return results, nil
}
