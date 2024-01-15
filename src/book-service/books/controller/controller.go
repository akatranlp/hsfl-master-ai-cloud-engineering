package books_controller

import (
	"net/http"

	"github.com/akatranlp/hsfl-master-ai-cloud-engineering/lib/router"
)

type Controller interface {
	GetBooks(http.ResponseWriter, *http.Request)
	GetBook(http.ResponseWriter, *http.Request)

	PostBook(http.ResponseWriter, *http.Request)
	PatchBook(http.ResponseWriter, *http.Request)
	DeleteBook(http.ResponseWriter, *http.Request)

	LoadBookMiddleware(http.ResponseWriter, *http.Request, router.Next)
}
