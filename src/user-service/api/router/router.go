package router

import (
	"github.com/akatranlp/hsfl-master-ai-cloud-engineering/lib/router"
	"github.com/akatranlp/hsfl-master-ai-cloud-engineering/user-service/user"
	"net/http"
)

type Router struct {
	router http.Handler
}

func New(
	userController user.Controller,
) *Router {
	r := router.New()

	r.POST("/api/v1/login", userController.Login)
	r.POST("/api/v1/register", userController.Register)

	r.USE("/api/v1/users", userController.AuthenticationMiddleWare)

	r.GET("/api/v1/users", userController.GetUsers)
	r.GET("/api/v1/users/me", userController.GetMe)
	r.PATCH("/api/v1/users/me", userController.PatchMe)
	r.DELETE("/api/v1/users/me", userController.DeleteMe)
	r.GET("/api/v1/users/:userid", userController.GetUser)

	// only accessible intern
	r.POST("/validate-token", userController.ValidateToken)
	r.POST("/move-user-amount", userController.MoveUserAmount)

	return &Router{r}
}

func (router *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	router.router.ServeHTTP(w, r)
}
