package handler

import (
	"context"
	"errors"

	"github.com/katedegree/spark/api/internal/usecase"
	"github.com/katedegree/spark/api/pkg/generated"
)

type ThingHandler struct {
	thingUsecase usecase.ThingUsecase
}

func NewThingHandler(thingUsecase usecase.ThingUsecase) *ThingHandler {
	return &ThingHandler{thingUsecase: thingUsecase}
}

func (h *ThingHandler) ListThings(ctx context.Context, req generated.ListThingsRequestObject) (generated.ListThingsResponseObject, error) {
	q := ""
	if req.Params.Q != nil {
		q = *req.Params.Q
	}

	records, err := h.thingUsecase.Search(ctx, q)
	if err != nil {
		return nil, err
	}

	things := make([]generated.ThingResponse, 0, len(records))
	for _, r := range records {
		id := int(r.ID)
		aliases := r.Aliases
		createdAt := r.CreatedAt
		things = append(things, generated.ThingResponse{
			Id:        &id,
			Name:      &r.Name,
			Aliases:   &aliases,
			CreatedAt: &createdAt,
		})
	}

	return generated.ListThings200JSONResponse(generated.ThingsResponse{
		Things: &things,
	}), nil
}

func (h *ThingHandler) CreateThing(ctx context.Context, req generated.CreateThingRequestObject) (generated.CreateThingResponseObject, error) {
	record, err := h.thingUsecase.Create(ctx, req.Body.Name)
	if err != nil {
		if errors.Is(err, usecase.ErrThingAlreadyExists) {
			return generated.CreateThing409JSONResponse{
				Code:    "THING_ALREADY_EXISTS",
				Message: "thing with this name already exists",
			}, nil
		}
		return nil, err
	}

	id := int(record.ID)
	aliases := record.Aliases
	createdAt := record.CreatedAt
	return generated.CreateThing201JSONResponse(generated.ThingResponse{
		Id:        &id,
		Name:      &record.Name,
		Aliases:   &aliases,
		CreatedAt: &createdAt,
	}), nil
}
