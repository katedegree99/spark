package handler

import (
	"context"

	authmw "github.com/katedegree99/spark/api/internal/adapter/middleware"
	"github.com/katedegree99/spark/api/internal/domain/repository"
	"github.com/katedegree99/spark/api/internal/usecase"
	"github.com/katedegree99/spark/api/pkg/generated"
)

type UsersHandler struct {
	pickupUsecase    usecase.PickupUsecase
	newUserUsecase   usecase.NewUserUsecase
	recommendUsecase usecase.RecommendUsecase
	listUsersUsecase usecase.ListUsersUsecase
}

func NewUsersHandler(
	pickupUsecase usecase.PickupUsecase,
	newUserUsecase usecase.NewUserUsecase,
	recommendUsecase usecase.RecommendUsecase,
	listUsersUsecase usecase.ListUsersUsecase,
) *UsersHandler {
	return &UsersHandler{
		pickupUsecase:    pickupUsecase,
		newUserUsecase:   newUserUsecase,
		recommendUsecase: recommendUsecase,
		listUsersUsecase: listUsersUsecase,
	}
}

func tagRecordToResponse(t *repository.ThingRecord) generated.TagResponse {
	id := int(t.ID)
	return generated.TagResponse{Id: &id, Name: &t.Name}
}

func (h *UsersHandler) ListPickupUsers(ctx context.Context, _ generated.ListPickupUsersRequestObject) (generated.ListPickupUsersResponseObject, error) {
	userID, ok := authmw.UserIDFromGoContext(ctx)
	if !ok {
		return generated.ListPickupUsers401JSONResponse{
			UnauthorizedJSONResponse: generated.UnauthorizedJSONResponse{
				Code:    "UNAUTHORIZED",
				Message: "unauthorized",
			},
		}, nil
	}

	results, err := h.pickupUsecase.ListPickup(ctx, userID)
	if err != nil {
		return nil, err
	}

	users := make([]generated.PickupUserResponse, 0, len(results))
	for _, r := range results {
		userIDInt := int(r.UserID)
		pu := generated.PickupUserResponse{
			UserId:  &userIDInt,
			Name:    &r.Name,
			Bio:     r.Bio,
			IconUrl: r.IconURL,
		}

		matched := make([]generated.TagResponse, 0, len(r.MatchedThings))
		for _, t := range r.MatchedThings {
			matched = append(matched, tagRecordToResponse(t))
		}
		pu.MatchedTags = &matched

		unmatched := make([]generated.TagResponse, 0, len(r.UnmatchedThings))
		for _, t := range r.UnmatchedThings {
			unmatched = append(unmatched, tagRecordToResponse(t))
		}
		pu.UnmatchedTags = &unmatched

		users = append(users, pu)
	}

	return generated.ListPickupUsers200JSONResponse{
		Users: &users,
	}, nil
}

func (h *UsersHandler) ListNewUsers(ctx context.Context, _ generated.ListNewUsersRequestObject) (generated.ListNewUsersResponseObject, error) {
	userID, _ := authmw.UserIDFromGoContext(ctx)

	results, err := h.newUserUsecase.ListNew(ctx, userID)
	if err != nil {
		return nil, err
	}

	users := make([]generated.NewUserResponse, 0, len(results))
	for _, r := range results {
		userIDInt := int(r.UserID)
		users = append(users, generated.NewUserResponse{
			UserId:  &userIDInt,
			Name:    &r.Name,
			IconUrl: r.IconURL,
		})
	}

	return generated.ListNewUsers200JSONResponse{
		Users: &users,
	}, nil
}

func (h *UsersHandler) ListRecommendUsers(ctx context.Context, _ generated.ListRecommendUsersRequestObject) (generated.ListRecommendUsersResponseObject, error) {
	userID, ok := authmw.UserIDFromGoContext(ctx)
	if !ok {
		return generated.ListRecommendUsers401JSONResponse{
			UnauthorizedJSONResponse: generated.UnauthorizedJSONResponse{
				Code:    "UNAUTHORIZED",
				Message: "unauthorized",
			},
		}, nil
	}

	results, err := h.recommendUsecase.ListRecommend(ctx, userID)
	if err != nil {
		return nil, err
	}

	users := make([]generated.RecommendUserResponse, 0, len(results))
	for _, r := range results {
		userIDInt := int(r.UserID)
		commonCount := r.CommonCount
		ru := generated.RecommendUserResponse{
			UserId:      &userIDInt,
			Name:        &r.Name,
			Bio:         r.Bio,
			IconUrl:     r.IconURL,
			CommonCount: &commonCount,
		}

		matched := make([]generated.TagResponse, 0, len(r.MatchedThings))
		for _, t := range r.MatchedThings {
			matched = append(matched, tagRecordToResponse(t))
		}
		ru.MatchedTags = &matched

		unmatched := make([]generated.TagResponse, 0, len(r.UnmatchedThings))
		for _, t := range r.UnmatchedThings {
			unmatched = append(unmatched, tagRecordToResponse(t))
		}
		ru.UnmatchedTags = &unmatched

		users = append(users, ru)
	}

	return generated.ListRecommendUsers200JSONResponse{
		Users: &users,
	}, nil
}

func (h *UsersHandler) ListUsers(ctx context.Context, request generated.ListUsersRequestObject) (generated.ListUsersResponseObject, error) {
	userID, ok := authmw.UserIDFromGoContext(ctx)
	if !ok {
		return generated.ListUsers401JSONResponse{
			UnauthorizedJSONResponse: generated.UnauthorizedJSONResponse{
				Code:    "UNAUTHORIZED",
				Message: "unauthorized",
			},
		}, nil
	}

	input := usecase.ListUsersInput{}
	if request.Params.Q != nil {
		input.Q = request.Params.Q
	}
	if request.Params.TagIds != nil {
		tagIDs := make([]uint, 0, len(*request.Params.TagIds))
		for _, id := range *request.Params.TagIds {
			tagIDs = append(tagIDs, uint(id))
		}
		input.TagIDs = tagIDs
	}
	if request.Params.Offset != nil {
		input.Offset = *request.Params.Offset
	}
	input.Limit = 20
	if request.Params.Limit != nil {
		input.Limit = *request.Params.Limit
	}

	results, err := h.listUsersUsecase.ListUsers(ctx, userID, input)
	if err != nil {
		return nil, err
	}

	users := make([]generated.RecommendUserResponse, 0, len(results))
	for _, r := range results {
		userIDInt := int(r.UserID)
		commonCount := r.CommonCount
		ru := generated.RecommendUserResponse{
			UserId:      &userIDInt,
			Name:        &r.Name,
			Bio:         r.Bio,
			IconUrl:     r.IconURL,
			CommonCount: &commonCount,
		}

		matched := make([]generated.TagResponse, 0, len(r.MatchedThings))
		for _, t := range r.MatchedThings {
			matched = append(matched, tagRecordToResponse(t))
		}
		ru.MatchedTags = &matched

		unmatched := make([]generated.TagResponse, 0, len(r.UnmatchedThings))
		for _, t := range r.UnmatchedThings {
			unmatched = append(unmatched, tagRecordToResponse(t))
		}
		ru.UnmatchedTags = &unmatched

		users = append(users, ru)
	}

	return generated.ListUsers200JSONResponse{
		Users: &users,
	}, nil
}
