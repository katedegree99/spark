package usecase

import (
	"context"
	"testing"
	"time"

	"github.com/katedegree/spark/api/internal/domain/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// --- mocks for pickup ---

type mockPickupRepository struct {
	findCandidates func(ctx context.Context, excludeUserID uint) ([]*repository.PickupCandidateRecord, error)
	findCache      func(ctx context.Context, userID uint, date time.Time) ([]uint, error)
	saveCache      func(ctx context.Context, userID uint, date time.Time, pickedUserIDs []uint) error
}

func (m *mockPickupRepository) FindCandidates(ctx context.Context, excludeUserID uint) ([]*repository.PickupCandidateRecord, error) {
	if m.findCandidates != nil {
		return m.findCandidates(ctx, excludeUserID)
	}
	return nil, nil
}

func (m *mockPickupRepository) FindCache(ctx context.Context, userID uint, date time.Time) ([]uint, error) {
	if m.findCache != nil {
		return m.findCache(ctx, userID, date)
	}
	return nil, nil
}

func (m *mockPickupRepository) SaveCache(ctx context.Context, userID uint, date time.Time, pickedUserIDs []uint) error {
	if m.saveCache != nil {
		return m.saveCache(ctx, userID, date, pickedUserIDs)
	}
	return nil
}

type mockThingRepository struct {
	findByIDs func(ctx context.Context, ids []uint) ([]*repository.ThingRecord, error)
}

func (m *mockThingRepository) Search(_ context.Context, _ string) ([]*repository.ThingRecord, error) {
	return nil, nil
}
func (m *mockThingRepository) FindByID(_ context.Context, _ uint) (*repository.ThingRecord, error) {
	return nil, nil
}
func (m *mockThingRepository) FindByName(_ context.Context, _ string) (*repository.ThingRecord, error) {
	return nil, nil
}
func (m *mockThingRepository) FindByIDs(ctx context.Context, ids []uint) ([]*repository.ThingRecord, error) {
	if m.findByIDs != nil {
		return m.findByIDs(ctx, ids)
	}
	return nil, nil
}
func (m *mockThingRepository) Create(_ context.Context, _ string) (*repository.ThingRecord, error) {
	return nil, nil
}
func (m *mockThingRepository) FindAliasesByThingID(_ context.Context, _ uint) ([]string, error) {
	return nil, nil
}
func (m *mockThingRepository) HasAlias(_ context.Context, _ uint) (bool, error) {
	return false, nil
}
func (m *mockThingRepository) AddAlias(_ context.Context, _ uint, _ string) error { return nil }
func (m *mockThingRepository) UpdateName(_ context.Context, _ uint, _ string) error {
	return nil
}
func (m *mockThingRepository) FindAllWithAliases(_ context.Context) ([]*repository.ThingRecord, error) {
	return nil, nil
}

type mockImageRepository struct {
	findByID func(ctx context.Context, id uint) (*repository.ImageRecord, error)
}

func (m *mockImageRepository) Create(_ context.Context, _, _ string) (*repository.ImageRecord, error) {
	return nil, nil
}
func (m *mockImageRepository) FindByID(ctx context.Context, id uint) (*repository.ImageRecord, error) {
	if m.findByID != nil {
		return m.findByID(ctx, id)
	}
	return nil, nil
}

// --- helpers ---

func thing(id uint, name string) *repository.ThingRecord {
	return &repository.ThingRecord{ID: id, Name: name, Aliases: []string{}, CreatedAt: time.Time{}}
}

func newPickupUsecase(
	pickupRepo repository.PickupRepository,
	profileRepo repository.ProfileRepository,
	thingRepo repository.ThingRepository,
	imageRepo repository.ImageRepository,
) PickupUsecase {
	return NewPickupUsecase(pickupRepo, profileRepo, thingRepo, imageRepo)
}

// --- tests ---

func TestListPickup_Success(t *testing.T) {
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

	uc := newPickupUsecase(pickupRepo, profileRepo, thingRepo, &mockImageRepository{})
	results, err := uc.ListPickup(context.Background(), 1)

	require.NoError(t, err)
	require.Len(t, results, 2)

	// Alice matches IDs 1,2 → 2 matched; Bob matches none → 0 matched
	assert.Equal(t, uint(10), results[0].UserID)
	assert.Len(t, results[0].MatchedThings, 2)
	assert.Len(t, results[0].UnmatchedThings, 1)

	assert.Equal(t, uint(11), results[1].UserID)
	assert.Len(t, results[1].MatchedThings, 0)
	assert.Len(t, results[1].UnmatchedThings, 2)
}

func TestListPickup_SortsByMatchedCount(t *testing.T) {
	profileRepo := &mockProfileRepository{
		findDoingIDsByUserID: func(_ context.Context, _ uint) ([]uint, error) {
			return []uint{1, 2, 3}, nil
		},
		findWantIDsByUserID: func(_ context.Context, _ uint) ([]uint, error) {
			return nil, nil
		},
	}
	pickupRepo := &mockPickupRepository{
		findCandidates: func(_ context.Context, _ uint) ([]*repository.PickupCandidateRecord, error) {
			return []*repository.PickupCandidateRecord{
				{UserID: 20, Name: "Low", ThingIDs: []uint{1}},
				{UserID: 21, Name: "High", ThingIDs: []uint{1, 2, 3}},
				{UserID: 22, Name: "Mid", ThingIDs: []uint{1, 2}},
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

	uc := newPickupUsecase(pickupRepo, profileRepo, thingRepo, &mockImageRepository{})
	results, err := uc.ListPickup(context.Background(), 1)

	require.NoError(t, err)
	require.Len(t, results, 3)
	assert.Equal(t, uint(21), results[0].UserID) // High: 3 matches
	assert.Equal(t, uint(22), results[1].UserID) // Mid:  2 matches
	assert.Equal(t, uint(20), results[2].UserID) // Low:  1 match
}

func TestListPickup_LimitsToFive(t *testing.T) {
	profileRepo := &mockProfileRepository{
		findDoingIDsByUserID: func(_ context.Context, _ uint) ([]uint, error) {
			return []uint{1}, nil
		},
		findWantIDsByUserID: func(_ context.Context, _ uint) ([]uint, error) {
			return nil, nil
		},
	}
	candidates := make([]*repository.PickupCandidateRecord, 8)
	for i := range candidates {
		candidates[i] = &repository.PickupCandidateRecord{
			UserID:   uint(100 + i),
			Name:     "User",
			ThingIDs: []uint{1},
		}
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

	uc := newPickupUsecase(pickupRepo, profileRepo, thingRepo, &mockImageRepository{})
	results, err := uc.ListPickup(context.Background(), 1)

	require.NoError(t, err)
	assert.Len(t, results, pickupLimit)
}

func TestListPickup_NoCandidates(t *testing.T) {
	profileRepo := &mockProfileRepository{
		findDoingIDsByUserID: func(_ context.Context, _ uint) ([]uint, error) {
			return []uint{1}, nil
		},
		findWantIDsByUserID: func(_ context.Context, _ uint) ([]uint, error) {
			return nil, nil
		},
	}
	pickupRepo := &mockPickupRepository{
		findCandidates: func(_ context.Context, _ uint) ([]*repository.PickupCandidateRecord, error) {
			return nil, nil
		},
	}

	uc := newPickupUsecase(pickupRepo, profileRepo, &mockThingRepository{}, &mockImageRepository{})
	results, err := uc.ListPickup(context.Background(), 1)

	require.NoError(t, err)
	assert.Empty(t, results)
}

func TestListPickup_WithIcon(t *testing.T) {
	iconID := uint(99)
	iconURL := "https://example.com/icon.jpg"
	profileRepo := &mockProfileRepository{
		findDoingIDsByUserID: func(_ context.Context, _ uint) ([]uint, error) {
			return nil, nil
		},
		findWantIDsByUserID: func(_ context.Context, _ uint) ([]uint, error) {
			return nil, nil
		},
	}
	pickupRepo := &mockPickupRepository{
		findCandidates: func(_ context.Context, _ uint) ([]*repository.PickupCandidateRecord, error) {
			return []*repository.PickupCandidateRecord{
				{UserID: 10, Name: "Alice", IconImageID: &iconID, ThingIDs: nil},
			}, nil
		},
	}
	imageRepo := &mockImageRepository{
		findByID: func(_ context.Context, id uint) (*repository.ImageRecord, error) {
			return &repository.ImageRecord{ID: id, URL: iconURL}, nil
		},
	}

	uc := newPickupUsecase(pickupRepo, profileRepo, &mockThingRepository{}, imageRepo)
	results, err := uc.ListPickup(context.Background(), 1)

	require.NoError(t, err)
	require.Len(t, results, 1)
	require.NotNil(t, results[0].IconURL)
	assert.Equal(t, iconURL, *results[0].IconURL)
}
