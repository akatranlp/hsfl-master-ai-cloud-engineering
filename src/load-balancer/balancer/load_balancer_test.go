package balancer

import (
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"net/url"
	"testing"

	mocks "github.com/akatranlp/hsfl-master-ai-cloud-engineering/load-balancer/_mocks"
	"github.com/akatranlp/hsfl-master-ai-cloud-engineering/load-balancer/balancer/target"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

type HandlerTest struct {
	statusCode int
	called     bool
}

func (h *HandlerTest) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(h.statusCode)
	h.called = true
}

func TestLoadBalancerStrategy(t *testing.T) {
	ctrl := gomock.NewController(t)
	strategyMock := mocks.NewMockStrategy(ctrl)

	serverCalled := false
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		serverCalled = true
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	url1, _ := url.Parse(server.URL)

	target1 := target.NewTarget(url1, httputil.NewSingleHostReverseProxy(url1))

	client := server.Client()
	lb := NewLoadBalancer([]*target.Target{target1}, 0, client, strategyMock)

	t.Run("Should call the server function", func(t *testing.T) {
		// given
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/api/v1/books", nil)

		handler := &HandlerTest{statusCode: http.StatusOK}
		target1 := target.NewTarget(url1, handler)

		strategyMock.
			EXPECT().
			NextTarget(r).
			Return(target1)

		// when
		lb.ServeHTTP(w, r)

		// then
		assert.Equal(t, http.StatusOK, w.Code)
		assert.True(t, handler.called)
	})

	t.Run("Test Health Check", func(t *testing.T) {
		// given
		serverCalled = false

		target1 := target.NewTarget(url1, nil)
		target1.Health = 1

		targets := []*target.Target{target1}
		lb.targets = targets

		strategyMock.
			EXPECT().
			SetTargets(targets)

		// when
		lb.HealthCheck()

		// then
		assert.True(t, serverCalled)
		assert.Equal(t, 0, target1.Health)
	})

	t.Run("Test Health Check", func(t *testing.T) {
		// given
		serverCalled = false

		url2, _ := url.Parse("http://localhost:6969")
		target1 := target.NewTarget(url2, nil)

		targets := []*target.Target{target1}
		lb.targets = targets

		strategyMock.
			EXPECT().
			SetTargets([]*target.Target{})

		// when
		lb.HealthCheck()

		// then
		assert.False(t, serverCalled)
		assert.Equal(t, 1, target1.Health)
	})
}
