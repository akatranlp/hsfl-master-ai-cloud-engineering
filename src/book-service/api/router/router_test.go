package router

import (
	"context"
	books_mocks "github.com/akatranlp/hsfl-master-ai-cloud-engineering/book-service/_mocks/books"
	chapters_mocks "github.com/akatranlp/hsfl-master-ai-cloud-engineering/book-service/_mocks/chapters"
	"github.com/akatranlp/hsfl-master-ai-cloud-engineering/lib/auth-middleware/_mocks"
	libRouter "github.com/akatranlp/hsfl-master-ai-cloud-engineering/lib/router"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRouter(t *testing.T) {
	ctrl := gomock.NewController(t)

	booksController := books_mocks.NewMockController(ctrl)
	chaptersController := chapters_mocks.NewMockController(ctrl)
	authController := mocks.NewMockController(ctrl)
	router := New(authController, booksController, chaptersController)

	t.Run("auth /api/v1/books", func(t *testing.T) {
		t.Run("Get Books should not be called", func(t *testing.T) {
			// given
			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", "/api/v1/books", nil)

			authController.
				EXPECT().
				AuthenticationMiddleware(w, r, gomock.Any()).
				Do(func(w http.ResponseWriter, r *http.Request, next libRouter.Next) {
					w.WriteHeader(http.StatusUnauthorized)
				}).Times(1)

			// when
			router.ServeHTTP(w, r)

			// then
			assert.Equal(t, http.StatusUnauthorized, w.Code)
		})

		t.Run("Get Books should be called", func(t *testing.T) {
			// given
			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", "/api/v1/books", nil)

			authController.
				EXPECT().
				AuthenticationMiddleware(w, r, gomock.Any()).
				Do(func(w http.ResponseWriter, r *http.Request, next libRouter.Next) {
					next(r)
				}).Times(1)

			booksController.
				EXPECT().
				GetBooks(w, r).
				Times(1)

			// when
			router.ServeHTTP(w, r)

			// then
			assert.Equal(t, http.StatusOK, w.Code)
		})
	})

	t.Run("middleware /api/v1/books/:bookid", func(t *testing.T) {
		t.Run("Get Book should not be called", func(t *testing.T) {
			// given
			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", "/api/v1/books/1", nil)

			authController.
				EXPECT().
				AuthenticationMiddleware(w, r, gomock.Any()).
				Do(func(w http.ResponseWriter, r *http.Request, next libRouter.Next) {
					next(r)
				}).Times(1)

			booksController.
				EXPECT().
				LoadBookMiddleware(w, r.WithContext(context.WithValue(r.Context(), "bookid", "1")), gomock.Any()).
				Do(func(w http.ResponseWriter, r *http.Request, next libRouter.Next) {
					w.WriteHeader(http.StatusNotFound)
				}).Times(1)

			// when
			router.ServeHTTP(w, r)

			// then
			assert.Equal(t, http.StatusNotFound, w.Code)
		})

		t.Run("Get Book should not be called", func(t *testing.T) {
			// given
			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", "/api/v1/books/1", nil)

			authController.
				EXPECT().
				AuthenticationMiddleware(w, r, gomock.Any()).
				Do(func(w http.ResponseWriter, r *http.Request, next libRouter.Next) {
					next(r)
				}).Times(1)

			req := r.WithContext(context.WithValue(r.Context(), "bookid", "1"))
			booksController.
				EXPECT().
				LoadBookMiddleware(w, req, gomock.Any()).
				Do(func(w http.ResponseWriter, r *http.Request, next libRouter.Next) {
					next(r)
				}).Times(1)

			booksController.EXPECT().
				GetBook(w, req).
				Times(1)

			// when
			router.ServeHTTP(w, r)

			// then
			assert.Equal(t, http.StatusOK, w.Code)
		})
	})

	t.Run("/api/v1/books", func(t *testing.T) {
		t.Run("should return 404 NOT FOUND if method is not GET or POST", func(t *testing.T) {
			tests := []string{"DELETE", "PUT", "HEAD", "CONNECT", "OPTIONS", "TRACE", "PATCH"}
			for _, test := range tests {
				// given
				w := httptest.NewRecorder()
				r := httptest.NewRequest(test, "/api/v1/books", nil)

				authController.
					EXPECT().
					AuthenticationMiddleware(w, r, gomock.Any()).
					Do(func(w http.ResponseWriter, r *http.Request, next libRouter.Next) {
						next(r)
					}).Times(1)

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

			authController.
				EXPECT().
				AuthenticationMiddleware(w, r, gomock.Any()).
				Do(func(w http.ResponseWriter, r *http.Request, next libRouter.Next) {
					next(r)
				}).Times(1)

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

			authController.
				EXPECT().
				AuthenticationMiddleware(w, r, gomock.Any()).
				Do(func(w http.ResponseWriter, r *http.Request, next libRouter.Next) {
					next(r)
				}).Times(1)

			booksController.
				EXPECT().
				PostBook(w, r).
				Times(1)

			// when
			router.ServeHTTP(w, r)

			// then
			assert.Equal(t, http.StatusOK, w.Code)
		})
	})

	t.Run("/api/v1/books/:bookid", func(t *testing.T) {
		t.Run("should return 404 NOT FOUND if method is not GET, DELETE or PATCH", func(t *testing.T) {
			tests := []string{"POST", "HEAD", "CONNECT", "OPTIONS", "TRACE", "PUT"}

			for _, test := range tests {
				// given
				w := httptest.NewRecorder()
				r := httptest.NewRequest(test, "/api/v1/books/1", nil)

				authController.
					EXPECT().
					AuthenticationMiddleware(w, r, gomock.Any()).
					Do(func(w http.ResponseWriter, r *http.Request, next libRouter.Next) {
						next(r)
					}).Times(1)

				booksController.
					EXPECT().
					LoadBookMiddleware(w, r.WithContext(context.WithValue(r.Context(), "bookid", "1")), gomock.Any()).
					Do(func(w http.ResponseWriter, r *http.Request, next libRouter.Next) {
						next(r)
					}).Times(1)

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

			authController.
				EXPECT().
				AuthenticationMiddleware(w, r, gomock.Any()).
				Do(func(w http.ResponseWriter, r *http.Request, next libRouter.Next) {
					next(r)
				}).Times(1)

			booksController.
				EXPECT().
				LoadBookMiddleware(w, r.WithContext(context.WithValue(r.Context(), "bookid", "1")), gomock.Any()).
				Do(func(w http.ResponseWriter, r *http.Request, next libRouter.Next) {
					next(r)
				}).Times(1)

			booksController.
				EXPECT().
				GetBook(w, r.WithContext(context.WithValue(r.Context(), "bookid", "1"))).
				Times(1)

			// when
			router.ServeHTTP(w, r)

			// then
			assert.Equal(t, http.StatusOK, w.Code)
		})

		t.Run("should call PATCH handler", func(t *testing.T) {
			// given
			w := httptest.NewRecorder()
			r := httptest.NewRequest("PATCH", "/api/v1/books/1", nil)

			authController.
				EXPECT().
				AuthenticationMiddleware(w, r, gomock.Any()).
				Do(func(w http.ResponseWriter, r *http.Request, next libRouter.Next) {
					next(r)
				}).Times(1)

			booksController.
				EXPECT().
				LoadBookMiddleware(w, r.WithContext(context.WithValue(r.Context(), "bookid", "1")), gomock.Any()).
				Do(func(w http.ResponseWriter, r *http.Request, next libRouter.Next) {
					next(r)
				}).Times(1)

			booksController.
				EXPECT().
				PatchBook(w, r.WithContext(context.WithValue(r.Context(), "bookid", "1"))).
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

			authController.
				EXPECT().
				AuthenticationMiddleware(w, r, gomock.Any()).
				Do(func(w http.ResponseWriter, r *http.Request, next libRouter.Next) {
					next(r)
				}).Times(1)

			booksController.
				EXPECT().
				LoadBookMiddleware(w, r.WithContext(context.WithValue(r.Context(), "bookid", "1")), gomock.Any()).
				Do(func(w http.ResponseWriter, r *http.Request, next libRouter.Next) {
					next(r)
				}).Times(1)

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
