package handler

import (
	"context"
	"errors"

	"github.com/katedegree/spark/api/internal/usecase"
	"github.com/katedegree/spark/api/pkg/generated"
)

type AuthHandler struct {
	authUsecase usecase.AuthUsecase
}

func NewAuthHandler(authUsecase usecase.AuthUsecase) *AuthHandler {
	return &AuthHandler{authUsecase: authUsecase}
}

func (h *AuthHandler) RegisterWithEmail(ctx context.Context, req generated.RegisterWithEmailRequestObject) (generated.RegisterWithEmailResponseObject, error) {
	email := string(req.Body.Email)
	password := req.Body.Password

	if err := h.authUsecase.Register(ctx, email, password); err != nil {
		if errors.Is(err, usecase.ErrEmailAlreadyExists) {
			return generated.RegisterWithEmail409JSONResponse{
				Code:    "EMAIL_ALREADY_EXISTS",
				Message: "email already registered",
			}, nil
		}
		return nil, err
	}

	expiresIn := 300
	message := "OTP sent to your email"
	return generated.RegisterWithEmail200JSONResponse{
		ExpiresIn: &expiresIn,
		Message:   &message,
	}, nil
}

func (h *AuthHandler) LoginWithEmail(ctx context.Context, req generated.LoginWithEmailRequestObject) (generated.LoginWithEmailResponseObject, error) {
	email := string(req.Body.Email)
	password := req.Body.Password

	if err := h.authUsecase.Login(ctx, email, password); err != nil {
		if errors.Is(err, usecase.ErrInvalidCredentials) {
			return generated.LoginWithEmail401JSONResponse{
				Code:    "INVALID_CREDENTIALS",
				Message: "invalid email or password",
			}, nil
		}
		return nil, err
	}

	expiresIn := 300
	message := "OTP sent to your email"
	return generated.LoginWithEmail200JSONResponse{
		ExpiresIn: &expiresIn,
		Message:   &message,
	}, nil
}

func (h *AuthHandler) VerifyOtp(ctx context.Context, req generated.VerifyOtpRequestObject) (generated.VerifyOtpResponseObject, error) {
	email := string(req.Body.Email)
	code := req.Body.Code

	pair, err := h.authUsecase.VerifyOTP(ctx, email, code)
	if err != nil {
		if errors.Is(err, usecase.ErrInvalidOTP) {
			return generated.VerifyOtp400JSONResponse{
				Code:    "INVALID_OTP",
				Message: err.Error(),
			}, nil
		}
		return nil, err
	}

	expiresIn := 900
	tokenType := generated.Bearer
	return generated.VerifyOtp200JSONResponse{
		AccessToken:   &pair.AccessToken,
		RefreshToken:  &pair.RefreshToken,
		ExpiresIn:     &expiresIn,
		TokenType:     &tokenType,
		ProfileExists: &pair.ProfileExists,
	}, nil
}

func (h *AuthHandler) LoginWithGoogle(ctx context.Context, req generated.LoginWithGoogleRequestObject) (generated.LoginWithGoogleResponseObject, error) {
	pair, err := h.authUsecase.LoginWithGoogle(ctx, req.Body.IdToken)
	if err != nil {
		if errors.Is(err, usecase.ErrInvalidGoogleToken) {
			return generated.LoginWithGoogle401JSONResponse{
				Code:    "INVALID_GOOGLE_TOKEN",
				Message: "invalid google id token",
			}, nil
		}
		return nil, err
	}

	expiresIn := 900
	tokenType := generated.Bearer
	return generated.LoginWithGoogle200JSONResponse{
		AccessToken:   &pair.AccessToken,
		RefreshToken:  &pair.RefreshToken,
		ExpiresIn:     &expiresIn,
		TokenType:     &tokenType,
		ProfileExists: &pair.ProfileExists,
	}, nil
}

func (h *AuthHandler) RefreshToken(ctx context.Context, req generated.RefreshTokenRequestObject) (generated.RefreshTokenResponseObject, error) {
	pair, err := h.authUsecase.RefreshToken(ctx, req.Body.RefreshToken)
	if err != nil {
		if errors.Is(err, usecase.ErrInvalidRefreshToken) {
			return generated.RefreshToken401JSONResponse{
				Code:    "INVALID_REFRESH_TOKEN",
				Message: err.Error(),
			}, nil
		}
		return nil, err
	}

	tokenType := generated.Bearer
	expiresIn := 900
	return generated.RefreshToken200JSONResponse(generated.AuthTokensResponse{
		AccessToken:  &pair.AccessToken,
		RefreshToken: &pair.RefreshToken,
		TokenType:    &tokenType,
		ExpiresIn:    &expiresIn,
	}), nil
}

func (h *AuthHandler) Logout(ctx context.Context, req generated.LogoutRequestObject) (generated.LogoutResponseObject, error) {
	if err := h.authUsecase.Logout(ctx, req.Body.RefreshToken); err != nil {
		return nil, err
	}
	return generated.Logout204Response{}, nil
}

