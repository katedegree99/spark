package repository

import (
	"context"
	"errors"
	"time"

	"github.com/katedegree99/spark/api/internal/domain/model"
	domainrepo "github.com/katedegree99/spark/api/internal/domain/repository"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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

func (r *pickupRepository) FindCache(ctx context.Context, userID uint, date time.Time) ([]uint, error) {
	var row model.UserPickupCache
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND cache_date = ?", userID, date.Format("2006-01-02")).
		First(&row).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return []uint(row.PickedUserIDs), nil
}

func (r *pickupRepository) SaveCache(ctx context.Context, userID uint, date time.Time, pickedUserIDs []uint) error {
	row := model.UserPickupCache{
		UserID:        userID,
		CacheDate:     date,
		PickedUserIDs: model.UintSlice(pickedUserIDs),
	}
	return r.db.WithContext(ctx).
		Where("user_id = ? AND cache_date = ?", userID, date.Format("2006-01-02")).
		Assign(model.UserPickupCache{PickedUserIDs: row.PickedUserIDs}).
		FirstOrCreate(&row).Error
}

func (r *pickupRepository) FindShownUserIDs(ctx context.Context, userID uint) ([]uint, error) {
	var rows []model.UserPickupHistory
	if err := r.db.WithContext(ctx).
		Select("shown_user_id").
		Where("user_id = ?", userID).
		Find(&rows).Error; err != nil {
		return nil, err
	}
	ids := make([]uint, 0, len(rows))
	for _, row := range rows {
		ids = append(ids, row.ShownUserID)
	}
	return ids, nil
}

func (r *pickupRepository) SaveHistory(ctx context.Context, userID uint, shownUserIDs []uint) error {
	if len(shownUserIDs) == 0 {
		return nil
	}
	rows := make([]model.UserPickupHistory, 0, len(shownUserIDs))
	for _, shownID := range shownUserIDs {
		rows = append(rows, model.UserPickupHistory{
			UserID:      userID,
			ShownUserID: shownID,
		})
	}
	// 重複（uk_user_shown 違反）は無視して挿入
	return r.db.WithContext(ctx).
		Clauses(clause.OnConflict{DoNothing: true}).
		Create(&rows).Error
}
