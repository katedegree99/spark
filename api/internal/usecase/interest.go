package usecase

import (
	"context"
	"errors"

	"github.com/katedegree99/spark/api/internal/domain/repository"
)

var ErrAlreadyInterested = errors.New("already interested")

type SendInterestResult struct {
	Matched bool
}

type InterestUsecase interface {
	SendInterest(ctx context.Context, fromUserID, toUserID uint) (*SendInterestResult, error)
}

type interestUsecase struct {
	interestRepo repository.InterestRepository
	profileRepo  repository.ProfileRepository
}

func NewInterestUsecase(
	interestRepo repository.InterestRepository,
	profileRepo repository.ProfileRepository,
) InterestUsecase {
	return &interestUsecase{
		interestRepo: interestRepo,
		profileRepo:  profileRepo,
	}
}

func (u *interestUsecase) SendInterest(ctx context.Context, fromUserID, toUserID uint) (*SendInterestResult, error) {
	// 対象ユーザーの存在確認
	exists, err := u.profileRepo.ExistsByUserID(ctx, toUserID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, ErrUserNotFound
	}

	// 重複チェック
	already, err := u.interestRepo.Exists(ctx, fromUserID, toUserID)
	if err != nil {
		return nil, err
	}
	if already {
		return nil, ErrAlreadyInterested
	}

	if err := u.interestRepo.Create(ctx, fromUserID, toUserID); err != nil {
		return nil, err
	}

	// 相手がすでに気になるを送っているか（マッチ判定）
	matched, err := u.interestRepo.Exists(ctx, toUserID, fromUserID)
	if err != nil {
		return nil, err
	}

	return &SendInterestResult{Matched: matched}, nil
}
