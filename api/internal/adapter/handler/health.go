package handler

import (
	"context"

	"github.com/katedegree99/spark/api/pkg/generated"
)

type HealthHandler struct{}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

func (h *HealthHandler) HealthCheck(_ context.Context, _ generated.HealthCheckRequestObject) (generated.HealthCheckResponseObject, error) {
	return generated.HealthCheck200JSONResponse{Status: "ok"}, nil
}
