package health

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDefaultController(t *testing.T) {
	controller := NewDefaultController()

	t.Run("ProvideHealth", func(t *testing.T) {
		// given
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/health", nil)

		// when
		controller.ProvideHealth(w, r)

		// then
		assert.Equal(t, http.StatusOK, w.Code)
	})
}
