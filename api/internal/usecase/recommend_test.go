package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/katedegree/spark/api/internal/domain/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newRecommendUsecase(
	pickupRepo repository.PickupRepository,
	profileRepo repository.ProfileRepository,
	thingRepo repository.ThingRepository,
	imageRepo repository.ImageRepository,
) RecommendUsecase {
	return NewRecommendUsecase(pickupRepo, profileRepo, thingRepo, imageRepo)
}

func TestListRecommend_Success(t *testing.T) {
	profileRepo := &mockProfileRepository{
		findDoingIDsByUserID: func(_ context.Context, _ uint) ([]uint, error) {
			return []uint{1, 2}, nil
		},
		findWantIDsByUserID: func(_ context.Context, _ uint) ([]uint, error) {
			return []uint{3}, nil
		},
	}
	pickupRepo := &mockPickupRepository{
		findCandidates: func(_ context.Context, _ uint) ([]*repository.PickupCandidateRecord, error) {
			return []*repository.PickupCandidateRecord{
				{UserID: 10, Name: "Alice", ThingIDs: []uint{1, 2, 4}},
				{UserID: 11, Name: "Bob", ThingIDs: []uint{5, 6}},
			}, nil
		},
	}
	thingRepo := &mockThingRepository{
		findByIDs: func(_ context.Context, ids []uint) ([]*repository.ThingRecord, error) {
			results := make([]*repository.ThingRecord, 0, len(ids))
			for _, id := range ids {
				results = append(results, thing(id, "thing"))
			}
			return results, nil
		},
	}

	uc := newRecommendUsecase(pickupRepo, profileRepo, thingRepo, &mockImageRepository{})
	results, err := uc.ListRecommend(context.Background(), 1)

	require.NoError(t, err)
	require.Len(t, results, 2)

	byID := make(map[uint]*RecommendUserResult, len(results))
	for _, r := range results {
		byID[r.UserID] = r
	}

	alice, ok := byID[10]
	require.True(t, ok, "Alice が含まれていない")
	assert.Equal(t, 2, alice.CommonCount)
	assert.Len(t, alice.MatchedThings, 2)
	assert.Len(t, alice.UnmatchedThings, 1)

	bob, ok := byID[11]
	require.True(t, ok, "Bob が含まれていない")
	assert.Equal(t, 0, bob.CommonCount)
	assert.Empty(t, bob.MatchedThings)
	assert.Len(t, bob.UnmatchedThings, 2)
}

func TestListRecommend_NoCandidates(t *testing.T) {
	profileRepo := &mockProfileRepository{
		findDoingIDsByUserID: func(_ context.Context, _ uint) ([]uint, error) { return []uint{1}, nil },
		findWantIDsByUserID:  func(_ context.Context, _ uint) ([]uint, error) { return nil, nil },
	}
	pickupRepo := &mockPickupRepository{
		findCandidates: func(_ context.Context, _ uint) ([]*repository.PickupCandidateRecord, error) {
			return nil, nil
		},
	}

	uc := newRecommendUsecase(pickupRepo, profileRepo, &mockThingRepository{}, &mockImageRepository{})
	results, err := uc.ListRecommend(context.Background(), 1)

	require.NoError(t, err)
	assert.Empty(t, results)
}

func TestListRecommend_LimitsToRecommendLimit(t *testing.T) {
	profileRepo := &mockProfileRepository{
		findDoingIDsByUserID: func(_ context.Context, _ uint) ([]uint, error) { return []uint{1}, nil },
		findWantIDsByUserID:  func(_ context.Context, _ uint) ([]uint, error) { return nil, nil },
	}
	// pool(36) を超える候補を用意
	candidates := make([]*repository.PickupCandidateRecord, 50)
	for i := range candidates {
		candidates[i] = &repository.PickupCandidateRecord{UserID: uint(100 + i), Name: "User", ThingIDs: []uint{1}}
	}
	pickupRepo := &mockPickupRepository{
		findCandidates: func(_ context.Context, _ uint) ([]*repository.PickupCandidateRecord, error) {
			return candidates, nil
		},
	}
	thingRepo := &mockThingRepository{
		findByIDs: func(_ context.Context, ids []uint) ([]*repository.ThingRecord, error) {
			results := make([]*repository.ThingRecord, 0, len(ids))
			for _, id := range ids {
				results = append(results, thing(id, "thing"))
			}
			return results, nil
		},
	}

	uc := newRecommendUsecase(pickupRepo, profileRepo, thingRepo, &mockImageRepository{})
	results, err := uc.ListRecommend(context.Background(), 1)

	require.NoError(t, err)
	assert.Len(t, results, recommendLimit)
}

func TestListRecommend_WithIcon(t *testing.T) {
	iconID := uint(99)
	iconURL := "https://example.com/icon.jpg"
	profileRepo := &mockProfileRepository{
		findDoingIDsByUserID: func(_ context.Context, _ uint) ([]uint, error) { return nil, nil },
		findWantIDsByUserID:  func(_ context.Context, _ uint) ([]uint, error) { return nil, nil },
	}
	pickupRepo := &mockPickupRepository{
		findCandidates: func(_ context.Context, _ uint) ([]*repository.PickupCandidateRecord, error) {
			return []*repository.PickupCandidateRecord{
				{UserID: 10, Name: "Alice", IconImageID: &iconID},
			}, nil
		},
	}
	imageRepo := &mockImageRepository{
		findByID: func(_ context.Context, id uint) (*repository.ImageRecord, error) {
			return &repository.ImageRecord{ID: id, URL: iconURL}, nil
		},
	}

	uc := newRecommendUsecase(pickupRepo, profileRepo, &mockThingRepository{}, imageRepo)
	results, err := uc.ListRecommend(context.Background(), 1)

	require.NoError(t, err)
	require.Len(t, results, 1)
	require.NotNil(t, results[0].IconURL)
	assert.Equal(t, iconURL, *results[0].IconURL)
}

func TestListRecommend_CommonCountIsMatchedLen(t *testing.T) {
	profileRepo := &mockProfileRepository{
		findDoingIDsByUserID: func(_ context.Context, _ uint) ([]uint, error) { return []uint{1, 2, 3}, nil },
		findWantIDsByUserID:  func(_ context.Context, _ uint) ([]uint, error) { return nil, nil },
	}
	pickupRepo := &mockPickupRepository{
		findCandidates: func(_ context.Context, _ uint) ([]*repository.PickupCandidateRecord, error) {
			return []*repository.PickupCandidateRecord{
				{UserID: 10, Name: "Alice", ThingIDs: []uint{1, 2, 3, 4, 5}},
			}, nil
		},
	}
	thingRepo := &mockThingRepository{
		findByIDs: func(_ context.Context, ids []uint) ([]*repository.ThingRecord, error) {
			results := make([]*repository.ThingRecord, 0, len(ids))
			for _, id := range ids {
				results = append(results, thing(id, "thing"))
			}
			return results, nil
		},
	}

	uc := newRecommendUsecase(pickupRepo, profileRepo, thingRepo, &mockImageRepository{})
	results, err := uc.ListRecommend(context.Background(), 1)

	require.NoError(t, err)
	require.Len(t, results, 1)
	assert.Equal(t, 3, results[0].CommonCount)
}

func TestListRecommend_RepositoryError(t *testing.T) {
	profileRepo := &mockProfileRepository{
		findDoingIDsByUserID: func(_ context.Context, _ uint) ([]uint, error) {
			return nil, errors.New("db error")
		},
	}

	uc := newRecommendUsecase(&mockPickupRepository{}, profileRepo, &mockThingRepository{}, &mockImageRepository{})
	_, err := uc.ListRecommend(context.Background(), 1)

	require.Error(t, err)
}
