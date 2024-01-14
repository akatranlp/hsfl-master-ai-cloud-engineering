package auth_middleware

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	mocks "github.com/akatranlp/hsfl-master-ai-cloud-engineering/lib/auth-middleware/_mocks"
	"github.com/stretchr/testify/assert"

	"go.uber.org/mock/gomock"
)

func TestDefaultController(t *testing.T) {
	ctrl := gomock.NewController(t)
	repository := mocks.NewMockRepository(ctrl)

	t.Run("Auth is deactivated", func(t *testing.T) {
		controller := NewDefaultController(repository, false)

		t.Run("AuthenticationMiddleware", func(t *testing.T) {
			t.Run("should return 200 if auth is not active", func(t *testing.T) {
				// given
				w := httptest.NewRecorder()
				r := httptest.NewRequest("GET", "/api/v1/books", nil)

				// when
				called := false
				controller.AuthenticationMiddleware(w, r, func(req *http.Request) {
					called = true
					r = req
				})

				// then
				assert.True(t, called)
				assert.Equal(t, uint64(1), r.Context().Value(AuthenticatedUserId))
				assert.Equal(t, http.StatusOK, w.Code)
			})
		})
	})

	t.Run("Auth is activated", func(t *testing.T) {
		controller := NewDefaultController(repository, true)

		t.Run("AuthenticationMiddleware", func(t *testing.T) {
			t.Run("should return 401 if token is not provided or invalid", func(t *testing.T) {
				// given
				w := httptest.NewRecorder()
				r := httptest.NewRequest("GET", "/api/v1/books", nil)
				r.Header.Set("Authorization", "")

				// when
				called := false
				controller.AuthenticationMiddleware(w, r, func(req *http.Request) {
					called = true
				})

				// then
				assert.False(t, called)
				assert.Equal(t, http.StatusUnauthorized, w.Code)
			})

			t.Run("should return 401 if token is invalid", func(t *testing.T) {
				// given
				w := httptest.NewRecorder()
				r := httptest.NewRequest("GET", "/api/v1/books", nil)
				r.Header.Set("Authorization", "Bearer invalid-token")

				repository.
					EXPECT().
					VerifyToken("invalid-token").
					Return(uint64(0), errors.New("invalid token"))

				// when
				called := false
				controller.AuthenticationMiddleware(w, r, func(req *http.Request) {
					called = true
				})

				// then
				assert.False(t, called)
				assert.Equal(t, http.StatusUnauthorized, w.Code)
			})

			t.Run("should return 200 if token is valid", func(t *testing.T) {
				// given
				userId := uint64(1)
				w := httptest.NewRecorder()
				r := httptest.NewRequest("GET", "/api/v1/books", nil)
				r.Header.Set("Authorization", "Bearer invalid-token")

				repository.
					EXPECT().
					VerifyToken("invalid-token").
					Return(userId, nil)

				// when
				called := false
				controller.AuthenticationMiddleware(w, r, func(req *http.Request) {
					called = true
					r = req
				})

				// then
				assert.True(t, called)
				assert.Equal(t, userId, r.Context().Value(AuthenticatedUserId))
				assert.Equal(t, http.StatusOK, w.Code)
			})
		})
	})

}
