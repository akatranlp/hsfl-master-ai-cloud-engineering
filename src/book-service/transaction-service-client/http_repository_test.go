package transaction_service_client

import (
	"bytes"
	"errors"
	mocks "github.com/akatranlp/hsfl-master-ai-cloud-engineering/lib/client/_mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"io"
	"net/http"
	"net/url"
	"testing"
)

func TestHTTPRepository(t *testing.T) {
	ctrl := gomock.NewController(t)
	client := mocks.NewMockClient(ctrl)

	testUrl, err := url.Parse("http://localhost:3000")
	if err != nil {
		t.Fatal(err)
	}
	repo := NewHTTPRepository(testUrl, client)

	t.Run("CheckChapterBought", func(t *testing.T) {
		t.Run("Return Error if Request errored", func(t *testing.T) {
			// given
			userId := uint64(1)
			chapterId := uint64(1)
			buf := []byte(`{"userId":1,"chapterId":1}`)

			client.EXPECT().Do(gomock.Any()).
				Do(func(req *http.Request) {
					assert.Equal(t, testUrl, req.URL)
					assert.Equal(t, io.NopCloser(bytes.NewBuffer(buf)), req.Body)
					assert.Equal(t, "POST", req.Method)
				}).
				Return(nil, errors.New("error with request"))

			// when
			err := repo.CheckChapterBought(userId, chapterId)

			// then
			assert.Error(t, err)
		})

		t.Run("Return Error if Response is unauthorized", func(t *testing.T) {
			// given
			userId := uint64(1)
			chapterId := uint64(1)
			buf := []byte(`{"userId":1,"chapterId":1}`)
			response := &http.Response{
				Status:        "404 Not Found",
				StatusCode:    http.StatusNotFound,
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
			err := repo.CheckChapterBought(userId, chapterId)

			// then
			assert.Error(t, err)
		})

		t.Run("Return Error if Response is not OK", func(t *testing.T) {
			// given
			userId := uint64(1)
			chapterId := uint64(1)
			buf := []byte(`{"userId":1,"chapterId":1}`)
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
			err := repo.CheckChapterBought(userId, chapterId)

			// then
			assert.Error(t, err)
		})

		t.Run("Return Error if Response body is not valid", func(t *testing.T) {
			// given
			userId := uint64(1)
			chapterId := uint64(1)
			buf := []byte(`{"userId":1,"chapterId":1}`)

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
			err := repo.CheckChapterBought(userId, chapterId)

			// then
			assert.Error(t, err)
		})

		t.Run("Return Error if Response body isn't success", func(t *testing.T) {
			// given
			userId := uint64(1)
			chapterId := uint64(1)
			buf := []byte(`{"userId":1,"chapterId":1}`)

			responseBodyContent := []byte(`{"success":false}`)
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
			err := repo.CheckChapterBought(userId, chapterId)

			// then
			assert.Error(t, err)
		})

		t.Run("Return UserId if Response body is success", func(t *testing.T) {
			// given
			userId := uint64(1)
			chapterId := uint64(1)
			buf := []byte(`{"userId":1,"chapterId":1}`)

			responseBodyContent := []byte(`{"success":true}`)
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
			err := repo.CheckChapterBought(userId, chapterId)

			// then
			assert.NoError(t, err)
		})
	})
}
