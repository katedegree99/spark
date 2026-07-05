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
