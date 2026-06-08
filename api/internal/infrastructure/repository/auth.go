package repository

import (
	"context"
	"errors"
	"time"

	"github.com/katedegree/spark/api/internal/domain/model"
	domainrepo "github.com/katedegree/spark/api/internal/domain/repository"
	"gorm.io/gorm"
)

type authRepository struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) domainrepo.AuthRepository {
	return &authRepository{db: db}
}

func (r *authRepository) FindUserByEmail(ctx context.Context, email string) (*domainrepo.AuthUser, error) {
	var user model.User
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &domainrepo.AuthUser{
		ID:             user.ID,
		Email:          user.Email,
		HashedPassword: user.HashedPassword,
	}, nil
}

func (r *authRepository) CreateUser(ctx context.Context, email, hashedPassword string) (*domainrepo.AuthUser, error) {
	user := model.User{
		Email:          email,
		HashedPassword: hashedPassword,
	}
	if err := r.db.WithContext(ctx).Create(&user).Error; err != nil {
		return nil, err
	}
	return &domainrepo.AuthUser{
		ID:             user.ID,
		Email:          user.Email,
		HashedPassword: user.HashedPassword,
	}, nil
}

func (r *authRepository) DeleteUserByEmail(ctx context.Context, email string) error {
	return r.db.WithContext(ctx).Where("email = ?", email).Delete(&model.User{}).Error
}

func (r *authRepository) FindOrCreateUserByEmail(ctx context.Context, email string) (*domainrepo.AuthUser, error) {
	user := model.User{Email: email, HashedPassword: ""}
	if err := r.db.WithContext(ctx).Where(model.User{Email: email}).FirstOrCreate(&user).Error; err != nil {
		return nil, err
	}
	return &domainrepo.AuthUser{
		ID:             user.ID,
		Email:          user.Email,
		HashedPassword: user.HashedPassword,
	}, nil
}

func (r *authRepository) SaveOTP(ctx context.Context, email, code string) error {
	otp := model.OTP{
		Email:     email,
		Code:      code,
		ExpiresAt: time.Now().Add(5 * time.Minute),
	}
	return r.db.WithContext(ctx).Create(&otp).Error
}

func (r *authRepository) VerifyOTP(ctx context.Context, email, code string) (bool, error) {
	var otp model.OTP
	err := r.db.WithContext(ctx).
		Where("email = ? AND code = ? AND expires_at > ?", email, code, time.Now()).
		First(&otp).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (r *authRepository) SaveRefreshToken(ctx context.Context, userID uint, token string, expiresAt time.Time) error {
	rt := model.RefreshToken{
		UserID:    userID,
		Token:     token,
		ExpiresAt: expiresAt,
	}
	return r.db.WithContext(ctx).Create(&rt).Error
}

func (r *authRepository) FindRefreshToken(ctx context.Context, token string) (*domainrepo.RefreshTokenRecord, error) {
	var rt model.RefreshToken
	err := r.db.WithContext(ctx).Where("token = ?", token).First(&rt).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &domainrepo.RefreshTokenRecord{
		UserID:    rt.UserID,
		Token:     rt.Token,
		ExpiresAt: rt.ExpiresAt,
	}, nil
}

func (r *authRepository) DeleteRefreshToken(ctx context.Context, token string) error {
	result := r.db.WithContext(ctx).Where("token = ?", token).Delete(&model.RefreshToken{})
	if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return result.Error
	}
	return nil
}
