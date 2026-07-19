package handler_test

import (
	"bytes"
	"context"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/katedegree99/spark/api/internal/adapter/handler"
	"github.com/katedegree99/spark/api/internal/adapter/router"
	"github.com/katedegree99/spark/api/internal/domain/repository"
	"github.com/katedegree99/spark/api/internal/usecase"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockImageUsecase struct {
	fn func(ctx context.Context, directory, filename string, r io.Reader) (*repository.ImageRecord, error)
}

func (m *mockImageUsecase) Upload(ctx context.Context, directory, filename string, r io.Reader) (*repository.ImageRecord, error) {
	if m.fn != nil {
		return m.fn(ctx, directory, filename, r)
	}
	return &repository.ImageRecord{ID: 1, Directory: directory, URL: "https://example.com/icon.jpg"}, nil
}

func newTestEchoImage(uc usecase.ImageUsecase) *echo.Echo {
	h := &handler.Handler{
		ImageHandler: handler.NewImageHandler(uc),
	}
	return router.NewRouter(h, handler.NewWsHub())
}

func buildMultipart(t *testing.T, directory, filename, content string) (*bytes.Buffer, string) {
	t.Helper()
	body := &bytes.Buffer{}
	w := multipart.NewWriter(body)

	if directory != "" {
		fw, err := w.CreateFormField("directory")
		require.NoError(t, err)
		_, _ = fw.Write([]byte(directory))
	}

	if filename != "" {
		fw, err := w.CreateFormFile("file", filename)
		require.NoError(t, err)
		_, _ = fw.Write([]byte(content))
	}

	w.Close()
	return body, w.FormDataContentType()
}

func doMultipart(e *echo.Echo, body *bytes.Buffer, contentType, auth string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(http.MethodPost, "/images", body)
	req.Header.Set("Content-Type", contentType)
	if auth != "" {
		req.Header.Set("Authorization", auth)
	}
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	return rec
}

// --- POST /images ---

func TestUploadImage_Success(t *testing.T) {
	uc := &mockImageUsecase{
		fn: func(_ context.Context, directory, filename string, _ io.Reader) (*repository.ImageRecord, error) {
			return &repository.ImageRecord{ID: 1, Directory: directory, URL: "https://cdn.example.com/" + filename}, nil
		},
	}
	body, ct := buildMultipart(t, "icon", "avatar.jpg", "fake-image-bytes")
	rec := doMultipart(newTestEchoImage(uc), body, ct, testBearerToken(1))
	assert.Equal(t, http.StatusCreated, rec.Code)
	respBody := rec.Body.String()
	require.Contains(t, respBody, "url")
	require.Contains(t, respBody, "directory")
}

func TestUploadImage_Unauthorized(t *testing.T) {
	uc := &mockImageUsecase{}
	body, ct := buildMultipart(t, "icon", "avatar.jpg", "fake")
	rec := doMultipart(newTestEchoImage(uc), body, ct, "")
	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestUploadImage_InvalidDirectory(t *testing.T) {
	uc := &mockImageUsecase{}
	body, ct := buildMultipart(t, "invalid_dir", "avatar.jpg", "fake")
	rec := doMultipart(newTestEchoImage(uc), body, ct, testBearerToken(1))
	assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
	assert.Contains(t, rec.Body.String(), "VALIDATION_ERROR")
}

func TestUploadImage_MissingDirectory(t *testing.T) {
	uc := &mockImageUsecase{}
	body, ct := buildMultipart(t, "", "avatar.jpg", "fake")
	rec := doMultipart(newTestEchoImage(uc), body, ct, testBearerToken(1))
	assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
}

func TestUploadImage_MissingFile(t *testing.T) {
	uc := &mockImageUsecase{}
	// directory のみで file なし
	body := &bytes.Buffer{}
	w := multipart.NewWriter(body)
	fw, err := w.CreateFormField("directory")
	require.NoError(t, err)
	_, _ = fw.Write([]byte("icon"))
	w.Close()

	rec := doMultipart(newTestEchoImage(uc), body, w.FormDataContentType(), testBearerToken(1))
	assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
}

func TestUploadImage_InternalError(t *testing.T) {
	uc := &mockImageUsecase{
		fn: func(_ context.Context, _, _ string, _ io.Reader) (*repository.ImageRecord, error) {
			return nil, errors.New("r2 error")
		},
	}
	body, ct := buildMultipart(t, "icon", "test.jpg", "fake")
	rec := doMultipart(newTestEchoImage(uc), body, ct, testBearerToken(1))
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}
