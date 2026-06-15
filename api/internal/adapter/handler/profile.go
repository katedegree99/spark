package handler

import (
	"context"

	"github.com/katedegree/spark/api/internal/usecase"
	"github.com/katedegree/spark/api/pkg/generated"
)

type ProfileHandler struct {
	profileUsecase usecase.ProfileUsecase
}

func NewProfileHandler(profileUsecase usecase.ProfileUsecase) *ProfileHandler {
	return &ProfileHandler{profileUsecase: profileUsecase}
}

func (h *ProfileHandler) CreateMyProfile(ctx context.Context, req generated.CreateMyProfileRequestObject) (generated.CreateMyProfileResponseObject, error) {
	panic("not implemented")
}

func (h *ProfileHandler) GetMyProfile(ctx context.Context, req generated.GetMyProfileRequestObject) (generated.GetMyProfileResponseObject, error) {
	panic("not implemented")
}

func (h *ProfileHandler) UpdateMyProfile(ctx context.Context, req generated.UpdateMyProfileRequestObject) (generated.UpdateMyProfileResponseObject, error) {
	panic("not implemented")
}
