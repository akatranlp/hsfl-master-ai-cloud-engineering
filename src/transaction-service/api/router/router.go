package router

import (
	auth_middleware "github.com/akatranlp/hsfl-master-ai-cloud-engineering/lib/auth-middleware"
	"github.com/akatranlp/hsfl-master-ai-cloud-engineering/lib/router"
	"github.com/akatranlp/hsfl-master-ai-cloud-engineering/transaction-service/transactions"
	"net/http"
)

type Router struct {
	router http.Handler
}

func New(
	transactionController transactions.Controller,
	authController auth_middleware.Controller,
) *Router {
	transactionsRouter := router.New()

	transactionsRouter.USE("/api/v1/transactions", authController.AuthenticationMiddleware)
	transactionsRouter.GET("/api/v1/transactions", transactionController.GetYourTransactions)
	transactionsRouter.POST("/api/v1/transactions", transactionController.CreateTransaction)
	return &Router{transactionsRouter}
}

func (router *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	router.router.ServeHTTP(w, r)
}
