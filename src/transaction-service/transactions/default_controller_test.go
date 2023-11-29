package transactions

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	auth_middleware "github.com/akatranlp/hsfl-master-ai-cloud-engineering/lib/auth-middleware"
	shared_types "github.com/akatranlp/hsfl-master-ai-cloud-engineering/lib/shared-types"
	book_service_client_mocks "github.com/akatranlp/hsfl-master-ai-cloud-engineering/transaction-service/_mocks/book-service-client"
	transaction_mocks "github.com/akatranlp/hsfl-master-ai-cloud-engineering/transaction-service/_mocks/transactions"
	user_service_client_mocks "github.com/akatranlp/hsfl-master-ai-cloud-engineering/transaction-service/_mocks/user-service-client"
	"github.com/akatranlp/hsfl-master-ai-cloud-engineering/transaction-service/transactions/model"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestTransactionDefaultController(t *testing.T) {
	ctrl := gomock.NewController(t)
	transactionRepository := transaction_mocks.NewMockRepository(ctrl)
	bookClientRepository := book_service_client_mocks.NewMockRepository(ctrl)
	userClientRepository := user_service_client_mocks.NewMockRepository(ctrl)

	controller := NewDefaultController(transactionRepository, bookClientRepository, userClientRepository)

	t.Run("GetYourTransactions", func(t *testing.T) {
		t.Run("should return 500 INTERNAL SERVER ERROR if query failed", func(t *testing.T) {
			// given
			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", "/api/v1/transactions", nil)
			r = r.WithContext(context.WithValue(r.Context(), auth_middleware.AuthenticatedUserId, uint64(1)))

			transactionRepository.
				EXPECT().
				FindAllForUserId(uint64(1)).
				Return(nil, errors.New("query failed")).
				Times(1)

			// when
			controller.GetYourTransactions(w, r)

			// then
			assert.Equal(t, http.StatusInternalServerError, w.Code)
		})

		t.Run("should return all transactions", func(t *testing.T) {
			// given
			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", "/api/v1/transactions", nil)
			r = r.WithContext(context.WithValue(r.Context(), auth_middleware.AuthenticatedUserId, uint64(1)))

			transactionRepository.
				EXPECT().
				FindAllForUserId(uint64(1)).
				Return([]*model.Transaction{{ID: 999}}, nil).
				Times(1)

			// when
			controller.GetYourTransactions(w, r)

			// then
			res := w.Result()
			var response []model.Transaction
			err := json.NewDecoder(res.Body).Decode(&response)

			assert.NoError(t, err)
			assert.Equal(t, http.StatusOK, w.Code)
			assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
			assert.Len(t, response, 1)
			assert.Equal(t, uint64(999), response[0].ID)
		})

		t.Run("should return 500 INTERNAL SERVER ERROR if query failed", func(t *testing.T) {
			// given
			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", "/api/v1/transactions?receiving=True", nil)
			r = r.WithContext(context.WithValue(r.Context(), auth_middleware.AuthenticatedUserId, uint64(1)))

			transactionRepository.
				EXPECT().
				FindAllForReceivingUserId(uint64(1)).
				Return(nil, errors.New("query failed")).
				Times(1)

			// when
			controller.GetYourTransactions(w, r)

			// then
			assert.Equal(t, http.StatusInternalServerError, w.Code)
		})

		t.Run("should return all transactions", func(t *testing.T) {
			// given
			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", "/api/v1/transactions?receiving=True", nil)
			r = r.WithContext(context.WithValue(r.Context(), auth_middleware.AuthenticatedUserId, uint64(1)))

			transactionRepository.
				EXPECT().
				FindAllForReceivingUserId(uint64(1)).
				Return([]*model.Transaction{{ID: 999}}, nil).
				Times(1)

			// when
			controller.GetYourTransactions(w, r)

			// then
			res := w.Result()
			var response []model.Transaction
			err := json.NewDecoder(res.Body).Decode(&response)

			assert.NoError(t, err)
			assert.Equal(t, http.StatusOK, w.Code)
			assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
			assert.Len(t, response, 1)
			assert.Equal(t, uint64(999), response[0].ID)
		})
	})

	t.Run("CreateTransaction", func(t *testing.T) {
		t.Run("should return 400 BAD REQUEST if payload is not json", func(t *testing.T) {
			tests := []io.Reader{
				nil,
				strings.NewReader(`{"invalid`),
			}

			for _, test := range tests {
				// given
				w := httptest.NewRecorder()
				r := httptest.NewRequest("POST", "/api/v1/transactions", test)
				r = r.WithContext(context.WithValue(r.Context(), auth_middleware.AuthenticatedUserId, uint64(1)))

				// when
				controller.CreateTransaction(w, r)

				// then
				assert.Equal(t, http.StatusBadRequest, w.Code)
			}
		})

		t.Run("should return 400 BAD REQUEST if payload is incomplete", func(t *testing.T) {
			tests := []io.Reader{
				strings.NewReader(`{}`),
			}

			for _, test := range tests {
				// given
				w := httptest.NewRecorder()
				r := httptest.NewRequest("POST", "/api/v1/transactions", test)
				r = r.WithContext(context.WithValue(r.Context(), auth_middleware.AuthenticatedUserId, uint64(1)))

				// when
				controller.CreateTransaction(w, r)

				// then
				assert.Equal(t, http.StatusBadRequest, w.Code)
			}
		})

		t.Run("should return 400 BAD REQUEST if transaction already exist", func(t *testing.T) {
			// given
			w := httptest.NewRecorder()
			r := httptest.NewRequest("POST", "/api/v1/transactions",
				strings.NewReader(`{"chapterID":1, "bookID":1}`))
			id := uint64(1)
			chapterId := uint64(1)
			bookId := uint64(1)
			r = r.WithContext(context.WithValue(r.Context(), auth_middleware.AuthenticatedUserId, id))

			transactionRepository.
				EXPECT().
				FindForUserIdAndChapterId(id, chapterId, bookId).
				Return(&model.Transaction{ID: 1}, nil)

			// when
			controller.CreateTransaction(w, r)

			// then
			assert.Equal(t, http.StatusBadRequest, w.Code)
		})

		t.Run("should return 400 INTERNAL SERVER ERROR if validate ChapterId failed", func(t *testing.T) {
			// given
			w := httptest.NewRecorder()
			r := httptest.NewRequest("POST", "/api/v1/transactions",
				strings.NewReader(`{"chapterID":1, "bookID":1}`))
			id := uint64(1)
			chapterId := uint64(1)
			bookId := uint64(1)
			r = r.WithContext(context.WithValue(r.Context(), auth_middleware.AuthenticatedUserId, id))

			transactionRepository.
				EXPECT().
				FindForUserIdAndChapterId(id, chapterId, bookId).
				Return(nil, errors.New("transaction doesn't exist"))

			bookClientRepository.
				EXPECT().
				ValidateChapterId(id, chapterId, bookId).
				Return(nil, errors.New("client error"))

			// when
			controller.CreateTransaction(w, r)

			// then
			assert.Equal(t, http.StatusBadRequest, w.Code)
		})

		t.Run("should return 500 INTERNAL SERVER ERROR if persisting failed", func(t *testing.T) {
			// given
			w := httptest.NewRecorder()
			r := httptest.NewRequest("POST", "/api/v1/transactions",
				strings.NewReader(`{"chapterID":1, "bookID":1}`))
			id := uint64(1)
			chapterId := uint64(1)
			bookId := uint64(1)
			r = r.WithContext(context.WithValue(r.Context(), auth_middleware.AuthenticatedUserId, id))

			transactionRepository.
				EXPECT().
				FindForUserIdAndChapterId(id, chapterId, bookId).
				Return(nil, errors.New("transaction doesn't exist"))

			bookClientRepository.
				EXPECT().
				ValidateChapterId(id, chapterId, bookId).
				Return(&shared_types.ValidateChapterIdResponse{ChapterId: 1, ReceivingUserId: 2, Amount: 100, BookId: 1}, nil)

			transactionRepository.
				EXPECT().
				Create([]*model.Transaction{{ChapterID: 1, Amount: 100, ReceivingUserID: 2, PayingUserID: 1, BookID: 1}}).
				Return(errors.New("database error"))

			// when
			controller.CreateTransaction(w, r)

			// then
			assert.Equal(t, http.StatusInternalServerError, w.Code)
		})

		t.Run("should return 500 INTERNAL SERVER ERROR if persisting failed", func(t *testing.T) {
			// given
			w := httptest.NewRecorder()
			r := httptest.NewRequest("POST", "/api/v1/transactions",
				strings.NewReader(`{"chapterID":1, "bookID":1}`))
			id := uint64(1)
			chapterId := uint64(1)
			bookId := uint64(1)
			r = r.WithContext(context.WithValue(r.Context(), auth_middleware.AuthenticatedUserId, id))

			transactionRepository.
				EXPECT().
				FindForUserIdAndChapterId(id, chapterId, bookId).
				Return(nil, errors.New("transaction doesn't exist"))

			bookClientRepository.
				EXPECT().
				ValidateChapterId(id, chapterId, bookId).
				Return(&shared_types.ValidateChapterIdResponse{ChapterId: 1, ReceivingUserId: 2, Amount: 100, BookId: 1}, nil)

			transactionRepository.
				EXPECT().
				Create([]*model.Transaction{{ChapterID: 1, Amount: 100, ReceivingUserID: 2, PayingUserID: 1, BookID: 1}}).
				Return(nil)

			userClientRepository.
				EXPECT().
				MoveBalance(id, uint64(2), int64(100)).
				Return(errors.New("failed"))

			// when
			controller.CreateTransaction(w, r)

			// then
			assert.Equal(t, http.StatusInternalServerError, w.Code)
		})

		t.Run("should create new transaction", func(t *testing.T) {
			// given
			w := httptest.NewRecorder()
			r := httptest.NewRequest("POST", "/api/v1/transactions",
				strings.NewReader(`{"chapterID":1, "bookID":1}`))
			id := uint64(1)
			chapterId := uint64(1)
			bookId := uint64(1)
			r = r.WithContext(context.WithValue(r.Context(), auth_middleware.AuthenticatedUserId, id))

			transactionRepository.
				EXPECT().
				FindForUserIdAndChapterId(id, chapterId, bookId).
				Return(nil, errors.New("transaction doesn't exist"))

			bookClientRepository.
				EXPECT().
				ValidateChapterId(id, chapterId, bookId).
				Return(&shared_types.ValidateChapterIdResponse{ChapterId: 1, ReceivingUserId: 2, Amount: 100, BookId: 1}, nil)

			transactionRepository.
				EXPECT().
				Create([]*model.Transaction{{ChapterID: 1, Amount: 100, ReceivingUserID: 2, PayingUserID: 1, BookID: 1}}).
				Return(nil)

			userClientRepository.
				EXPECT().
				MoveBalance(id, uint64(2), int64(100)).
				Return(nil)

			// when
			controller.CreateTransaction(w, r)

			// then
			assert.Equal(t, http.StatusOK, w.Code)
		})
	})

	t.Run("CheckChapterBought", func(t *testing.T) {
		t.Run("should return 400 BAD REQUEST if payload is not json", func(t *testing.T) {
			tests := []io.Reader{
				nil,
				strings.NewReader(`{"invalid`),
			}

			for _, test := range tests {
				// given
				w := httptest.NewRecorder()
				r := httptest.NewRequest("POST", "/check-chapter-bought", test)
				r = r.WithContext(context.WithValue(r.Context(), auth_middleware.AuthenticatedUserId, uint64(1)))

				// when
				controller.CheckChapterBought(w, r)

				// then
				assert.Equal(t, http.StatusBadRequest, w.Code)
			}
		})

		t.Run("should return 400 BAD REQUEST if payload is incomplete", func(t *testing.T) {
			tests := []io.Reader{
				strings.NewReader(`{}`),
			}

			for _, test := range tests {
				// given
				w := httptest.NewRecorder()
				r := httptest.NewRequest("POST", "/check-chapter-bought", test)
				r = r.WithContext(context.WithValue(r.Context(), auth_middleware.AuthenticatedUserId, uint64(1)))

				// when
				controller.CheckChapterBought(w, r)

				// then
				assert.Equal(t, http.StatusBadRequest, w.Code)
			}
		})

		t.Run("should return 404 if transaction was not found", func(t *testing.T) {
			// given
			w := httptest.NewRecorder()
			r := httptest.NewRequest("POST", "/check-chapter-bought", strings.NewReader(`{"userId":1,"chapterId":1, "bookId":1}`))

			transactionRepository.
				EXPECT().
				FindForUserIdAndChapterId(uint64(1), uint64(1), uint64(1)).
				Return(nil, errors.New("transaction doesn't exist"))

			// when
			controller.CheckChapterBought(w, r)

			// then
			assert.Equal(t, http.StatusNotFound, w.Code)
		})

		t.Run("should return 200 if transaction was found", func(t *testing.T) {
			// given
			w := httptest.NewRecorder()
			r := httptest.NewRequest("POST", "/check-chapter-bought", strings.NewReader(`{"userId":1,"chapterId":1,"bookId":1}`))

			transactionRepository.
				EXPECT().
				FindForUserIdAndChapterId(uint64(1), uint64(1), uint64(1)).
				Return(&model.Transaction{}, nil)

			// when
			controller.CheckChapterBought(w, r)

			// then
			assert.Equal(t, http.StatusOK, w.Code)
		})
	})
}
