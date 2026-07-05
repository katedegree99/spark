package repository

import (
	"context"
	"errors"

	"github.com/katedegree99/spark/api/internal/domain/model"
	domainrepo "github.com/katedegree99/spark/api/internal/domain/repository"
	"gorm.io/gorm"
)

type interestRepository struct {
	db *gorm.DB
}

func NewInterestRepository(db *gorm.DB) domainrepo.InterestRepository {
	return &interestRepository{db: db}
}

func (r *interestRepository) Exists(ctx context.Context, fromUserID, toUserID uint) (bool, error) {
	var row model.UserInterest
	err := r.db.WithContext(ctx).
		Where("from_user_id = ? AND to_user_id = ?", fromUserID, toUserID).
		First(&row).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (r *interestRepository) Create(ctx context.Context, fromUserID, toUserID uint) error {
	return r.db.WithContext(ctx).Create(&model.UserInterest{
		FromUserID: fromUserID,
		ToUserID:   toUserID,
	}).Error
}

func (r *interestRepository) ListSent(ctx context.Context, fromUserID uint, offset, limit int) ([]*domainrepo.InterestRecord, error) {
	var interests []model.UserInterest
	if err := r.db.WithContext(ctx).
		Where("from_user_id = ?", fromUserID).
		Order("created_at DESC").
		Offset(offset).Limit(limit).
		Find(&interests).Error; err != nil {
		return nil, err
	}
	toUserIDs := make([]uint, len(interests))
	for i, v := range interests {
		toUserIDs[i] = v.ToUserID
	}
	return r.fetchProfileRecords(ctx, toUserIDs)
}

func (r *interestRepository) ListReceived(ctx context.Context, toUserID uint, offset, limit int) ([]*domainrepo.InterestRecord, error) {
	var interests []model.UserInterest
	if err := r.db.WithContext(ctx).
		Where("to_user_id = ?", toUserID).
		Order("created_at DESC").
		Offset(offset).Limit(limit).
		Find(&interests).Error; err != nil {
		return nil, err
	}
	fromUserIDs := make([]uint, len(interests))
	for i, v := range interests {
		fromUserIDs[i] = v.FromUserID
	}
	return r.fetchProfileRecords(ctx, fromUserIDs)
}

func (r *interestRepository) fetchProfileRecords(ctx context.Context, userIDs []uint) ([]*domainrepo.InterestRecord, error) {
	if len(userIDs) == 0 {
		return nil, nil
	}

	var profiles []model.Profile
	if err := r.db.WithContext(ctx).Where("user_id IN ?", userIDs).Find(&profiles).Error; err != nil {
		return nil, err
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

	// userIDs の順序を保持して返す
	profileByUserID := make(map[uint]model.Profile, len(profiles))
	for _, p := range profiles {
		profileByUserID[p.UserID] = p
	}

	records := make([]*domainrepo.InterestRecord, 0, len(userIDs))
	for _, uid := range userIDs {
		p, ok := profileByUserID[uid]
		if !ok {
			continue
		}
		thingSet := thingIDsByUser[p.UserID]
		thingIDs := make([]uint, 0, len(thingSet))
		for id := range thingSet {
			thingIDs = append(thingIDs, id)
		}
		records = append(records, &domainrepo.InterestRecord{
			UserID:      p.UserID,
			Name:        p.Name,
			Bio:         p.Bio,
			IconImageID: p.IconImageID,
			ThingIDs:    thingIDs,
		})
	}
	return records, nil
}
