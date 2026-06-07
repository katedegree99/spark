package handler

import (
	"context"

	"github.com/katedegree/spark/api/internal/usecase"
	"github.com/katedegree/spark/api/pkg/generated"
)

type AuthHandler struct {
	authUsecase usecase.AuthUsecase
}

func NewAuthHandler(authUsecase usecase.AuthUsecase) *AuthHandler {
	return &AuthHandler{authUsecase: authUsecase}
}

func (h *AuthHandler) RegisterWithEmail(_ context.Context, _ generated.RegisterWithEmailRequestObject) (generated.RegisterWithEmailResponseObject, error) {
	panic("not implemented")
}

func (h *AuthHandler) LoginWithEmail(_ context.Context, _ generated.LoginWithEmailRequestObject) (generated.LoginWithEmailResponseObject, error) {
	panic("not implemented")
}

func (h *AuthHandler) VerifyOtp(_ context.Context, _ generated.VerifyOtpRequestObject) (generated.VerifyOtpResponseObject, error) {
	panic("not implemented")
}

func (h *AuthHandler) LoginWithGoogle(_ context.Context, _ generated.LoginWithGoogleRequestObject) (generated.LoginWithGoogleResponseObject, error) {
	panic("not implemented")
}

func (h *AuthHandler) RefreshToken(_ context.Context, _ generated.RefreshTokenRequestObject) (generated.RefreshTokenResponseObject, error) {
	panic("not implemented")
}

func (h *AuthHandler) Logout(_ context.Context, _ generated.LogoutRequestObject) (generated.LogoutResponseObject, error) {
	panic("not implemented")
}

var _ generated.StrictServerInterface = (*AuthHandler)(nil)
