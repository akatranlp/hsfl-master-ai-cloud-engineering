package books

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/akatranlp/hsfl-master-ai-cloud-engineering/book-service/books/model"
	auth_middleware "github.com/akatranlp/hsfl-master-ai-cloud-engineering/lib/auth-middleware"
	"github.com/akatranlp/hsfl-master-ai-cloud-engineering/lib/router"
	"golang.org/x/sync/singleflight"
)

type bookContext string

const (
	MiddleWareBook bookContext = "book"
)

type DefaultController struct {
	bookRepository Repository
	g              *singleflight.Group
}

func NewDefaultController(
	bookRepository Repository,
) *DefaultController {
	g := &singleflight.Group{}
	return &DefaultController{bookRepository, g}
}

func (ctrl *DefaultController) GetBooks(w http.ResponseWriter, r *http.Request) {
	userId := r.URL.Query().Get("userId")

	var books []*model.Book
	if userId != "" {
		id, err := strconv.ParseUint(userId, 10, 64)
		if err != nil {
			http.Error(w, "could not parse userid", http.StatusBadRequest)
			return
		}
		newBooks, err, _ := ctrl.g.Do(fmt.Sprintf("books-%d", id), func() (interface{}, error) {
			return ctrl.bookRepository.FindAllByUserId(id)
		})
		if err != nil {
			http.Error(w, "Error while getting the books", http.StatusBadRequest)
			return
		}
		books = newBooks.([]*model.Book)
	} else {
		newBooks, err, _ := ctrl.g.Do("books", func() (interface{}, error) {
			return ctrl.bookRepository.FindAll()
		})
		if err != nil {
			http.Error(w, "Error while getting the books", http.StatusInternalServerError)
			return
		}
		books = newBooks.([]*model.Book)
	}

	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(books)
}

type createBookRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (r createBookRequest) isValid() bool {
	return r.Name != "" && r.Description != ""
}

func (ctrl *DefaultController) PostBook(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(auth_middleware.AuthenticatedUserId).(uint64)

	var request createBookRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if !request.isValid() {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := ctrl.bookRepository.Create([]*model.Book{{
		Name:        request.Name,
		AuthorID:    userId,
		Description: request.Description,
	}}); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (ctrl *DefaultController) GetBook(w http.ResponseWriter, r *http.Request) {
	book := r.Context().Value(MiddleWareBook).(*model.Book)

	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(book)
}

type updateBookRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (ctrl *DefaultController) PatchBook(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(auth_middleware.AuthenticatedUserId).(uint64)
	book := r.Context().Value(MiddleWareBook).(*model.Book)

	if userId != book.AuthorID {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	var request updateBookRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var patchBook model.BookPatch
	if request.Name != "" {
		patchBook.Name = &request.Name
	}
	if request.Description != "" {
		patchBook.Description = &request.Description
	}

	if err := ctrl.bookRepository.Update(book.ID, &patchBook); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (ctrl *DefaultController) DeleteBook(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(auth_middleware.AuthenticatedUserId).(uint64)
	book := r.Context().Value(MiddleWareBook).(*model.Book)

	if userId != book.AuthorID {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if err := ctrl.bookRepository.Delete([]*model.Book{book}); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (ctrl *DefaultController) LoadBookMiddleware(w http.ResponseWriter, r *http.Request, next router.Next) {
	bookId := r.Context().Value("bookid").(string)

	id, err := strconv.ParseUint(bookId, 10, 64)
	if err != nil {
		http.Error(w, "can't parse the bookId", http.StatusBadRequest)
		return
	}

	newBook, err, _ := ctrl.g.Do(fmt.Sprintf("book-%d", id), func() (interface{}, error) {
		return ctrl.bookRepository.FindById(id)
	})
	if err != nil {
		http.Error(w, "can't find the book", http.StatusNotFound)
		return
	}
	book := newBook.(*model.Book)

	ctx := context.WithValue(r.Context(), MiddleWareBook, book)
	next(r.WithContext(ctx))
}
