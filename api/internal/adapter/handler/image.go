package handler

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/google/uuid"
	"github.com/katedegree99/spark/api/internal/usecase"
	"github.com/katedegree99/spark/api/pkg/generated"
)

type ImageHandler struct {
	imageUsecase usecase.ImageUsecase
}

func NewImageHandler(imageUsecase usecase.ImageUsecase) *ImageHandler {
	return &ImageHandler{imageUsecase: imageUsecase}
}

func (h *ImageHandler) UploadImage(ctx context.Context, req generated.UploadImageRequestObject) (generated.UploadImageResponseObject, error) {
	var directory string
	var filename string
	var fileReader io.Reader

	mr := req.Body
	for {
		part, err := mr.NextPart()
		if err == io.EOF {
			break
		}
		if err != nil {
			return generated.UploadImage422JSONResponse{
				ValidationErrorJSONResponse: generated.ValidationErrorJSONResponse{
					Message: "failed to read multipart: " + err.Error(),
					Code:    "VALIDATION_ERROR",
				},
			}, nil
		}

		switch part.FormName() {
		case "directory":
			data, err := io.ReadAll(part)
			if err != nil {
				return generated.UploadImage422JSONResponse{
					ValidationErrorJSONResponse: generated.ValidationErrorJSONResponse{
						Message: "failed to read directory field",
						Code:    "VALIDATION_ERROR",
					},
				}, nil
			}
			dir := generated.UploadImageMultipartBodyDirectory(data)
			if !dir.Valid() {
				return generated.UploadImage422JSONResponse{
					ValidationErrorJSONResponse: generated.ValidationErrorJSONResponse{
						Message: "directory must be one of: icon",
						Code:    "VALIDATION_ERROR",
					},
				}, nil
			}
			directory = string(dir)
		case "file":
			originalFilename := part.FileName()
			if originalFilename == "" {
				originalFilename = "upload"
			}
			filename = uuid.New().String() + "_" + originalFilename
			fileReader = part
		}
	}

	if directory == "" || fileReader == nil {
		return generated.UploadImage422JSONResponse{
			ValidationErrorJSONResponse: generated.ValidationErrorJSONResponse{
				Message: "directory and file are required",
				Code:    "VALIDATION_ERROR",
			},
		}, nil
	}

	record, err := h.imageUsecase.Upload(ctx, directory, filename, fileReader)
	if err != nil {
		return nil, fmt.Errorf("upload image: %w", err)
	}

	id := int(record.ID)
	dir := record.Directory
	url := record.URL
	now := time.Now()

	return generated.UploadImage201JSONResponse(generated.ImageResponse{
		Id:        &id,
		Directory: &dir,
		Url:       &url,
		CreatedAt: &now,
	}), nil
}
