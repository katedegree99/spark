package router

import (
	"github.com/katedegree/spark/api/internal/adapter/handler"
	"github.com/katedegree/spark/api/pkg/generated"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func NewRouter(authHandler *handler.AuthHandler) *echo.Echo {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	strict := generated.NewStrictHandler(authHandler, nil)
	generated.RegisterHandlers(e, strict)

	return e
}
