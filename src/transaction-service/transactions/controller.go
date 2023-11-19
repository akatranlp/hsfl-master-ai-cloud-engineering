package transactions

import (
	"net/http"
)

type Controller interface {
	GetYourTransactions(http.ResponseWriter, *http.Request)
	CreateTransaction(http.ResponseWriter, *http.Request)
	CheckChapterBought(http.ResponseWriter, *http.Request)
}
