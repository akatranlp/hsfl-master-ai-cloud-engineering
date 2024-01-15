package router

import (
	"net/http"

	"github.com/akatranlp/hsfl-master-ai-cloud-engineering/lib/health"
	"github.com/akatranlp/hsfl-master-ai-cloud-engineering/lib/router"

	"github.com/akatranlp/hsfl-master-ai-cloud-engineering/test-data-service/controller"
)

type Router struct {
	router http.Handler
}

func New(
	controller controller.Controller,
	healthController health.Controller,
) *Router {
	router := router.New()

	router.GET("/health", healthController.ProvideHealth)
	router.POST("/api/v1/reset", controller.ResetDatabase)

	return &Router{router}
}

func (router *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	router.router.ServeHTTP(w, r)
}
