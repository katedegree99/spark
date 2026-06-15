package handler

import (
	"context"

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
	panic("not implemented")
}

func (h *ThingHandler) CreateThing(ctx context.Context, req generated.CreateThingRequestObject) (generated.CreateThingResponseObject, error) {
	panic("not implemented")
}
