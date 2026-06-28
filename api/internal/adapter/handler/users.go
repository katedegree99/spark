package handler

import (
	"context"

	authmw "github.com/katedegree99/spark/api/internal/adapter/middleware"
	"github.com/katedegree99/spark/api/internal/usecase"
	"github.com/katedegree99/spark/api/pkg/generated"
)

type UsersHandler struct {
	pickupUsecase usecase.PickupUsecase
}

func NewUsersHandler(pickupUsecase usecase.PickupUsecase) *UsersHandler {
	return &UsersHandler{pickupUsecase: pickupUsecase}
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

		matched := make([]generated.ThingResponse, 0, len(r.MatchedThings))
		for _, t := range r.MatchedThings {
			matched = append(matched, thingRecordToResponse(t))
		}
		pu.MatchedTags = &matched

		unmatched := make([]generated.ThingResponse, 0, len(r.UnmatchedThings))
		for _, t := range r.UnmatchedThings {
			unmatched = append(unmatched, thingRecordToResponse(t))
		}
		pu.UnmatchedTags = &unmatched

		users = append(users, pu)
	}

	return generated.ListPickupUsers200JSONResponse{
		Users: &users,
	}, nil
}
