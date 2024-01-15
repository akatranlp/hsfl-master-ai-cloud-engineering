package router

import (
	"net/http"

	auth_middleware "github.com/akatranlp/hsfl-master-ai-cloud-engineering/lib/auth-middleware"
	"github.com/akatranlp/hsfl-master-ai-cloud-engineering/lib/health"

	books_controller "github.com/akatranlp/hsfl-master-ai-cloud-engineering/book-service/books/controller"
	chaptersController "github.com/akatranlp/hsfl-master-ai-cloud-engineering/book-service/chapters/controller"
	"github.com/akatranlp/hsfl-master-ai-cloud-engineering/lib/router"
)

type Router struct {
	router http.Handler
}

func New(
	authController auth_middleware.Controller,
	booksController books_controller.Controller,
	chapterController chaptersController.Controller,
	healthController health.Controller,
) *Router {
	booksRouter := router.New()
	booksRouter.GET("/health", healthController.ProvideHealth)
	booksRouter.POST("/valdiate-chapter-id", chapterController.ValidateChapterId)

	booksRouter.USE("/api/v1/books", authController.AuthenticationMiddleware)
	booksRouter.GET("/api/v1/books", booksController.GetBooks)
	booksRouter.POST("/api/v1/books", booksController.PostBook)

	booksRouter.USE("/api/v1/books/:bookid", booksController.LoadBookMiddleware)
	booksRouter.GET("/api/v1/books/:bookid", booksController.GetBook)
	booksRouter.PATCH("/api/v1/books/:bookid", booksController.PatchBook)
	booksRouter.DELETE("/api/v1/books/:bookid", booksController.DeleteBook)

	booksRouter.GET("/api/v1/books/:bookid/chapters", chapterController.GetChaptersForBook)
	booksRouter.POST("/api/v1/books/:bookid/chapters", chapterController.PostChapter)

	booksRouter.USE("/api/v1/books/:bookid/chapters/:chapterid", chapterController.LoadChapterMiddleware)
	booksRouter.GET("/api/v1/books/:bookid/chapters/:chapterid", chapterController.GetChapterForBook)
	booksRouter.PATCH("/api/v1/books/:bookid/chapters/:chapterid", chapterController.PatchChapter)
	booksRouter.DELETE("/api/v1/books/:bookid/chapters/:chapterid", chapterController.DeleteChapter)

	return &Router{booksRouter}
}

func (router *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	router.router.ServeHTTP(w, r)
}
