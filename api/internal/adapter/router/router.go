package router

import (
	"github.com/katedegree99/spark/api/internal/adapter/handler"
	authmw "github.com/katedegree99/spark/api/internal/adapter/middleware"
	"github.com/katedegree99/spark/api/pkg/generated"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func NewRouter(h *handler.Handler) *echo.Echo {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	strict := generated.NewStrictHandler(h, []generated.StrictMiddlewareFunc{
		jwtMiddlewareAdapter(authmw.JWTAuth()),
	})
	generated.RegisterHandlers(e, strict)

	return e
}

// jwtMiddlewareAdapter applies Echo middleware only to operations that require bearerAuth.
func jwtMiddlewareAdapter(mw echo.MiddlewareFunc) generated.StrictMiddlewareFunc {
	protected := map[string]bool{
		"UploadImage":      true,
		"ListThings":       true,
		"CreateThing":      true,
		"CreateMyProfile":  true,
		"GetMyProfile":     true,
		"UpdateMyProfile":  true,
		"Logout":           true,
		"ListPickupUsers":    true,
		"ListRecommendUsers": true,
	}
	return func(f generated.StrictHandlerFunc, operationID string) generated.StrictHandlerFunc {
		if !protected[operationID] {
			return f
		}
		return func(c echo.Context, req any) (any, error) {
			var nextCalled bool
			if err := mw(func(c echo.Context) error {
				nextCalled = true
				return nil
			})(c); err != nil || !nextCalled {
				return nil, err
			}
			return f(c, req)
		}
	}
}
