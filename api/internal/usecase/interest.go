package usecase

import (
	"context"
	"errors"

	"github.com/katedegree99/spark/api/internal/domain/repository"
)

var ErrAlreadyInterested = errors.New("already interested")

type SendInterestResult struct {
	Matched bool
	RoomID  *uint
}

type InterestUsecase interface {
	SendInterest(ctx context.Context, fromUserID, toUserID uint) (*SendInterestResult, error)
}

type interestUsecase struct {
	interestRepo repository.InterestRepository
	profileRepo  repository.ProfileRepository
	roomRepo     repository.RoomRepository
}

func NewInterestUsecase(
	interestRepo repository.InterestRepository,
	profileRepo repository.ProfileRepository,
	roomRepo repository.RoomRepository,
) InterestUsecase {
	return &interestUsecase{
		interestRepo: interestRepo,
		profileRepo:  profileRepo,
		roomRepo:     roomRepo,
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
	if !matched {
		return &SendInterestResult{Matched: false}, nil
	}

	// マッチ成立：ルームが未作成なら作成する
	room, err := u.roomRepo.FindByUsers(ctx, fromUserID, toUserID)
	if err != nil {
		return nil, err
	}
	if room == nil {
		room, err = u.roomRepo.Create(ctx, fromUserID, toUserID)
		if err != nil {
			return nil, err
		}
	}
	return &SendInterestResult{Matched: true, RoomID: &room.ID}, nil
}
