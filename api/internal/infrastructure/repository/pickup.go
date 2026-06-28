package repository

import (
	"context"

	"github.com/katedegree/spark/api/internal/domain/model"
	domainrepo "github.com/katedegree/spark/api/internal/domain/repository"
	"gorm.io/gorm"
)

type pickupRepository struct {
	db *gorm.DB
}

func NewPickupRepository(db *gorm.DB) domainrepo.PickupRepository {
	return &pickupRepository{db: db}
}

func (r *pickupRepository) FindCandidates(ctx context.Context, excludeUserID uint) ([]*domainrepo.PickupCandidateRecord, error) {
	var profiles []model.Profile
	if err := r.db.WithContext(ctx).Where("user_id != ?", excludeUserID).Find(&profiles).Error; err != nil {
		return nil, err
	}
	if len(profiles) == 0 {
		return nil, nil
	}

	userIDs := make([]uint, len(profiles))
	for i, p := range profiles {
		userIDs[i] = p.UserID
	}

	var doings []model.UserDoing
	if err := r.db.WithContext(ctx).Where("user_id IN ?", userIDs).Find(&doings).Error; err != nil {
		return nil, err
	}

	var wants []model.UserWant
	if err := r.db.WithContext(ctx).Where("user_id IN ?", userIDs).Find(&wants).Error; err != nil {
		return nil, err
	}

	thingIDsByUser := make(map[uint]map[uint]struct{}, len(profiles))
	for _, d := range doings {
		if _, ok := thingIDsByUser[d.UserID]; !ok {
			thingIDsByUser[d.UserID] = map[uint]struct{}{}
		}
		thingIDsByUser[d.UserID][d.ThingID] = struct{}{}
	}
	for _, w := range wants {
		if _, ok := thingIDsByUser[w.UserID]; !ok {
			thingIDsByUser[w.UserID] = map[uint]struct{}{}
		}
		thingIDsByUser[w.UserID][w.ThingID] = struct{}{}
	}

	records := make([]*domainrepo.PickupCandidateRecord, 0, len(profiles))
	for _, p := range profiles {
		thingSet := thingIDsByUser[p.UserID]
		thingIDs := make([]uint, 0, len(thingSet))
		for id := range thingSet {
			thingIDs = append(thingIDs, id)
		}
		records = append(records, &domainrepo.PickupCandidateRecord{
			UserID:      p.UserID,
			Name:        p.Name,
			Bio:         p.Bio,
			IconImageID: p.IconImageID,
			ThingIDs:    thingIDs,
		})
	}
	return records, nil
}
