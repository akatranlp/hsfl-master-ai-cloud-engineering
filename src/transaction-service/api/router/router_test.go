package router

import (
	"fmt"
	auth_mocks "github.com/akatranlp/hsfl-master-ai-cloud-engineering/lib/auth-middleware/_mocks"
	libRouter "github.com/akatranlp/hsfl-master-ai-cloud-engineering/lib/router"
	mocks "github.com/akatranlp/hsfl-master-ai-cloud-engineering/transaction-service/_mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRouter(t *testing.T) {
	ctrl := gomock.NewController(t)

	transactionsController := mocks.NewMockController(ctrl)
	authController := auth_mocks.NewMockController(ctrl)
	router := New(transactionsController, authController)

	t.Run("middleware /api/v1/transactions", func(t *testing.T) {
		t.Run("GetYourTransactions should not be called", func(t *testing.T) {
			// given
			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", "/api/v1/transactions", nil)

			authController.
				EXPECT().
				AuthenticationMiddleware(w, r, gomock.Any()).Do(func(w http.ResponseWriter, r *http.Request, next libRouter.Next) {
				w.WriteHeader(http.StatusUnauthorized)
			}).
				Times(1)

			// when
			router.ServeHTTP(w, r)

			// then
			assert.Equal(t, http.StatusUnauthorized, w.Code)
		})

		t.Run("GetYourTransactions should be called", func(t *testing.T) {
			// given
			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", "/api/v1/transactions", nil)

			authController.
				EXPECT().
				AuthenticationMiddleware(w, r, gomock.Any()).Do(func(w http.ResponseWriter, r *http.Request, next libRouter.Next) {
				next(r)
			}).
				Times(1)

			transactionsController.
				EXPECT().
				GetYourTransactions(w, r).
				Times(1)

			// when
			router.ServeHTTP(w, r)

			// then
			assert.Equal(t, http.StatusOK, w.Code)
		})
	})

	t.Run("/api/v1/transactions", func(t *testing.T) {
		t.Run("should return 404 NOT FOUND if method is not GET or POST", func(t *testing.T) {
			tests := []string{"DELETE", "PUT", "HEAD", "CONNECT", "OPTIONS", "TRACE", "PATCH"}
			for _, test := range tests {
				fmt.Println(test)
				// given
				w := httptest.NewRecorder()
				r := httptest.NewRequest(test, "/api/v1/books", nil)

				// when
				router.ServeHTTP(w, r)

				// then
				assert.Equal(t, http.StatusNotFound, w.Code)
			}
		})
		t.Run("should call GET handler", func(t *testing.T) {
			// given
			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", "/api/v1/transactions", nil)

			authController.
				EXPECT().
				AuthenticationMiddleware(w, r, gomock.Any()).Do(func(w http.ResponseWriter, r *http.Request, next libRouter.Next) {
				next(r)
			}).
				Times(1)

			transactionsController.
				EXPECT().
				GetYourTransactions(w, r).
				Times(1)

			// when
			router.ServeHTTP(w, r)

			// then
			assert.Equal(t, http.StatusOK, w.Code)
		})

		t.Run("should call POST handler", func(t *testing.T) {
			// given
			w := httptest.NewRecorder()
			r := httptest.NewRequest("POST", "/api/v1/transactions", nil)

			authController.
				EXPECT().
				AuthenticationMiddleware(w, r, gomock.Any()).Do(func(w http.ResponseWriter, r *http.Request, next libRouter.Next) {
				next(r)
			}).
				Times(1)

			transactionsController.
				EXPECT().
				CreateTransaction(w, r).
				Times(1)

			// when
			router.ServeHTTP(w, r)

			// then
			assert.Equal(t, http.StatusOK, w.Code)
		})
	})
}
