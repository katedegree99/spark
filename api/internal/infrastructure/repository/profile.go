package repository

import (
	"context"
	"errors"

	"github.com/katedegree99/spark/api/internal/domain/model"
	domainrepo "github.com/katedegree99/spark/api/internal/domain/repository"
	"gorm.io/gorm"
)

type profileRepository struct {
	db *gorm.DB
}

func NewProfileRepository(db *gorm.DB) domainrepo.ProfileRepository {
	return &profileRepository{db: db}
}

func (r *profileRepository) ExistsByUserID(ctx context.Context, userID uint) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&model.Profile{}).Where("user_id = ?", userID).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *profileRepository) Create(ctx context.Context, p *domainrepo.ProfileRecord) error {
	m := model.Profile{
		UserID:      p.UserID,
		Name:        p.Name,
		Bio:         p.Bio,
		IconImageID: p.IconImageID,
	}
	if err := r.db.WithContext(ctx).Create(&m).Error; err != nil {
		return err
	}
	p.CreatedAt = m.CreatedAt
	p.UpdatedAt = m.UpdatedAt
	return nil
}

func (r *profileRepository) FindByUserID(ctx context.Context, userID uint) (*domainrepo.ProfileRecord, error) {
	var m model.Profile
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&m).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &domainrepo.ProfileRecord{
		UserID:      m.UserID,
		Name:        m.Name,
		Bio:         m.Bio,
		IconImageID: m.IconImageID,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}, nil
}

func (r *profileRepository) Update(ctx context.Context, p *domainrepo.ProfileRecord) error {
	result := r.db.WithContext(ctx).Model(&model.Profile{}).Where("user_id = ?", p.UserID).Updates(map[string]any{
		"name":          p.Name,
		"bio":           p.Bio,
		"icon_image_id": p.IconImageID,
	})
	if result.Error != nil {
		return result.Error
	}
	// Reload to get updated timestamps
	var m model.Profile
	if err := r.db.WithContext(ctx).Where("user_id = ?", p.UserID).First(&m).Error; err != nil {
		return err
	}
	p.CreatedAt = m.CreatedAt
	p.UpdatedAt = m.UpdatedAt
	return nil
}

func (r *profileRepository) SetDoings(ctx context.Context, userID uint, thingIDs []uint) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("user_id = ?", userID).Delete(&model.UserDoing{}).Error; err != nil {
			return err
		}
		if len(thingIDs) == 0 {
			return nil
		}
		rows := make([]model.UserDoing, 0, len(thingIDs))
		for _, id := range thingIDs {
			rows = append(rows, model.UserDoing{UserID: userID, ThingID: id})
		}
		return tx.Create(&rows).Error
	})
}

func (r *profileRepository) SetWants(ctx context.Context, userID uint, thingIDs []uint) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("user_id = ?", userID).Delete(&model.UserWant{}).Error; err != nil {
			return err
		}
		if len(thingIDs) == 0 {
			return nil
		}
		rows := make([]model.UserWant, 0, len(thingIDs))
		for _, id := range thingIDs {
			rows = append(rows, model.UserWant{UserID: userID, ThingID: id})
		}
		return tx.Create(&rows).Error
	})
}

func (r *profileRepository) FindDoingIDsByUserID(ctx context.Context, userID uint) ([]uint, error) {
	var rows []model.UserDoing
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).Find(&rows).Error; err != nil {
		return nil, err
	}
	ids := make([]uint, 0, len(rows))
	for _, row := range rows {
		ids = append(ids, row.ThingID)
	}
	return ids, nil
}

func (r *profileRepository) FindWantIDsByUserID(ctx context.Context, userID uint) ([]uint, error) {
	var rows []model.UserWant
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).Find(&rows).Error; err != nil {
		return nil, err
	}
	ids := make([]uint, 0, len(rows))
	for _, row := range rows {
		ids = append(ids, row.ThingID)
	}
	return ids, nil
}
