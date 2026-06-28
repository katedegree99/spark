package repository

import (
	"context"

	domainrepo "github.com/katedegree/spark/api/internal/domain/repository"
	"gorm.io/gorm"
)

type newUserRepository struct {
	db *gorm.DB
}

func NewNewUserRepository(db *gorm.DB) domainrepo.NewUserRepository {
	return &newUserRepository{db: db}
}

func (r *newUserRepository) FindNew(ctx context.Context, excludeUserID uint, limit int) ([]*domainrepo.NewUserRecord, error) {
	type row struct {
		UserID  uint
		Name    string
		IconURL *string
	}

	var rows []row
	err := r.db.WithContext(ctx).
		Table("profiles p").
		Select("p.user_id, p.name, i.url as icon_url").
		Joins("LEFT JOIN images i ON p.icon_image_id = i.id").
		Where("p.user_id != ?", excludeUserID).
		Order("p.created_at DESC").
		Limit(limit).
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}

	records := make([]*domainrepo.NewUserRecord, 0, len(rows))
	for _, r := range rows {
		records = append(records, &domainrepo.NewUserRecord{
			UserID:  r.UserID,
			Name:    r.Name,
			IconURL: r.IconURL,
		})
	}
	return records, nil
}
