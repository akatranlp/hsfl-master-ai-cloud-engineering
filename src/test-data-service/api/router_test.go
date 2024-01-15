package router

import (
	"net/http"
	"net/http/httptest"
	"testing"

	health_mocks "github.com/akatranlp/hsfl-master-ai-cloud-engineering/lib/health/_mocks"
	mocks "github.com/akatranlp/hsfl-master-ai-cloud-engineering/test-data-service/_mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestRouter(t *testing.T) {
	ctrl := gomock.NewController(t)

	controller := mocks.NewMockController(ctrl)
	healthController := health_mocks.NewMockController(ctrl)
	router := New(controller, healthController)

	t.Run("/health", func(t *testing.T) {
		t.Run("should return 404 NOT FOUND if method is not GET", func(t *testing.T) {
			tests := []string{"HEAD", "POST", "PUT", "DELETE", "CONNECT", "OPTIONS", "TRACE", "PATCH"}
			for _, test := range tests {
				// given
				w := httptest.NewRecorder()
				r := httptest.NewRequest(test, "/health", nil)

				// when
				router.ServeHTTP(w, r)

				// then
				assert.Equal(t, http.StatusNotFound, w.Code)
			}
		})

		t.Run("should call GET handler", func(t *testing.T) {
			// given
			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", "/health", nil)

			healthController.
				EXPECT().
				ProvideHealth(w, r).
				Times(1)

			// when
			router.ServeHTTP(w, r)

			// then
			assert.Equal(t, http.StatusOK, w.Code)
		})
	})

	t.Run("/api/v1/reset", func(t *testing.T) {
		t.Run("should return 404 NOT FOUND if method is not POST", func(t *testing.T) {

			tests := []string{"GET", "HEAD", "PUT", "DELETE", "CONNECT", "OPTIONS", "TRACE", "PATCH"}

			for _, test := range tests {
				// given
				w := httptest.NewRecorder()
				r := httptest.NewRequest(test, "/api/v1/reset", nil)

				// when
				router.ServeHTTP(w, r)

				// then
				assert.Equal(t, http.StatusNotFound, w.Code)
			}
		})

		t.Run("should call POST handler", func(t *testing.T) {
			// given
			w := httptest.NewRecorder()
			r := httptest.NewRequest("POST", "/api/v1/reset", nil)

			controller.
				EXPECT().
				ResetDatabase(w, r).
				Times(1)

			// when
			router.ServeHTTP(w, r)

			// then
			assert.Equal(t, http.StatusOK, w.Code)
		})
	})

}
