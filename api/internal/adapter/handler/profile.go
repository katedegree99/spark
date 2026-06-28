package handler

import (
	"context"
	"errors"

	authmw "github.com/katedegree/spark/api/internal/adapter/middleware"
	"github.com/katedegree/spark/api/internal/domain/repository"
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
	userID, ok := authmw.UserIDFromGoContext(ctx)
	if !ok {
		return generated.CreateMyProfile401JSONResponse{
			UnauthorizedJSONResponse: generated.UnauthorizedJSONResponse{
				Code:    "UNAUTHORIZED",
				Message: "unauthorized",
			},
		}, nil
	}

	input := usecase.ProfileInput{
		Name:     req.Body.Name,
		Bio:      req.Body.Bio,
		DoingIDs: toUintSlice(req.Body.DoingThingIds),
		WantIDs:  toUintSlice(req.Body.WantThingIds),
	}
	if req.Body.IconImageId != nil {
		v := uint(*req.Body.IconImageId)
		input.IconImageID = &v
	}

	result, err := h.profileUsecase.Create(ctx, userID, input)
	if err != nil {
		if errors.Is(err, usecase.ErrProfileAlreadyExists) {
			return generated.CreateMyProfile409JSONResponse{
				Code:    "PROFILE_ALREADY_EXISTS",
				Message: "profile already exists",
			}, nil
		}
		return nil, err
	}

	return generated.CreateMyProfile201JSONResponse(buildProfileResponse(result)), nil
}

func (h *ProfileHandler) GetMyProfile(ctx context.Context, req generated.GetMyProfileRequestObject) (generated.GetMyProfileResponseObject, error) {
	userID, ok := authmw.UserIDFromGoContext(ctx)
	if !ok {
		return generated.GetMyProfile401JSONResponse{
			UnauthorizedJSONResponse: generated.UnauthorizedJSONResponse{
				Code:    "UNAUTHORIZED",
				Message: "unauthorized",
			},
		}, nil
	}

	result, err := h.profileUsecase.Get(ctx, userID)
	if err != nil {
		if errors.Is(err, usecase.ErrProfileNotFound) {
			return generated.GetMyProfile404JSONResponse{
				Code:    "PROFILE_NOT_FOUND",
				Message: "profile not found",
			}, nil
		}
		return nil, err
	}

	return generated.GetMyProfile200JSONResponse(buildProfileResponse(result)), nil
}

func (h *ProfileHandler) UpdateMyProfile(ctx context.Context, req generated.UpdateMyProfileRequestObject) (generated.UpdateMyProfileResponseObject, error) {
	userID, ok := authmw.UserIDFromGoContext(ctx)
	if !ok {
		return generated.UpdateMyProfile401JSONResponse{
			UnauthorizedJSONResponse: generated.UnauthorizedJSONResponse{
				Code:    "UNAUTHORIZED",
				Message: "unauthorized",
			},
		}, nil
	}

	// Get current profile to apply partial updates.
	current, err := h.profileUsecase.Get(ctx, userID)
	if err != nil {
		if errors.Is(err, usecase.ErrProfileNotFound) {
			return generated.UpdateMyProfile404JSONResponse{
				Code:    "PROFILE_NOT_FOUND",
				Message: "profile not found",
			}, nil
		}
		return nil, err
	}

	input := usecase.ProfileInput{
		Name:        current.Profile.Name,
		Bio:         current.Profile.Bio,
		IconImageID: current.Profile.IconImageID,
		DoingIDs:    thingRecordsToIDs(current.Doings),
		WantIDs:     thingRecordsToIDs(current.Wants),
	}

	if req.Body.Name != nil {
		input.Name = *req.Body.Name
	}
	if req.Body.Bio != nil {
		input.Bio = req.Body.Bio
	}
	if req.Body.IconImageId != nil {
		v := uint(*req.Body.IconImageId)
		input.IconImageID = &v
	}
	if req.Body.DoingThingIds != nil {
		input.DoingIDs = toUintSlice(req.Body.DoingThingIds)
	}
	if req.Body.WantThingIds != nil {
		input.WantIDs = toUintSlice(req.Body.WantThingIds)
	}

	result, err := h.profileUsecase.Update(ctx, userID, input)
	if err != nil {
		if errors.Is(err, usecase.ErrProfileNotFound) {
			return generated.UpdateMyProfile404JSONResponse{
				Code:    "PROFILE_NOT_FOUND",
				Message: "profile not found",
			}, nil
		}
		return nil, err
	}

	return generated.UpdateMyProfile200JSONResponse(buildProfileResponse(result)), nil
}

// buildProfileResponse converts a ProfileResult to a generated.ProfileResponse.
func buildProfileResponse(r *usecase.ProfileResult) generated.ProfileResponse {
	userIDInt := int(r.Profile.UserID)

	resp := generated.ProfileResponse{
		UserId:    &userIDInt,
		Name:      &r.Profile.Name,
		Bio:       r.Profile.Bio,
		CreatedAt: &r.Profile.CreatedAt,
		UpdatedAt: &r.Profile.UpdatedAt,
	}

	if r.Image != nil {
		imgID := int(r.Image.ID)
		resp.IconImage = &generated.ImageResponse{
			Id:        &imgID,
			Directory: &r.Image.Directory,
			Url:       &r.Image.URL,
		}
	}

	doings := make([]generated.ThingResponse, 0, len(r.Doings))
	for _, t := range r.Doings {
		doings = append(doings, thingRecordToResponse(t))
	}
	resp.Doings = &doings

	wants := make([]generated.ThingResponse, 0, len(r.Wants))
	for _, t := range r.Wants {
		wants = append(wants, thingRecordToResponse(t))
	}
	resp.Wants = &wants

	return resp
}

func thingRecordToResponse(t *repository.ThingRecord) generated.ThingResponse {
	id := int(t.ID)
	aliases := append([]string{t.Name}, t.Aliases...)
	return generated.ThingResponse{
		Id:      &id,
		Name:    &t.Name,
		Aliases: &aliases,
	}
}

// toUintSlice converts *[]int to []uint.
func toUintSlice(ids *[]int) []uint {
	if ids == nil {
		return []uint{}
	}
	result := make([]uint, 0, len(*ids))
	for _, id := range *ids {
		result = append(result, uint(id))
	}
	return result
}

// thingRecordsToIDs extracts IDs from a ThingRecord slice.
func thingRecordsToIDs(things []*repository.ThingRecord) []uint {
	ids := make([]uint, 0, len(things))
	for _, t := range things {
		ids = append(ids, t.ID)
	}
	return ids
}
