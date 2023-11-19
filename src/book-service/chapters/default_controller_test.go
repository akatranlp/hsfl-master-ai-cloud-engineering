package chapters

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	chapters_mocks "github.com/akatranlp/hsfl-master-ai-cloud-engineering/book-service/_mocks/chapters"
	transaction_service_client_mocks "github.com/akatranlp/hsfl-master-ai-cloud-engineering/book-service/_mocks/transaction-service-client"
	"github.com/akatranlp/hsfl-master-ai-cloud-engineering/book-service/books"
	booksModel "github.com/akatranlp/hsfl-master-ai-cloud-engineering/book-service/books/model"
	"github.com/akatranlp/hsfl-master-ai-cloud-engineering/book-service/chapters/model"
	authMiddleware "github.com/akatranlp/hsfl-master-ai-cloud-engineering/lib/auth-middleware"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestChapterDefaultController(t *testing.T) {
	ctrl := gomock.NewController(t)
	chapterRepository := chapters_mocks.NewMockRepository(ctrl)
	transactionServiceClient := transaction_service_client_mocks.NewMockRepository(ctrl)
	controller := NewDefaultController(chapterRepository, transactionServiceClient)

	t.Run("GetChapters", func(t *testing.T) {
		t.Run("should return 500 INTERNAL SERVER ERROR if query failed", func(t *testing.T) {
			// given
			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", "/api/v1/chapters", nil)
			dbBook := &booksModel.Book{
				ID:          1,
				Name:        "Book One",
				AuthorID:    1,
				Description: "! good book",
			}
			r = r.WithContext(context.WithValue(r.Context(), books.MiddleWareBook, dbBook))

			chapterRepository.
				EXPECT().
				FindAllPreviewsByBookId(dbBook.ID).
				Return(nil, errors.New("query failed")).
				Times(1)

			// when
			controller.GetChaptersForBook(w, r)

			// then
			assert.Equal(t, http.StatusInternalServerError, w.Code)
		})

		t.Run("should return all chapters", func(t *testing.T) {
			// given
			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", "/api/v1/chapters", nil)
			dbBook := &booksModel.Book{
				ID:          1,
				Name:        "Book One",
				AuthorID:    1,
				Description: "! good book",
			}
			r = r.WithContext(context.WithValue(r.Context(), books.MiddleWareBook, dbBook))

			chapterRepository.
				EXPECT().
				FindAllPreviewsByBookId(uint64(1)).
				Return([]*model.ChapterPreview{{ID: 999}}, nil).
				Times(1)

			// when
			controller.GetChaptersForBook(w, r)

			// then
			res := w.Result()
			var response []model.ChapterPreview
			err := json.NewDecoder(res.Body).Decode(&response)

			assert.NoError(t, err)
			assert.Equal(t, http.StatusOK, w.Code)
			assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
			assert.Len(t, response, 1)
			assert.Equal(t, uint64(999), response[0].ID)
		})
	})

	t.Run("PostChapters", func(t *testing.T) {
		t.Run("should return 401 If yoa not the author of the book", func(t *testing.T) {

			// given
			w := httptest.NewRecorder()
			r := httptest.NewRequest("POST", "/api/v1/chapters", nil)
			r = r.WithContext(context.WithValue(r.Context(), authMiddleware.AuthenticatedUserId, uint64(1)))
			dbBook := &booksModel.Book{
				ID:          1,
				Name:        "Book One",
				AuthorID:    2,
				Description: "! good book",
			}
			r = r.WithContext(context.WithValue(r.Context(), books.MiddleWareBook, dbBook))

			// when
			controller.PostChapter(w, r)

			// then
			assert.Equal(t, http.StatusUnauthorized, w.Code)

		})

		t.Run("should return 400 BAD REQUEST if payload is not json", func(t *testing.T) {
			tests := []io.Reader{
				nil,
				strings.NewReader(`{"invalid`),
			}

			for _, test := range tests {
				// given
				w := httptest.NewRecorder()
				r := httptest.NewRequest("POST", "/api/v1/chapters", test)
				r = r.WithContext(context.WithValue(r.Context(), authMiddleware.AuthenticatedUserId, uint64(1)))
				dbBook := &booksModel.Book{
					ID:          1,
					Name:        "Book One",
					AuthorID:    1,
					Description: "! good book",
				}
				r = r.WithContext(context.WithValue(r.Context(), books.MiddleWareBook, dbBook))

				// when
				controller.PostChapter(w, r)

				// then
				assert.Equal(t, http.StatusBadRequest, w.Code)
			}
		})

		t.Run("should return 400 BAD REQUEST if payload is incomplete", func(t *testing.T) {
			tests := []io.Reader{
				strings.NewReader(`{"description": "amazing chapter"}`),
				strings.NewReader(`{"authorid": 1}`),
			}

			for _, test := range tests {
				// given
				w := httptest.NewRecorder()
				r := httptest.NewRequest("POST", "/api/v1/chapters", test)
				r = r.WithContext(context.WithValue(r.Context(), authMiddleware.AuthenticatedUserId, uint64(1)))
				dbBook := &booksModel.Book{
					ID:          1,
					Name:        "Book One",
					AuthorID:    1,
					Description: "! good book",
				}
				r = r.WithContext(context.WithValue(r.Context(), books.MiddleWareBook, dbBook))

				// when
				controller.PostChapter(w, r)

				// then
				assert.Equal(t, http.StatusBadRequest, w.Code)
			}
		})

		t.Run("should return 500 INTERNAL SERVER ERROR if persisting failed", func(t *testing.T) {
			// given
			w := httptest.NewRecorder()
			r := httptest.NewRequest("POST", "/api/v1/chapters",
				strings.NewReader(`{"name":"test chapter","price":10,"content":"amazing chapter"}`))
			userId := uint64(1)
			r = r.WithContext(context.WithValue(r.Context(), authMiddleware.AuthenticatedUserId, userId))
			dbBook := &booksModel.Book{
				ID:          1,
				Name:        "Book One",
				AuthorID:    1,
				Description: "! good book",
			}
			r = r.WithContext(context.WithValue(r.Context(), books.MiddleWareBook, dbBook))

			chapterRepository.
				EXPECT().
				Create([]*model.Chapter{{
					BookID:  1,
					Name:    "test chapter",
					Price:   10,
					Content: "amazing chapter",
				}}).
				Return(errors.New("database error"))

			// when
			controller.PostChapter(w, r)

			// then
			assert.Equal(t, http.StatusInternalServerError, w.Code)
		})

		t.Run("should create new chapter", func(t *testing.T) {
			// given
			w := httptest.NewRecorder()
			r := httptest.NewRequest("POST", "/api/v1/chapters",
				strings.NewReader(`{"name":"test chapter","price":10,"content":"amazing chapter"}`))
			userId := uint64(1)
			r = r.WithContext(context.WithValue(r.Context(), authMiddleware.AuthenticatedUserId, userId))
			dbBook := &booksModel.Book{
				ID:          1,
				Name:        "Book One",
				AuthorID:    1,
				Description: "! good book",
			}
			r = r.WithContext(context.WithValue(r.Context(), books.MiddleWareBook, dbBook))

			chapterRepository.
				EXPECT().
				Create([]*model.Chapter{{
					BookID:  1,
					Name:    "test chapter",
					Price:   10,
					Content: "amazing chapter",
				}}).
				Return(nil)

			// when
			controller.PostChapter(w, r)

			// then
			assert.Equal(t, http.StatusOK, w.Code)
		})
	})

	t.Run("GetChapter", func(t *testing.T) {
		t.Run("should return 200 OK and chapter", func(t *testing.T) {
			// given
			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", "/api/v1/chapters/1", nil)
			r = r.WithContext(context.WithValue(r.Context(), authMiddleware.AuthenticatedUserId, uint64(1)))
			dbBook := &booksModel.Book{
				ID:          1,
				Name:        "Book One",
				AuthorID:    1,
				Description: "! good book",
			}
			r = r.WithContext(context.WithValue(r.Context(), books.MiddleWareBook, dbBook))
			dbChapter := &model.Chapter{
				ID:      1,
				BookID:  1,
				Name:    "Chapter One",
				Price:   100,
				Content: "Nice chapter",
			}
			r = r.WithContext(context.WithValue(r.Context(), middleWareChapter, dbChapter))

			// when
			controller.GetChapterForBook(w, r)

			// then
			res := w.Result()
			var response model.Chapter
			err := json.NewDecoder(res.Body).Decode(&response)

			assert.NoError(t, err)
			assert.Equal(t, http.StatusOK, w.Code)
			assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
			assert.Equal(t, dbChapter, &response)
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
				r := httptest.NewRequest("PATCH", "/api/v1/chapters/1", test)
				dbBook := &booksModel.Book{
					ID:          1,
					Name:        "Book One",
					AuthorID:    1,
					Description: "! good book",
				}
				r = r.WithContext(context.WithValue(r.Context(), books.MiddleWareBook, dbBook))
				dbChapter := &model.Chapter{
					ID:      1,
					BookID:  1,
					Name:    "Chapter One",
					Price:   100,
					Content: "Nice chapter",
				}
				r = r.WithContext(context.WithValue(r.Context(), authMiddleware.AuthenticatedUserId, uint64(1)))
				r = r.WithContext(context.WithValue(r.Context(), middleWareChapter, dbChapter))

				// when
				controller.PatchChapter(w, r)

				// then
				assert.Equal(t, http.StatusBadRequest, w.Code)
			}
		})

		t.Run("should return 500 INTERNAL SERVER ERROR if query failed", func(t *testing.T) {
			// given
			w := httptest.NewRecorder()
			r := httptest.NewRequest("PUT", "/api/v1/chapters/1",
				strings.NewReader(`{"id": 999}`))
			dbBook := &booksModel.Book{
				ID:          1,
				Name:        "Book One",
				AuthorID:    1,
				Description: "! good book",
			}
			r = r.WithContext(context.WithValue(r.Context(), books.MiddleWareBook, dbBook))
			dbChapter := &model.Chapter{
				ID:      1,
				BookID:  1,
				Name:    "Chapter One",
				Price:   100,
				Content: "Nice chapter",
			}
			r = r.WithContext(context.WithValue(r.Context(), authMiddleware.AuthenticatedUserId, uint64(1)))
			r = r.WithContext(context.WithValue(r.Context(), middleWareChapter, dbChapter))

			chapterRepository.
				EXPECT().
				Update(uint64(1), &model.ChapterPatch{}).
				Return(errors.New("database error"))

			// when
			controller.PatchChapter(w, r)

			// then
			assert.Equal(t, http.StatusInternalServerError, w.Code)
		})

		t.Run("should return 500 INTERNAL SERVER ERROR if query failed", func(t *testing.T) {
			// given
			w := httptest.NewRecorder()
			r := httptest.NewRequest("PUT", "/api/v1/chapters/1",
				strings.NewReader(`{"id": 999}`))
			dbBook := &booksModel.Book{
				ID:          1,
				Name:        "Book One",
				AuthorID:    1,
				Description: "! good book",
			}
			r = r.WithContext(context.WithValue(r.Context(), books.MiddleWareBook, dbBook))
			dbChapter := &model.Chapter{
				ID:      1,
				BookID:  1,
				Name:    "Chapter One",
				Price:   100,
				Content: "Nice chapter",
			}
			r = r.WithContext(context.WithValue(r.Context(), authMiddleware.AuthenticatedUserId, uint64(1)))
			r = r.WithContext(context.WithValue(r.Context(), middleWareChapter, dbChapter))

			chapterRepository.
				EXPECT().
				Update(uint64(1), &model.ChapterPatch{}).
				Return(errors.New("database error"))

			// when
			controller.PatchChapter(w, r)

			// then
			assert.Equal(t, http.StatusInternalServerError, w.Code)
		})

		t.Run("should return 401 If you are not the creator of the chapter", func(t *testing.T) {
			// given
			w := httptest.NewRecorder()
			r := httptest.NewRequest("PUT", "/api/v1/chapters/1",
				strings.NewReader(`{"id": 999}`))
			dbBook := &booksModel.Book{
				ID:          1,
				Name:        "Book One",
				AuthorID:    1,
				Description: "! good book",
			}
			r = r.WithContext(context.WithValue(r.Context(), books.MiddleWareBook, dbBook))
			dbChapter := &model.Chapter{
				ID:      1,
				BookID:  1,
				Name:    "Chapter One",
				Price:   100,
				Content: "Nice chapter",
			}
			r = r.WithContext(context.WithValue(r.Context(), authMiddleware.AuthenticatedUserId, uint64(2)))
			r = r.WithContext(context.WithValue(r.Context(), middleWareChapter, dbChapter))

			// when
			controller.PatchChapter(w, r)

			// then
			assert.Equal(t, http.StatusUnauthorized, w.Code)
		})

		t.Run("should update one chapter", func(t *testing.T) {
			// given
			w := httptest.NewRecorder()
			r := httptest.NewRequest("PATCH", "/api/v1/chapters/1",
				strings.NewReader(`{"content": "a fine chapter"}`))
			dbBook := &booksModel.Book{
				ID:          1,
				Name:        "Book One",
				AuthorID:    1,
				Description: "! good book",
			}
			r = r.WithContext(context.WithValue(r.Context(), books.MiddleWareBook, dbBook))
			dbChapter := &model.Chapter{
				ID:      1,
				BookID:  1,
				Name:    "Chapter One",
				Price:   100,
				Content: "Nice chapter",
			}
			r = r.WithContext(context.WithValue(r.Context(), authMiddleware.AuthenticatedUserId, uint64(1)))
			r = r.WithContext(context.WithValue(r.Context(), middleWareChapter, dbChapter))

			newDescription := "a fine chapter"
			chapterRepository.
				EXPECT().
				Update(uint64(1), &model.ChapterPatch{Content: &newDescription}).
				Return(nil)

			// when
			controller.PatchChapter(w, r)

			// then
			assert.Equal(t, http.StatusOK, w.Code)
		})
	})

	t.Run("DeleteChapter", func(t *testing.T) {
		t.Run("should return 500 INTERNAL SERVER ERROR if query fails", func(t *testing.T) {
			// given
			w := httptest.NewRecorder()
			r := httptest.NewRequest("DELETE", "/api/v1/chapters/1", nil)
			dbBook := &booksModel.Book{
				ID:          1,
				Name:        "Book One",
				AuthorID:    1,
				Description: "! good book",
			}
			r = r.WithContext(context.WithValue(r.Context(), books.MiddleWareBook, dbBook))
			dbChapter := &model.Chapter{
				ID:      1,
				BookID:  1,
				Name:    "Chapter One",
				Price:   100,
				Content: "Nice chapter",
			}
			r = r.WithContext(context.WithValue(r.Context(), authMiddleware.AuthenticatedUserId, uint64(1)))
			r = r.WithContext(context.WithValue(r.Context(), middleWareChapter, dbChapter))

			chapterRepository.
				EXPECT().
				Delete([]*model.Chapter{dbChapter}).
				Return(errors.New("database error"))

			// when
			controller.DeleteChapter(w, r)

			// then
			assert.Equal(t, http.StatusInternalServerError, w.Code)
		})

		t.Run("should return 401 if not the user who created the chapter", func(t *testing.T) {
			// given
			w := httptest.NewRecorder()
			r := httptest.NewRequest("DELETE", "/api/v1/chapters/1", nil)
			dbBook := &booksModel.Book{
				ID:          1,
				Name:        "Book One",
				AuthorID:    1,
				Description: "! good book",
			}
			r = r.WithContext(context.WithValue(r.Context(), books.MiddleWareBook, dbBook))
			dbChapter := &model.Chapter{
				ID:      1,
				BookID:  1,
				Name:    "Chapter One",
				Price:   100,
				Content: "Nice chapter",
			}
			r = r.WithContext(context.WithValue(r.Context(), authMiddleware.AuthenticatedUserId, uint64(2)))
			r = r.WithContext(context.WithValue(r.Context(), middleWareChapter, dbChapter))

			// when
			controller.DeleteChapter(w, r)

			// then
			assert.Equal(t, http.StatusUnauthorized, w.Code)
		})

		t.Run("should return 200 OK", func(t *testing.T) {
			// given
			w := httptest.NewRecorder()
			r := httptest.NewRequest("DELETE", "/api/v1/chapters/1", nil)
			dbBook := &booksModel.Book{
				ID:          1,
				Name:        "Book One",
				AuthorID:    1,
				Description: "! good book",
			}
			r = r.WithContext(context.WithValue(r.Context(), books.MiddleWareBook, dbBook))
			dbChapter := &model.Chapter{
				ID:      1,
				BookID:  1,
				Name:    "Chapter One",
				Price:   100,
				Content: "Nice chapter",
			}
			r = r.WithContext(context.WithValue(r.Context(), authMiddleware.AuthenticatedUserId, uint64(1)))
			r = r.WithContext(context.WithValue(r.Context(), middleWareChapter, dbChapter))

			chapterRepository.
				EXPECT().
				Delete([]*model.Chapter{dbChapter}).
				Return(nil)

			// when
			controller.DeleteChapter(w, r)

			// then
			assert.Equal(t, http.StatusOK, w.Code)
		})
	})

	t.Run("LoadChapterMiddleware", func(t *testing.T) {
		t.Run("Should return 400 if the chapterid cannot be parsed", func(t *testing.T) {
			// given
			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", "/api/v1/chapters/aaa", nil)
			r = r.WithContext(context.WithValue(r.Context(), "chapterid", "aaa"))
			dbBook := &booksModel.Book{}

			r = r.WithContext(context.WithValue(r.Context(), books.MiddleWareBook, dbBook))

			// when
			called := false
			controller.LoadChapterMiddleware(w, r, func(r *http.Request) {
				called = true
			})

			assert.Equal(t, false, called)
			assert.Equal(t, http.StatusBadRequest, w.Code)
		})

		t.Run("Should return 404 if the query fails", func(t *testing.T) {
			// given
			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", "/api/v1/chapters/1", nil)
			dbBook := &booksModel.Book{
				ID:          1,
				Name:        "Book One",
				AuthorID:    1,
				Description: "! good book",
			}

			r = r.WithContext(context.WithValue(r.Context(), books.MiddleWareBook, dbBook))
			r = r.WithContext(context.WithValue(r.Context(), "chapterid", "1"))
			dbChapter := &model.Chapter{
				ID:      1,
				BookID:  1,
				Name:    "Chapter One",
				Price:   100,
				Content: "Nice chapter",
			}

			chapterRepository.
				EXPECT().
				FindByIdAndBookId(dbChapter.BookID, dbBook.ID).
				Return(nil, errors.New("database error"))

			// when
			called := false
			controller.LoadChapterMiddleware(w, r, func(r *http.Request) {
				called = true
			})

			assert.Equal(t, false, called)
			assert.Equal(t, http.StatusNotFound, w.Code)
		})

		t.Run("Should return 200 if it succeeds", func(t *testing.T) {
			// given
			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", "/api/v1/chapters/1", nil)
			dbBook := &booksModel.Book{
				ID:          1,
				Name:        "Book One",
				AuthorID:    1,
				Description: "! good book",
			}

			r = r.WithContext(context.WithValue(r.Context(), books.MiddleWareBook, dbBook))
			r = r.WithContext(context.WithValue(r.Context(), "chapterid", "1"))
			dbChapter := &model.Chapter{
				ID:      1,
				BookID:  1,
				Name:    "Chapter One",
				Price:   100,
				Content: "Nice chapter",
			}

			chapterRepository.
				EXPECT().
				FindByIdAndBookId(dbChapter.BookID, dbBook.ID).
				Return(dbChapter, nil)

			// when
			called := false
			controller.LoadChapterMiddleware(w, r, func(req *http.Request) {
				called = true
				r = req
			})

			assert.Equal(t, true, called)
			assert.Equal(t, dbChapter, r.Context().Value(middleWareChapter))
			assert.Equal(t, http.StatusOK, w.Code)
		})
	})
}
