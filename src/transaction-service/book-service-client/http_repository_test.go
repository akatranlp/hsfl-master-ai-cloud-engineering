package book_service_client

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

	t.Run("MoveBalance", func(t *testing.T) {
		t.Run("Return Error if Request errored", func(t *testing.T) {
			// given
			userId := uint64(1)
			chapterId := uint64(1)
			bookId := uint64(1)
			buf := []byte(`{"userId":1,"chapterId":1,"bookId":1}`)

			client.EXPECT().Do(gomock.Any()).
				Do(func(req *http.Request) {
					assert.Equal(t, testUrl, req.URL)
					assert.Equal(t, io.NopCloser(bytes.NewBuffer(buf)), req.Body)
					assert.Equal(t, "POST", req.Method)
				}).
				Return(nil, errors.New("error with request"))

			// when
			res, err := repo.ValidateChapterId(userId, chapterId, bookId)

			// then
			assert.Error(t, err)
			assert.Nil(t, res)
		})

		t.Run("Return Error if Response is unauthorized", func(t *testing.T) {
			// given
			userId := uint64(1)
			chapterId := uint64(1)
			bookId := uint64(1)
			buf := []byte(`{"userId":1,"chapterId":1,"bookId":1}`)
			response := &http.Response{
				Status:        "400 Bad Request",
				StatusCode:    http.StatusBadRequest,
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
			res, err := repo.ValidateChapterId(userId, chapterId, bookId)

			// then
			assert.Error(t, err)
			assert.Nil(t, res)
		})

		t.Run("Return Error if Response is not OK", func(t *testing.T) {
			// given
			userId := uint64(1)
			chapterId := uint64(1)
			bookId := uint64(1)
			buf := []byte(`{"userId":1,"chapterId":1,"bookId":1}`)

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

			// then
			res, err := repo.ValidateChapterId(userId, chapterId, bookId)

			// then
			assert.Error(t, err)
			assert.Nil(t, res)
		})

		t.Run("Return Error if Response body is not valid", func(t *testing.T) {
			// given
			userId := uint64(1)
			chapterId := uint64(1)
			bookId := uint64(1)
			buf := []byte(`{"userId":1,"chapterId":1,"bookId":1}`)

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
			res, err := repo.ValidateChapterId(userId, chapterId, bookId)

			// then
			assert.Error(t, err)
			assert.Nil(t, res)
		})

		t.Run("Return UserId if Response body is success", func(t *testing.T) {
			// given
			userId := uint64(1)
			chapterId := uint64(1)
			bookId := uint64(1)
			buf := []byte(`{"userId":1,"chapterId":1,"bookId":1}`)

			responseBodyContent := []byte(`{"chapterId":1,"bookId":1,"receivingUserId":2,"amount":100}`)
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
			res, err := repo.ValidateChapterId(userId, chapterId, bookId)

			// then
			assert.NoError(t, err)
			assert.NotNil(t, res)
			assert.Equal(t, chapterId, res.ChapterId)
			assert.Equal(t, uint64(1), res.BookId)
			assert.Equal(t, uint64(2), res.ReceivingUserId)
			assert.Equal(t, uint64(100), res.Amount)
		})
	})
}
