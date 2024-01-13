package auth_middleware

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"net/url"
	"testing"

	mocks "github.com/akatranlp/hsfl-master-ai-cloud-engineering/lib/client/_mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestHTTPRepository(t *testing.T) {
	ctrl := gomock.NewController(t)
	client := mocks.NewMockClient(ctrl)

	testUrl, err := url.Parse("http://localhost:3000")
	if err != nil {
		t.Fatal(err)
	}
	repo := NewHTTPRepository(testUrl, client)

	t.Run("VerifyToken", func(t *testing.T) {
		t.Run("Return Error if Request errored", func(t *testing.T) {
			// given
			token := ""
			buf := []byte(`{"token":""}`)

			client.EXPECT().Do(gomock.Any()).
				Do(func(req *http.Request) {
					assert.Equal(t, testUrl, req.URL)
					assert.Equal(t, io.NopCloser(bytes.NewBuffer(buf)), req.Body)
					assert.Equal(t, "POST", req.Method)
				}).
				Return(nil, errors.New("error with request"))

			// when
			userId, err := repo.VerifyToken(token)

			// then
			assert.Error(t, err)
			assert.Equal(t, uint64(0), userId)
		})

		t.Run("Return Error if Response is unauthorized", func(t *testing.T) {
			// given
			token := ""
			buf := []byte(`{"token":""}`)
			response := &http.Response{
				Status:        "401 Unauthorized",
				StatusCode:    http.StatusUnauthorized,
				Header:        http.Header{},
				Body:          http.NoBody,
				ContentLength: 0,
			}

			client.EXPECT().Do(gomock.Any()).
				Do(func(req *http.Request) {
					assert.Equal(t, testUrl, req.URL)
					assert.Equal(t, io.NopCloser(bytes.NewBuffer(buf)), req.Body)
					assert.Equal(t, "POST", req.Method)
				}).
				Return(response, nil)

			// when
			userId, err := repo.VerifyToken(token)

			// then
			assert.Error(t, err)
			assert.Equal(t, uint64(0), userId)
		})

		t.Run("Return Error if Response is not OK", func(t *testing.T) {
			// given
			token := ""
			buf := []byte(`{"token":""}`)
			response := &http.Response{
				Status:        "500 Internal Server Error",
				StatusCode:    http.StatusInternalServerError,
				Header:        http.Header{},
				Body:          http.NoBody,
				ContentLength: 0,
			}

			client.EXPECT().Do(gomock.Any()).
				Do(func(req *http.Request) {
					assert.Equal(t, testUrl, req.URL)
					assert.Equal(t, io.NopCloser(bytes.NewBuffer(buf)), req.Body)
					assert.Equal(t, "POST", req.Method)
				}).
				Return(response, nil)

			// when
			userId, err := repo.VerifyToken(token)

			// then
			assert.Error(t, err)
			assert.Equal(t, uint64(0), userId)
		})

		t.Run("Return Error if Response body is not valid", func(t *testing.T) {
			// given
			token := ""
			buf := []byte(`{"token":""}`)

			responseBodyContent := []byte("invalid json")
			response := &http.Response{
				Status:        "200 OK",
				StatusCode:    http.StatusOK,
				Header:        http.Header{},
				Body:          io.NopCloser(bytes.NewBuffer(responseBodyContent)),
				ContentLength: int64(len(responseBodyContent)),
			}

			client.EXPECT().Do(gomock.Any()).
				Do(func(req *http.Request) {
					assert.Equal(t, testUrl, req.URL)
					assert.Equal(t, io.NopCloser(bytes.NewBuffer(buf)), req.Body)
					assert.Equal(t, "POST", req.Method)
				}).
				Return(response, nil)

			// when
			userId, err := repo.VerifyToken(token)

			// then
			assert.Error(t, err)
			assert.Equal(t, uint64(0), userId)
		})

		t.Run("Return Error if Response body isn't success", func(t *testing.T) {
			// given
			token := ""
			buf := []byte(`{"token":""}`)

			responseBodyContent := []byte(`{"success":false,"userId":1}`)
			response := &http.Response{
				Status:        "200 OK",
				StatusCode:    http.StatusOK,
				Header:        http.Header{},
				Body:          io.NopCloser(bytes.NewBuffer(responseBodyContent)),
				ContentLength: int64(len(responseBodyContent)),
			}

			client.EXPECT().Do(gomock.Any()).
				Do(func(req *http.Request) {
					assert.Equal(t, testUrl, req.URL)
					assert.Equal(t, io.NopCloser(bytes.NewBuffer(buf)), req.Body)
					assert.Equal(t, "POST", req.Method)
				}).
				Return(response, nil)

			// when
			userId, err := repo.VerifyToken(token)

			// then
			assert.Error(t, err)
			assert.Equal(t, uint64(0), userId)
		})

		t.Run("Return UserId if Response body is success", func(t *testing.T) {
			// given
			token := ""
			buf := []byte(`{"token":""}`)

			responseBodyContent := []byte(`{"success":true,"userId":1}`)
			response := &http.Response{
				Status:        "200 OK",
				StatusCode:    http.StatusOK,
				Header:        http.Header{},
				Body:          io.NopCloser(bytes.NewBuffer(responseBodyContent)),
				ContentLength: int64(len(responseBodyContent)),
			}

			client.EXPECT().Do(gomock.Any()).
				Do(func(req *http.Request) {
					assert.Equal(t, testUrl, req.URL)
					assert.Equal(t, io.NopCloser(bytes.NewBuffer(buf)), req.Body)
					assert.Equal(t, "POST", req.Method)
				}).
				Return(response, nil)

			// when
			userId, err := repo.VerifyToken(token)

			// then
			assert.NoError(t, err)
			assert.Equal(t, uint64(1), userId)
		})
	})
}
