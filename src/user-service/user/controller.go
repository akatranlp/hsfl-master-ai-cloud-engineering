package user

import (
	"github.com/akatranlp/hsfl-master-ai-cloud-engineering/lib/router"
	"net/http"
)

type Controller interface {
	Login(http.ResponseWriter, *http.Request)
	Register(http.ResponseWriter, *http.Request)
	ValidateToken(http.ResponseWriter, *http.Request)
	ChangeUserBalance(http.ResponseWriter, *http.Request)
	GetUsers(http.ResponseWriter, *http.Request)
	GetMe(http.ResponseWriter, *http.Request)
	PatchMe(http.ResponseWriter, *http.Request)
	DeleteMe(http.ResponseWriter, *http.Request)
	GetUser(http.ResponseWriter, *http.Request)
	AuthenticationMiddleWare(http.ResponseWriter, *http.Request, router.Next)
}
