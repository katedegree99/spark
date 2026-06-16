package middleware

import (
	"context"
	"errors"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type contextKey string

const userIDKey contextKey = "userID"

func JWTAuth() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			header := c.Request().Header.Get("Authorization")
			if !strings.HasPrefix(header, "Bearer ") {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"code":    "UNAUTHORIZED",
					"message": "missing or invalid authorization header",
				})
			}
			raw := strings.TrimPrefix(header, "Bearer ")

			token, err := jwt.Parse(raw, func(t *jwt.Token) (any, error) {
				if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, errors.New("unexpected signing method")
				}
				return []byte(os.Getenv("JWT_SECRET")), nil
			})
			if err != nil || !token.Valid {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"code":    "UNAUTHORIZED",
					"message": "invalid or expired token",
				})
			}

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"code":    "UNAUTHORIZED",
					"message": "invalid token claims",
				})
			}

			sub, ok := claims["sub"].(float64)
			if !ok {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"code":    "UNAUTHORIZED",
					"message": "invalid subject claim",
				})
			}

			userID := uint(sub)
			c.Set(string(userIDKey), userID)
			// Also propagate into the standard context.Context so StrictServerInterface
			// handlers can retrieve userID via UserIDFromGoContext.
			newCtx := context.WithValue(c.Request().Context(), userIDKey, userID)
			c.SetRequest(c.Request().WithContext(newCtx))
			return next(c)
		}
	}
}

func UserIDFromContext(c echo.Context) (uint, bool) {
	v := c.Get(string(userIDKey))
	if v == nil {
		return 0, false
	}
	id, ok := v.(uint)
	return id, ok
}

// UserIDFromGoContext extracts the userID from a standard context.Context.
// This is intended for use in StrictServerInterface handlers where only
// context.Context is available (not echo.Context).
func UserIDFromGoContext(ctx context.Context) (uint, bool) {
	v := ctx.Value(userIDKey)
	if v == nil {
		return 0, false
	}
	id, ok := v.(uint)
	return id, ok
}
