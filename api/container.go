package main

import (
	"github.com/katedegree/spark/api/internal/adapter/handler"
	"github.com/katedegree/spark/api/internal/adapter/router"
	infradb "github.com/katedegree/spark/api/internal/infrastructure/db"
	infraemail "github.com/katedegree/spark/api/internal/infrastructure/email"
	infrallm "github.com/katedegree/spark/api/internal/infrastructure/llm"
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
	c.Provide(infrarepo.NewPickupRepository)
	c.Provide(infraemail.NewResendEmailService)
	c.Provide(infrar2.NewR2Service)
	c.Provide(infrallm.NewClaudeAliasService)

	// usecase
	c.Provide(usecase.NewAuthUsecase)
	c.Provide(usecase.NewProfileUsecase)
	c.Provide(usecase.NewThingUsecase)
	c.Provide(usecase.NewImageUsecase)
	c.Provide(usecase.NewPickupUsecase)

	// adapter
	c.Provide(handler.NewAuthHandler)
	c.Provide(handler.NewProfileHandler)
	c.Provide(handler.NewThingHandler)
	c.Provide(handler.NewImageHandler)
	c.Provide(handler.NewUsersHandler)
	c.Provide(func(
		auth *handler.AuthHandler,
		image *handler.ImageHandler,
		profile *handler.ProfileHandler,
		thing *handler.ThingHandler,
		users *handler.UsersHandler,
	) *handler.Handler {
		return &handler.Handler{
			AuthHandler:    auth,
			ImageHandler:   image,
			ProfileHandler: profile,
			ThingHandler:   thing,
			UsersHandler:   users,
		}
	})
	c.Provide(router.NewRouter)

	return c
}
