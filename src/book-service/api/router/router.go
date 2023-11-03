package router

import (
	"net/http"

	"github.com/akatranlp/hsfl-master-ai-cloud-engineering/book-service/books"
	"github.com/akatranlp/hsfl-master-ai-cloud-engineering/lib/router"
)

type Router struct {
	router http.Handler
}

func New(
	booksController books.Controller,
) *Router {
	booksRouter := router.New()

	booksRouter.GET("/api/v1/books", booksController.GetBooks)
	booksRouter.POST("/api/v1/books", booksController.PostBooks)
	booksRouter.GET("/api/v1/books/:bookid", booksController.GetBook)
	booksRouter.PUT("/api/v1/books/:bookid", booksController.PutBook)
	booksRouter.DELETE("/api/v1/books/:bookid", booksController.DeleteBook)
	booksRouter.GET("/api/v1/books/:bookid/chapters", booksController.GetChapters)
	booksRouter.POST("/api/v1/books/:bookid/chapters", booksController.PostChapters)
	booksRouter.GET("/api/v1/books/:bookid/chapters/:chapterid", booksController.GetChapter)
	booksRouter.PUT("/api/v1/books/:bookid/chapters/:chapterid", booksController.PutChapter)
	booksRouter.DELETE("/api/v1/books/:bookid/chapters/:chapterid", booksController.DeleteChapter)
	return &Router{booksRouter}
}

func (router *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	router.router.ServeHTTP(w, r)
}
