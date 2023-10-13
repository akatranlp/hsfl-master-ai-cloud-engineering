package router

import (
	"context"
	"fmt"
	mocks "github.com/akatranlp/hsfl-master-ai-cloud-engineering/book-service/_mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRouter(t *testing.T) {
	ctrl := gomock.NewController(t)

	booksController := mocks.NewMockController(ctrl)
	router := New(booksController)

	t.Run("/api/v1/books", func(t *testing.T) {
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
			r := httptest.NewRequest("GET", "/api/v1/books", nil)

			fmt.Println("booksController")
			fmt.Println(w, r)
			booksController.
				EXPECT().
				GetBooks(w, r).
				Times(1)

			// when
			router.ServeHTTP(w, r)

			// then
			assert.Equal(t, http.StatusOK, w.Code)
		})

		t.Run("should call POST handler", func(t *testing.T) {
			// given
			w := httptest.NewRecorder()
			r := httptest.NewRequest("POST", "/api/v1/books", nil)
			
			booksController.
				EXPECT().
				PostBooks(w, r).
				Times(1)

			// when
			router.ServeHTTP(w, r)

			// then
			assert.Equal(t, http.StatusOK, w.Code)
		})
	})

	t.Run("/api/v1/books/:bookid", func(t *testing.T) {
		t.Run("should return 404 NOT FOUND if method is not GET, DELETE or PUT", func(t *testing.T) {
			tests := []string{"POST", "HEAD", "CONNECT", "OPTIONS", "TRACE", "PATCH"}

			for _, test := range tests {
				// given
				w := httptest.NewRecorder()
				r := httptest.NewRequest(test, "/api/v1/books/1", nil)

				// when
				router.ServeHTTP(w, r)

				// then
				assert.Equal(t, http.StatusNotFound, w.Code)
			}
		})

		t.Run("should call GET handler", func(t *testing.T) {
			// given
			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", "/api/v1/books/1", nil)

			booksController.
				EXPECT().
				GetBook(w, r.WithContext(context.WithValue(r.Context(), "bookid", "1"))).
				Times(1)

			// when
			router.ServeHTTP(w, r)

			// then
			assert.Equal(t, http.StatusOK, w.Code)
		})

		t.Run("should call PUT handler", func(t *testing.T) {
			// given
			w := httptest.NewRecorder()
			r := httptest.NewRequest("PUT", "/api/v1/books/1", nil)

			booksController.
				EXPECT().
				PutBook(w, r.WithContext(context.WithValue(r.Context(), "bookid", "1"))).
				Times(1)

			// when
			router.ServeHTTP(w, r)

			// then
			assert.Equal(t, http.StatusOK, w.Code)
		})

		t.Run("should call DELETE handler", func(t *testing.T) {
			// given
			w := httptest.NewRecorder()
			r := httptest.NewRequest("DELETE", "/api/v1/books/1", nil)

			booksController.
				EXPECT().
				DeleteBook(w, r.WithContext(context.WithValue(r.Context(), "bookid", "1"))).
				Times(1)

			// when
			router.ServeHTTP(w, r)

			// then
			assert.Equal(t, http.StatusOK, w.Code)
		})
	})
}
