package transactions

import (
	"encoding/json"
	auth_middleware "github.com/akatranlp/hsfl-master-ai-cloud-engineering/lib/auth-middleware"
	book_service_client "github.com/akatranlp/hsfl-master-ai-cloud-engineering/transaction-service/book-service-client"
	"github.com/akatranlp/hsfl-master-ai-cloud-engineering/transaction-service/transactions/model"
	user_service_client "github.com/akatranlp/hsfl-master-ai-cloud-engineering/transaction-service/user-service-client"
	"net/http"
)

type DefaultController struct {
	transactionRepository Repository
	bookClientRepository  book_service_client.Repository
	userClientRepository  user_service_client.Repository
}

func NewDefaultController(
	transactionRepository Repository,
	bookClientRepository book_service_client.Repository,
	userClientRepository user_service_client.Repository,
) *DefaultController {
	return &DefaultController{transactionRepository, bookClientRepository, userClientRepository}
}

func (ctrl *DefaultController) GetYourTransactions(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(auth_middleware.AuthenticatedUserId).(uint64)

	receiving := r.URL.Query().Get("receiving") != ""

	var transactions []*model.Transaction
	var err error
	if receiving {
		transactions, err = ctrl.transactionRepository.FindAllForReceivingUserId(userId)
	} else {
		transactions, err = ctrl.transactionRepository.FindAllForUserId(userId)
	}

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(transactions)
}

type createTransactionRequest struct {
	ChapterID uint64 `json:"chapterID"`
}

func (r createTransactionRequest) isValid() bool {
	return r.ChapterID != 0
}

func (ctrl *DefaultController) CreateTransaction(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(auth_middleware.AuthenticatedUserId).(uint64)

	var request createTransactionRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if !request.isValid() {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	validatedInfo, err := ctrl.bookClientRepository.ValidateChapterId(userId, request.ChapterID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := ctrl.transactionRepository.Create([]*model.Transaction{{
		ChapterID:       request.ChapterID,
		PayingUserID:    userId,
		ReceivingUserID: validatedInfo.ReceivingUserId,
		BookID:          validatedInfo.BookId,
		Amount:          validatedInfo.Amount,
	}}); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := ctrl.userClientRepository.MoveBalance(userId, validatedInfo.ReceivingUserId, int64(validatedInfo.Amount)); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

}
