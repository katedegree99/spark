package handler_test

import (
	"net/http"
	"testing"

	"github.com/katedegree99/spark/api/internal/adapter/handler"
	"github.com/katedegree99/spark/api/internal/adapter/router"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func newTestEchoHealth() *echo.Echo {
	h := &handler.Handler{
		HealthHandler: handler.NewHealthHandler(),
	}
	return router.NewRouter(h, handler.NewWsHub())
}

func TestHealthCheck_AlwaysReturns200(t *testing.T) {
	rec := do(newTestEchoHealth(), http.MethodGet, "/health", "")
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), "ok")
}

func TestHealthCheck_NoAuthRequired(t *testing.T) {
	rec := do(newTestEchoHealth(), http.MethodGet, "/health", "")
	assert.Equal(t, http.StatusOK, rec.Code)
}
