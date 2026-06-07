package main

import (
	"github.com/katedegree/spark/api/internal/adapter/handler"
	"github.com/katedegree/spark/api/internal/adapter/router"
	infradb "github.com/katedegree/spark/api/internal/infrastructure/db"
	infrarepo "github.com/katedegree/spark/api/internal/infrastructure/repository"
	"github.com/katedegree/spark/api/internal/usecase"
	"go.uber.org/dig"
)

func NewContainer() *dig.Container {
	c := dig.New()

	// infrastructure
	c.Provide(infradb.NewDB)
	c.Provide(infrarepo.NewAuthRepository)

	// usecase
	c.Provide(usecase.NewAuthUsecase)

	// adapter
	c.Provide(handler.NewAuthHandler)
	c.Provide(router.NewRouter)

	return c
}
