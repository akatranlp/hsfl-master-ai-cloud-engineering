package httpproxy

import (
	"bytes"
	"errors"
	mocks "github.com/akatranlp/hsfl-master-ai-cloud-engineering/lib/client/_mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHTTPProxy(t *testing.T) {
	ctrl := gomock.NewController(t)
	client := mocks.NewMockClient(ctrl)

	t.Run("throw an error if no hosts are provided", func(t *testing.T) {
		// given
		proxy := NewHTTPProxy(client)

		// when
		err := AddToProxy(proxy, "*", "/the/route", []string{})

		// then
		assert.NotNil(t, err)
		assert.Equal(t, "there was no host provided", err.Error())
	})

	t.Run("throw an error if hosts is not a valid url", func(t *testing.T) {
		// given
		proxy := NewHTTPProxy(client)

		// when
		err := AddToProxy(proxy, "*", "/the/route", []string{"\n"})

		// then
		assert.NotNil(t, err)
		assert.Equal(t, "invalid origin server URL", err.Error())
	})

	t.Run("should return 404 NOT FOUND if host is unknown", func(t *testing.T) {
		// given
		proxy := NewHTTPProxy(client)
		AddToProxy(proxy, "fake.host.com", "/the/route", []string{"http://new-host:3000"})

		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/the/route", nil)

		// when
		proxy.ServeHTTP(w, r)

		// then
		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("should return 404 NOT FOUND if path is unknown", func(t *testing.T) {
		// given
		proxy := NewHTTPProxy(client)
		AddToProxy(proxy, "example.com", "/the/route", []string{"http://new-host:3000"})

		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/false/route", nil)

		// when
		proxy.ServeHTTP(w, r)

		// then
		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("should call the client because it matched the path and error", func(t *testing.T) {
		// given
		proxy := NewHTTPProxy(client)
		AddToProxy(proxy, "(.*)", "/the/route", []string{"http://new-host:3000"})

		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/the/route", nil)

		// when
		client.EXPECT().Do(r).Return(nil, errors.New("got an Error"))
		proxy.ServeHTTP(w, r)

		// then
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("should call the client because it matched the path", func(t *testing.T) {
		// given
		proxy := NewHTTPProxy(client)
		AddToProxy(proxy, "*", "/the/route", []string{"http://new-host:3000"})

		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/the/route", nil)

		response := &http.Response{
			Status:        "200 OK",
			StatusCode:    http.StatusOK,
			Header:        http.Header{},
			Body:          http.NoBody,
			ContentLength: 0,
		}

		// when
		client.EXPECT().Do(r).Return(response, nil)
		proxy.ServeHTTP(w, r)

		// then
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "example.com", r.Host)
		assert.Equal(t, "http", r.URL.Scheme)
		assert.Equal(t, strings.Split(r.RemoteAddr, ":")[0], r.Header.Get("X-Forwarded-For"))
		assert.Equal(t, "example.com", r.Header.Get("X-Forwarded-Host"))
	})

	t.Run("should copy all header values from the server", func(t *testing.T) {
		// given
		proxy := NewHTTPProxy(client)
		AddToProxy(proxy, "*", "/the/route", []string{"http://new-host:3000"})

		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/the/route", nil)

		headers := http.Header{}
		headers.Set("Set-Cookie", "refresh-token=here_should_be_a_token; Domain=example.com; Secure; HttpOnly; expires=Thu, 18 Dec 2023 20:00:00 UTC; path=/refresh-token")
		headers.Set("Custom-Header", "custom-value")

		response := &http.Response{
			Status:        "200 OK",
			StatusCode:    http.StatusOK,
			Header:        headers,
			Body:          http.NoBody,
			ContentLength: 0,
		}

		// when
		client.EXPECT().Do(r).Return(response, nil)
		proxy.ServeHTTP(w, r)

		// then
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, headers.Get("Set-Cookie"), w.Header().Get("Set-Cookie"))
		assert.Equal(t, headers.Get("Custom-Header"), w.Header().Get("Custom-Header"))
	})

	t.Run("should copy the whole Body from the server", func(t *testing.T) {
		// given
		proxy := NewHTTPProxy(client)
		AddToProxy(proxy, "*", "/the/route", []string{"http://new-host:3000"})

		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/the/route", nil)

		requestBodyContent := []byte("Hello World")

		response := &http.Response{
			Status:        "200 OK",
			StatusCode:    http.StatusOK,
			Header:        http.Header{},
			Body:          io.NopCloser(bytes.NewBuffer(requestBodyContent)),
			ContentLength: int64(len(requestBodyContent)),
		}

		// when
		client.EXPECT().Do(r).Return(response, nil)
		proxy.ServeHTTP(w, r)

		responseBodyContent, _ := io.ReadAll(w.Result().Body)

		// then
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, requestBodyContent, responseBodyContent)
	})

	t.Run("should call both hosts after another", func(t *testing.T) {
		// given
		proxy := NewHTTPProxy(client)
		AddToProxy(proxy, "*", "/the/route", []string{"http://new-host:3000", "http://second-host:8000"})

		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/the/route", nil)

		response := &http.Response{
			Status:        "200 OK",
			StatusCode:    http.StatusOK,
			Header:        http.Header{},
			Body:          http.NoBody,
			ContentLength: 0,
		}

		// when
		client.EXPECT().Do(r).Return(response, nil).Times(2)
		proxy.ServeHTTP(w, r)

		// then
		assert.Equal(t, http.StatusOK, w.Code)

		// given a second time
		w = httptest.NewRecorder()
		r.Host = "example.com"
		r.Header.Del("X-Forwarded-For")
		r.Header.Del("X-Forwarded-Host")

		// when a second time
		proxy.ServeHTTP(w, r)

		// then
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("should use the whole path from the new host and the normal path", func(t *testing.T) {
		// given
		proxy := NewHTTPProxy(client)
		AddToProxy(proxy, "*", "/the/route", []string{"http://new-host:3000/append"})

		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/the/route", nil)

		response := &http.Response{
			Status:        "200 OK",
			StatusCode:    http.StatusOK,
			Header:        http.Header{},
			Body:          http.NoBody,
			ContentLength: 0,
		}

		// when
		client.EXPECT().Do(r).Return(response, nil)
		proxy.ServeHTTP(w, r)

		// then
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "/append/the/route", r.URL.Path)
	})
}
