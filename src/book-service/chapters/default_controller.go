package chapters

import (
	"context"
	"encoding/json"
	"github.com/akatranlp/hsfl-master-ai-cloud-engineering/book-service/books"
	booksModel "github.com/akatranlp/hsfl-master-ai-cloud-engineering/book-service/books/model"
	"github.com/akatranlp/hsfl-master-ai-cloud-engineering/book-service/chapters/model"
	authMiddleware "github.com/akatranlp/hsfl-master-ai-cloud-engineering/lib/auth-middleware"
	"github.com/akatranlp/hsfl-master-ai-cloud-engineering/lib/router"
	"net/http"
	"strconv"
)

type chapterContext string

const (
	middleWareChapter chapterContext = "chapter"
)

type DefaultController struct {
	chapterRepository Repository
}

func NewDefaultController(
	chapterRepository Repository,
) *DefaultController {
	return &DefaultController{chapterRepository}
}
func (ctrl *DefaultController) GetChaptersForBook(w http.ResponseWriter, r *http.Request) {
	book := r.Context().Value(books.MiddleWareBook).(*booksModel.Book)

	chapters, err := ctrl.chapterRepository.FindAllPreviewsByBookId(book.ID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(chapters)
}

type createChapterRequest struct {
	Name    string  `json:"name"`
	Price   *uint64 `json:"price"`
	Content string  `json:"content"`
}

func (r createChapterRequest) isValid() bool {
	return r.Name != "" && r.Price != nil && r.Content != ""
}

func (ctrl *DefaultController) PostChapter(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(authMiddleware.AuthenticatedUserId).(uint64)
	book := r.Context().Value(books.MiddleWareBook).(*booksModel.Book)

	if userId != book.AuthorID {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	var request createChapterRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if !request.isValid() {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := ctrl.chapterRepository.Create([]*model.Chapter{{
		BookID:  book.ID,
		Name:    request.Name,
		Price:   *request.Price,
		Content: request.Content,
	}}); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (ctrl *DefaultController) GetChapter(w http.ResponseWriter, r *http.Request) {
	chapterId := r.Context().Value("chapterid").(string)

	id, err := strconv.ParseUint(chapterId, 10, 64)
	if err != nil {
		http.Error(w, "can't parse the chapterId", http.StatusBadRequest)
		return
	}

	chapter, err := ctrl.chapterRepository.FindById(id)
	if err != nil {
		http.Error(w, "can't find the chapter", http.StatusNotFound)
		return
	}

	preview := model.ChapterPreview{
		ID:     chapter.ID,
		BookID: chapter.BookID,
		Name:   chapter.Name,
		Price:  chapter.Price,
	}

	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(preview)
}

func (ctrl *DefaultController) GetChapterForBook(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(authMiddleware.AuthenticatedUserId).(uint64)
	book := r.Context().Value(books.MiddleWareBook).(*booksModel.Book)
	chapter := r.Context().Value(middleWareChapter).(*model.Chapter)

	if userId == book.AuthorID {
		w.Header().Add("Content-Type", "application/json")
		json.NewEncoder(w).Encode(chapter)
		return
	}

	// Check if the user has bought the chapter

	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(chapter)
}

type updateChapterRequest struct {
	Name    string  `json:"name"`
	Price   *uint64 `json:"price"`
	Content string  `json:"content"`
}

func (ctrl *DefaultController) PatchChapter(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(authMiddleware.AuthenticatedUserId).(uint64)
	book := r.Context().Value(books.MiddleWareBook).(*booksModel.Book)
	chapter := r.Context().Value(middleWareChapter).(*model.Chapter)

	if userId != book.AuthorID {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	var request updateChapterRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var patchChapter model.ChapterPatch
	if request.Name != "" {
		patchChapter.Name = &request.Name
	}
	if request.Content != "" {
		patchChapter.Content = &request.Content
	}
	if request.Price != nil {
		patchChapter.Price = request.Price
	}

	if err := ctrl.chapterRepository.Update(chapter.ID, &patchChapter); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (ctrl *DefaultController) DeleteChapter(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(authMiddleware.AuthenticatedUserId).(uint64)
	book := r.Context().Value(books.MiddleWareBook).(*booksModel.Book)
	chapter := r.Context().Value(middleWareChapter).(*model.Chapter)

	if userId != book.AuthorID {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if err := ctrl.chapterRepository.Delete([]*model.Chapter{chapter}); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (ctrl *DefaultController) LoadChapterMiddleware(w http.ResponseWriter, r *http.Request, next router.Next) {
	book := r.Context().Value(books.MiddleWareBook).(*booksModel.Book)
	chapterId := r.Context().Value("chapterid").(string)

	id, err := strconv.ParseUint(chapterId, 10, 64)
	if err != nil {
		http.Error(w, "can't parse the chapterId", http.StatusBadRequest)
		return
	}

	chapter, err := ctrl.chapterRepository.FindByIdAndBookId(id, book.ID)
	if err != nil {
		http.Error(w, "can't find the chapter", http.StatusNotFound)
		return
	}

	ctx := context.WithValue(r.Context(), middleWareChapter, chapter)
	next(r.WithContext(ctx))
}
