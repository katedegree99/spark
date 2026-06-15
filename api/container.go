package main

import (
	"github.com/katedegree/spark/api/internal/adapter/handler"
	"github.com/katedegree/spark/api/internal/adapter/router"
	infradb "github.com/katedegree/spark/api/internal/infrastructure/db"
	infraemail "github.com/katedegree/spark/api/internal/infrastructure/email"
	infrar2 "github.com/katedegree/spark/api/internal/infrastructure/r2"
	infrarepo "github.com/katedegree/spark/api/internal/infrastructure/repository"
	"github.com/katedegree/spark/api/internal/usecase"
	"go.uber.org/dig"
)

func NewContainer() *dig.Container {
	c := dig.New()

	// infrastructure
	c.Provide(infradb.NewDB)
	c.Provide(infrarepo.NewAuthRepository)
	c.Provide(infrarepo.NewProfileRepository)
	c.Provide(infrarepo.NewThingRepository)
	c.Provide(infrarepo.NewImageRepository)
	c.Provide(infraemail.NewResendEmailService)
	c.Provide(infrar2.NewR2Service)

	// usecase
	c.Provide(usecase.NewAuthUsecase)
	c.Provide(usecase.NewProfileUsecase)
	c.Provide(usecase.NewThingUsecase)
	c.Provide(usecase.NewImageUsecase)

	// adapter
	c.Provide(handler.NewAuthHandler)
	c.Provide(handler.NewProfileHandler)
	c.Provide(handler.NewThingHandler)
	c.Provide(handler.NewImageHandler)
	c.Provide(func(
		auth *handler.AuthHandler,
		image *handler.ImageHandler,
		profile *handler.ProfileHandler,
		thing *handler.ThingHandler,
	) *handler.Handler {
		return &handler.Handler{
			AuthHandler:    auth,
			ImageHandler:   image,
			ProfileHandler: profile,
			ThingHandler:   thing,
		}
	})
	c.Provide(router.NewRouter)

	return c
}
