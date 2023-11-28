package user

import (
	"net/http"

	"github.com/akatranlp/hsfl-master-ai-cloud-engineering/lib/router"
)

type Controller interface {
	Login(http.ResponseWriter, *http.Request)
	Register(http.ResponseWriter, *http.Request)
	RefreshToken(http.ResponseWriter, *http.Request)
	Logout(http.ResponseWriter, *http.Request)
	ValidateToken(http.ResponseWriter, *http.Request)
	MoveUserAmount(http.ResponseWriter, *http.Request)
	GetUsers(http.ResponseWriter, *http.Request)
	GetMe(http.ResponseWriter, *http.Request)
	PatchMe(http.ResponseWriter, *http.Request)
	DeleteMe(http.ResponseWriter, *http.Request)
	GetUser(http.ResponseWriter, *http.Request)
	AuthenticationMiddleWare(http.ResponseWriter, *http.Request, router.Next)
}
