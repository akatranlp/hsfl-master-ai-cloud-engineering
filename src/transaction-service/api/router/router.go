package router

import (
	"net/http"

	auth_middleware "github.com/akatranlp/hsfl-master-ai-cloud-engineering/lib/auth-middleware"
	"github.com/akatranlp/hsfl-master-ai-cloud-engineering/lib/health"
	"github.com/akatranlp/hsfl-master-ai-cloud-engineering/lib/router"
	transactions_controller "github.com/akatranlp/hsfl-master-ai-cloud-engineering/transaction-service/transactions/controller"
)

type Router struct {
	router http.Handler
}

func New(
	transactionController transactions_controller.Controller,
	authController auth_middleware.Controller,
	healthController health.Controller,
) *Router {
	transactionsRouter := router.New()

	transactionsRouter.GET("/health", healthController.ProvideHealth)
	transactionsRouter.POST("/check-chapter-bought", transactionController.CheckChapterBought)

	transactionsRouter.USE("/api/v1/transactions", authController.AuthenticationMiddleware)
	transactionsRouter.GET("/api/v1/transactions", transactionController.GetYourTransactions)
	transactionsRouter.POST("/api/v1/transactions", transactionController.CreateTransaction)

	return &Router{transactionsRouter}
}

func (router *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	router.router.ServeHTTP(w, r)
}
