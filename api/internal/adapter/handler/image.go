package handler

import (
	"context"

	"github.com/katedegree/spark/api/internal/usecase"
	"github.com/katedegree/spark/api/pkg/generated"
)

type ImageHandler struct {
	imageUsecase usecase.ImageUsecase
}

func NewImageHandler(imageUsecase usecase.ImageUsecase) *ImageHandler {
	return &ImageHandler{imageUsecase: imageUsecase}
}

func (h *ImageHandler) UploadImage(ctx context.Context, req generated.UploadImageRequestObject) (generated.UploadImageResponseObject, error) {
	panic("not implemented")
}
