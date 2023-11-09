package books

import (
	"context"
	"encoding/json"
	"errors"
	books_mocks "github.com/akatranlp/hsfl-master-ai-cloud-engineering/book-service/_mocks/books"
	"github.com/akatranlp/hsfl-master-ai-cloud-engineering/book-service/books/model"
	authMiddleware "github.com/akatranlp/hsfl-master-ai-cloud-engineering/lib/auth-middleware"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestBookDefaultController(t *testing.T) {
	ctrl := gomock.NewController(t)
	bookRepository := books_mocks.NewMockRepository(ctrl)
	controller := NewDefaultController(bookRepository)

	t.Run("GetBooks", func(t *testing.T) {
		t.Run("should return 500 INTERNAL SERVER ERROR if query failed", func(t *testing.T) {
			// given
			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", "/api/v1/books", nil)

			bookRepository.
				EXPECT().
				FindAll().
				Return(nil, errors.New("query failed")).
				Times(1)

			// when
			controller.GetBooks(w, r)

			// then
			assert.Equal(t, http.StatusInternalServerError, w.Code)
		})

		t.Run("should return all books", func(t *testing.T) {
			// given
			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", "/api/v1/books", nil)

			bookRepository.
				EXPECT().
				FindAll().
				Return([]*model.Book{{ID: 999}}, nil).
				Times(1)

			// when
			controller.GetBooks(w, r)

			// then
			res := w.Result()
			var response []model.Book
			err := json.NewDecoder(res.Body).Decode(&response)

			assert.NoError(t, err)
			assert.Equal(t, http.StatusOK, w.Code)
			assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
			assert.Len(t, response, 1)
			assert.Equal(t, uint64(999), response[0].ID)
		})

		t.Run("should return 500 INTERNAL SERVER ERROR if userId cant be parsed failed", func(t *testing.T) {
			// given
			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", "/api/v1/books?userId=invalid", nil)

			// when
			controller.GetBooks(w, r)

			// then
			assert.Equal(t, http.StatusBadRequest, w.Code)
		})

		t.Run("should return all your books", func(t *testing.T) {
			// given
			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", "/api/v1/books?userId=1", nil)

			bookRepository.
				EXPECT().
				FindAllByUserId(uint64(1)).
				Return([]*model.Book{{ID: 999}}, nil).
				Times(1)

			// when
			controller.GetBooks(w, r)

			// then
			res := w.Result()
			var response []model.Book
			err := json.NewDecoder(res.Body).Decode(&response)

			assert.NoError(t, err)
			assert.Equal(t, http.StatusOK, w.Code)
			assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
			assert.Len(t, response, 1)
			assert.Equal(t, uint64(999), response[0].ID)
		})
	})

	t.Run("PostBooks", func(t *testing.T) {
		t.Run("should return 400 BAD REQUEST if payload is not json", func(t *testing.T) {
			tests := []io.Reader{
				nil,
				strings.NewReader(`{"invalid`),
			}

			for _, test := range tests {
				// given
				w := httptest.NewRecorder()
				r := httptest.NewRequest("POST", "/api/v1/books", test)
				r = r.WithContext(context.WithValue(r.Context(), authMiddleware.AuthenticatedUserId, uint64(1)))

				// when
				controller.PostBook(w, r)

				// then
				assert.Equal(t, http.StatusBadRequest, w.Code)
			}
		})

		t.Run("should return 400 BAD REQUEST if payload is incomplete", func(t *testing.T) {
			tests := []io.Reader{
				strings.NewReader(`{"description": "amazing book"}`),
				strings.NewReader(`{"authorid": 1}`),
			}

			for _, test := range tests {
				// given
				w := httptest.NewRecorder()
				r := httptest.NewRequest("POST", "/api/v1/books", test)
				r = r.WithContext(context.WithValue(r.Context(), authMiddleware.AuthenticatedUserId, uint64(1)))

				// when
				controller.PostBook(w, r)

				// then
				assert.Equal(t, http.StatusBadRequest, w.Code)
			}
		})

		t.Run("should return 500 INTERNAL SERVER ERROR if persisting failed", func(t *testing.T) {
			// given
			w := httptest.NewRecorder()
			r := httptest.NewRequest("POST", "/api/v1/books",
				strings.NewReader(`{"name":"test book","description":"amazing book"}`))
			userId := uint64(1)
			r = r.WithContext(context.WithValue(r.Context(), authMiddleware.AuthenticatedUserId, userId))

			bookRepository.
				EXPECT().
				Create([]*model.Book{{Name: "test book", AuthorID: userId, Description: "amazing book"}}).
				Return(errors.New("database error"))

			// when
			controller.PostBook(w, r)

			// then
			assert.Equal(t, http.StatusInternalServerError, w.Code)
		})

		t.Run("should create new book", func(t *testing.T) {
			// given
			w := httptest.NewRecorder()
			r := httptest.NewRequest("POST", "/api/v1/books",
				strings.NewReader(`{"name":"test book","description":"amazing book"}`))
			userId := uint64(1)
			r = r.WithContext(context.WithValue(r.Context(), authMiddleware.AuthenticatedUserId, userId))

			bookRepository.
				EXPECT().
				Create([]*model.Book{{Name: "test book", AuthorID: userId, Description: "amazing book"}}).
				Return(nil)

			// when
			controller.PostBook(w, r)

			// then
			assert.Equal(t, http.StatusOK, w.Code)
		})
	})

	t.Run("GetBook", func(t *testing.T) {
		t.Run("should return 200 OK and book", func(t *testing.T) {
			// given
			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", "/api/v1/books/1", nil)
			dbBook := &model.Book{
				ID:          1,
				Name:        "Book One",
				AuthorID:    1,
				Description: "! good book",
			}
			r = r.WithContext(context.WithValue(r.Context(), MiddleWareBook, dbBook))

			// when
			controller.GetBook(w, r)

			// then
			res := w.Result()
			var response model.Book
			err := json.NewDecoder(res.Body).Decode(&response)

			assert.NoError(t, err)
			assert.Equal(t, http.StatusOK, w.Code)
			assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
			assert.Equal(t, dbBook, &response)
		})
	})

	t.Run("Patch", func(t *testing.T) {
		t.Run("should return 400 BAD REQUEST if payload is not json", func(t *testing.T) {
			tests := []io.Reader{
				nil,
				strings.NewReader(`{"invalid`),
			}

			for _, test := range tests {
				// given
				w := httptest.NewRecorder()
				r := httptest.NewRequest("PATCH", "/api/v1/books/1", test)
				dbBook := &model.Book{
					ID:          1,
					Name:        "Book One",
					AuthorID:    1,
					Description: "! good book",
				}
				r = r.WithContext(context.WithValue(r.Context(), authMiddleware.AuthenticatedUserId, uint64(1)))
				r = r.WithContext(context.WithValue(r.Context(), MiddleWareBook, dbBook))

				// when
				controller.PatchBook(w, r)

				// then
				assert.Equal(t, http.StatusBadRequest, w.Code)
			}
		})

		t.Run("should return 500 INTERNAL SERVER ERROR if query failed", func(t *testing.T) {
			// given
			w := httptest.NewRecorder()
			r := httptest.NewRequest("PATCH", "/api/v1/books/1",
				strings.NewReader(`{"id": 999}`))
			dbBook := &model.Book{
				ID:          1,
				Name:        "Book One",
				AuthorID:    1,
				Description: "! good book",
			}
			r = r.WithContext(context.WithValue(r.Context(), authMiddleware.AuthenticatedUserId, uint64(1)))
			r = r.WithContext(context.WithValue(r.Context(), MiddleWareBook, dbBook))

			bookRepository.
				EXPECT().
				Update(uint64(1), &model.BookPatch{}).
				Return(errors.New("database error"))

			// when
			controller.PatchBook(w, r)

			// then
			assert.Equal(t, http.StatusInternalServerError, w.Code)
		})

		t.Run("should return 500 INTERNAL SERVER ERROR if query failed", func(t *testing.T) {
			// given
			w := httptest.NewRecorder()
			r := httptest.NewRequest("PUT", "/api/v1/books/1",
				strings.NewReader(`{"id": 999}`))
			dbBook := &model.Book{
				ID:          1,
				Name:        "Book One",
				AuthorID:    1,
				Description: "! good book",
			}
			r = r.WithContext(context.WithValue(r.Context(), authMiddleware.AuthenticatedUserId, uint64(1)))
			r = r.WithContext(context.WithValue(r.Context(), MiddleWareBook, dbBook))

			bookRepository.
				EXPECT().
				Update(uint64(1), &model.BookPatch{}).
				Return(errors.New("database error"))

			// when
			controller.PatchBook(w, r)

			// then
			assert.Equal(t, http.StatusInternalServerError, w.Code)
		})

		t.Run("should return 401 If you are not the creator of the book", func(t *testing.T) {
			// given
			w := httptest.NewRecorder()
			r := httptest.NewRequest("PUT", "/api/v1/books/1",
				strings.NewReader(`{"id": 999}`))
			dbBook := &model.Book{
				ID:          1,
				Name:        "Book One",
				AuthorID:    1,
				Description: "! good book",
			}
			r = r.WithContext(context.WithValue(r.Context(), authMiddleware.AuthenticatedUserId, uint64(2)))
			r = r.WithContext(context.WithValue(r.Context(), MiddleWareBook, dbBook))

			// when
			controller.PatchBook(w, r)

			// then
			assert.Equal(t, http.StatusUnauthorized, w.Code)
		})

		t.Run("should update one book", func(t *testing.T) {
			// given
			w := httptest.NewRecorder()
			r := httptest.NewRequest("PUT", "/api/v1/books/1",
				strings.NewReader(`{"description": "a fine book"}`))
			dbBook := &model.Book{
				ID:          1,
				Name:        "Book One",
				AuthorID:    1,
				Description: "! good book",
			}
			r = r.WithContext(context.WithValue(r.Context(), authMiddleware.AuthenticatedUserId, uint64(1)))
			r = r.WithContext(context.WithValue(r.Context(), MiddleWareBook, dbBook))

			newDescription := "a fine book"
			bookRepository.
				EXPECT().
				Update(uint64(1), &model.BookPatch{Description: &newDescription}).
				Return(nil)

			// when
			controller.PatchBook(w, r)

			// then
			assert.Equal(t, http.StatusOK, w.Code)
		})
	})

	t.Run("DeleteBook", func(t *testing.T) {
		t.Run("should return 500 INTERNAL SERVER ERROR if query fails", func(t *testing.T) {
			// given
			w := httptest.NewRecorder()
			r := httptest.NewRequest("DELETE", "/api/v1/books/1", nil)
			dbBook := &model.Book{
				ID:          1,
				Name:        "Book One",
				AuthorID:    1,
				Description: "! good book",
			}
			r = r.WithContext(context.WithValue(r.Context(), authMiddleware.AuthenticatedUserId, uint64(1)))
			r = r.WithContext(context.WithValue(r.Context(), MiddleWareBook, dbBook))

			bookRepository.
				EXPECT().
				Delete([]*model.Book{dbBook}).
				Return(errors.New("database error"))

			// when
			controller.DeleteBook(w, r)

			// then
			assert.Equal(t, http.StatusInternalServerError, w.Code)
		})

		t.Run("should return 401 if not the user who created the book", func(t *testing.T) {
			// given
			w := httptest.NewRecorder()
			r := httptest.NewRequest("DELETE", "/api/v1/books/1", nil)
			dbBook := &model.Book{
				ID:          1,
				Name:        "Book One",
				AuthorID:    2,
				Description: "! good book",
			}
			r = r.WithContext(context.WithValue(r.Context(), authMiddleware.AuthenticatedUserId, uint64(1)))
			r = r.WithContext(context.WithValue(r.Context(), MiddleWareBook, dbBook))

			// when
			controller.DeleteBook(w, r)

			// then
			assert.Equal(t, http.StatusUnauthorized, w.Code)
		})

		t.Run("should return 200 OK", func(t *testing.T) {
			// given
			w := httptest.NewRecorder()
			r := httptest.NewRequest("DELETE", "/api/v1/books/1", nil)
			dbBook := &model.Book{
				ID:          1,
				Name:        "Book One",
				AuthorID:    1,
				Description: "! good book",
			}
			r = r.WithContext(context.WithValue(r.Context(), authMiddleware.AuthenticatedUserId, uint64(1)))
			r = r.WithContext(context.WithValue(r.Context(), MiddleWareBook, dbBook))

			bookRepository.
				EXPECT().
				Delete([]*model.Book{dbBook}).
				Return(nil)

			// when
			controller.DeleteBook(w, r)

			// then
			assert.Equal(t, http.StatusOK, w.Code)
		})
	})

	t.Run("LoadBookMiddleware", func(t *testing.T) {
		t.Run("Should return 400 if the bookid cannot be parsed", func(t *testing.T) {
			// given
			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", "/api/v1/books/aaa", nil)
			r = r.WithContext(context.WithValue(r.Context(), "bookid", "aaa"))

			// when
			called := false
			controller.LoadBookMiddleware(w, r, func(r *http.Request) {
				called = true
			})

			assert.Equal(t, false, called)
			assert.Equal(t, http.StatusBadRequest, w.Code)
		})

		t.Run("Should return 404 if the query fails", func(t *testing.T) {
			// given
			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", "/api/v1/books/1", nil)
			r = r.WithContext(context.WithValue(r.Context(), "bookid", "1"))

			bookRepository.
				EXPECT().
				FindById(uint64(1)).
				Return(nil, errors.New("database error"))

			// when
			called := false
			controller.LoadBookMiddleware(w, r, func(r *http.Request) {
				called = true
			})

			assert.Equal(t, false, called)
			assert.Equal(t, http.StatusNotFound, w.Code)
		})

		t.Run("Should return 200 if it succeeds", func(t *testing.T) {
			// given
			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", "/api/v1/books/1", nil)
			r = r.WithContext(context.WithValue(r.Context(), "bookid", "1"))
			dbBook := &model.Book{
				ID:          1,
				Name:        "Book One",
				AuthorID:    1,
				Description: "! good book",
			}

			bookRepository.
				EXPECT().
				FindById(uint64(1)).
				Return(dbBook, nil)

			// when
			called := false
			controller.LoadBookMiddleware(w, r, func(req *http.Request) {
				called = true
				r = req
			})

			assert.Equal(t, true, called)
			assert.Equal(t, dbBook, r.Context().Value(MiddleWareBook))
			assert.Equal(t, http.StatusOK, w.Code)
		})
	})
}
