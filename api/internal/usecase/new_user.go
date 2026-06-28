package usecase

import (
	"context"

	"github.com/katedegree99/spark/api/internal/domain/repository"
)

const newUserLimit = 10

type NewUserResult struct {
	UserID  uint
	Name    string
	IconURL *string
}

type NewUserUsecase interface {
	ListNew(ctx context.Context, userID uint) ([]*NewUserResult, error)
}

type newUserUsecase struct {
	repo repository.NewUserRepository
}

func NewNewUserUsecase(repo repository.NewUserRepository) NewUserUsecase {
	return &newUserUsecase{repo: repo}
}

func (u *newUserUsecase) ListNew(ctx context.Context, userID uint) ([]*NewUserResult, error) {
	records, err := u.repo.FindNew(ctx, userID, newUserLimit)
	if err != nil {
		return nil, err
	}

	results := make([]*NewUserResult, 0, len(records))
	for _, r := range records {
		results = append(results, &NewUserResult{
			UserID:  r.UserID,
			Name:    r.Name,
			IconURL: r.IconURL,
		})
	}
	return results, nil
}
