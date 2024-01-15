package controller

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	mocks "github.com/akatranlp/hsfl-master-ai-cloud-engineering/test-data-service/_mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestDefaultController(t *testing.T) {
	ctrl := gomock.NewController(t)
	repository := mocks.NewMockRepository(ctrl)

	controller := NewDefaultController(repository)

	t.Run("ResetDatabase", func(t *testing.T) {
		t.Run("should return 500 INTERNAL SERVER ERROR if reset database fails", func(t *testing.T) {
			// given
			w := httptest.NewRecorder()
			r := httptest.NewRequest("POST", "/api/v1/reset", nil)

			repository.
				EXPECT().
				ResetDatabase().
				Return(errors.New("error")).
				Times(1)

			// when
			controller.ResetDatabase(w, r)

			// then
			assert.Equal(t, http.StatusInternalServerError, w.Code)
		})

		t.Run("should call repository.ResetDatabase", func(t *testing.T) {
			// given
			w := httptest.NewRecorder()
			r := httptest.NewRequest("POST", "/api/v1/reset", nil)

			repository.
				EXPECT().
				ResetDatabase().
				Return(nil).
				Times(1)

			// when
			controller.ResetDatabase(w, r)

			// then
			assert.Equal(t, http.StatusOK, w.Code)
		})
	})
}
