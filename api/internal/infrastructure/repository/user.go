package repository

import (
	"context"

	"github.com/katedegree99/spark/api/internal/domain/model"
	domainrepo "github.com/katedegree99/spark/api/internal/domain/repository"
	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) domainrepo.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) ListUsers(ctx context.Context, excludeUserID uint, params domainrepo.ListUsersParams) ([]*domainrepo.UserListRecord, error) {
	db := r.db.WithContext(ctx).Model(&model.Profile{}).Where("user_id != ?", excludeUserID)

	if params.Q != nil && *params.Q != "" {
		db = db.Where("name LIKE ?", *params.Q+"%")
	}

	if len(params.TagIDs) > 0 {
		// AND 条件: 指定されたタグを全て持つユーザーのみ
		sub := r.db.WithContext(ctx).Raw(
			`SELECT user_id FROM (
				SELECT user_id, thing_id FROM user_doings WHERE thing_id IN ?
				UNION
				SELECT user_id, thing_id FROM user_wants WHERE thing_id IN ?
			) combined GROUP BY user_id HAVING COUNT(DISTINCT thing_id) = ?`,
			params.TagIDs, params.TagIDs, len(params.TagIDs),
		)
		db = db.Where("user_id IN (?)", sub)
	}

	var profiles []model.Profile
	if err := db.Offset(params.Offset).Limit(params.Limit).Find(&profiles).Error; err != nil {
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

	records := make([]*domainrepo.UserListRecord, 0, len(profiles))
	for _, p := range profiles {
		thingSet := thingIDsByUser[p.UserID]
		thingIDs := make([]uint, 0, len(thingSet))
		for id := range thingSet {
			thingIDs = append(thingIDs, id)
		}
		records = append(records, &domainrepo.UserListRecord{
			UserID:      p.UserID,
			Name:        p.Name,
			Bio:         p.Bio,
			IconImageID: p.IconImageID,
			ThingIDs:    thingIDs,
		})
	}

	return records, nil
}
