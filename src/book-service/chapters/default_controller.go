package chapters

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/akatranlp/hsfl-master-ai-cloud-engineering/book-service/books"
	booksModel "github.com/akatranlp/hsfl-master-ai-cloud-engineering/book-service/books/model"
	"github.com/akatranlp/hsfl-master-ai-cloud-engineering/book-service/chapters/model"
	transaction_service_client "github.com/akatranlp/hsfl-master-ai-cloud-engineering/book-service/transaction-service-client"
	authMiddleware "github.com/akatranlp/hsfl-master-ai-cloud-engineering/lib/auth-middleware"
	"github.com/akatranlp/hsfl-master-ai-cloud-engineering/lib/router"
	shared_types "github.com/akatranlp/hsfl-master-ai-cloud-engineering/lib/shared-types"
)

type chapterContext string

const (
	middleWareChapter chapterContext = "chapter"
)

type DefaultController struct {
	chapterRepository Repository
	transactionServiceClient transaction_service_client.Repository
}

func NewDefaultController(
	chapterRepository Repository,
	transactionServiceClient transaction_service_client.Repository,
) *DefaultController {
	return &DefaultController{chapterRepository, transactionServiceClient}
}
func (ctrl *DefaultController) GetChaptersForBook(w http.ResponseWriter, r *http.Request) {
	book := r.Context().Value(books.MiddleWareBook).(*booksModel.Book)

	chapters, err := ctrl.chapterRepository.FindAllPreviewsByBookId(book.ID)
	if err != nil {
		log.Println("ERROR [GetChaptersForBook - FindAllPreviewsByBookId]: ", err.Error())
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
		log.Println("ERROR [PostChapter - userId != book.AuthorID]: ", "You are not the owner of the book")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	var request createChapterRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		log.Println("ERROR [PostChapter - Decode createChapterRequest]: ", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if !request.isValid() {
		log.Println("ERROR [PostChapter - Validate createChapterRequest]: ", "Invalid request")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := ctrl.chapterRepository.Create([]*model.Chapter{{
		BookID:  book.ID,
		Name:    request.Name,
		Price:   *request.Price,
		Content: request.Content,
	}}); err != nil {
		log.Println("ERROR [PostChapter - Create]: ", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (ctrl *DefaultController) GetChapter(w http.ResponseWriter, r *http.Request) {
	chapterId := r.Context().Value("chapterid").(string)

	id, err := strconv.ParseUint(chapterId, 10, 64)
	if err != nil {
		log.Println("ERROR [GetChapter - ParseUint]: ", err.Error())
		http.Error(w, "can't parse the chapterId", http.StatusBadRequest)
		return
	}

	chapter, err := ctrl.chapterRepository.FindById(id)
	if err != nil {
		log.Println("ERROR [GetChapter - FindById]: ", err.Error())
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

	response, err :=  ctrl.transactionServiceClient.CheckChapterBought(userId, chapter.ID)
	if err != nil || !response.Success {
		log.Println("ERROR [GetChapterForBook - CheckChapterBought]: ", err.Error())
		w.WriteHeader(http.StatusPaymentRequired)
		return
	}

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
		log.Println("ERROR [PatchChapter - userId != book.AuthorID]: ", "You are not the owner of the book")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	var request updateChapterRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		log.Println("ERROR [PatchChapter - Decode updateChapterRequest]: ", err.Error())
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
		log.Println("ERROR [PatchChapter - Update]: ", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (ctrl *DefaultController) DeleteChapter(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(authMiddleware.AuthenticatedUserId).(uint64)
	book := r.Context().Value(books.MiddleWareBook).(*booksModel.Book)
	chapter := r.Context().Value(middleWareChapter).(*model.Chapter)

	if userId != book.AuthorID {
		log.Println("ERROR [DeleteChapter - userId != book.AuthorID]: ", "You are not the owner of the book")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if err := ctrl.chapterRepository.Delete([]*model.Chapter{chapter}); err != nil {
		log.Println("ERROR [DeleteChapter - Delete]: ", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (ctrl *DefaultController) LoadChapterMiddleware(w http.ResponseWriter, r *http.Request, next router.Next) {
	book := r.Context().Value(books.MiddleWareBook).(*booksModel.Book)
	chapterId := r.Context().Value("chapterid").(string)

	id, err := strconv.ParseUint(chapterId, 10, 64)
	if err != nil {
		log.Println("ERROR [LoadChapterMiddleware - ParseUint]: ", err.Error())
		http.Error(w, "can't parse the chapterId", http.StatusBadRequest)
		return
	}

	chapter, err := ctrl.chapterRepository.FindByIdAndBookId(id, book.ID)
	if err != nil {
		log.Println("ERROR [LoadChapterMiddleware - FindByIdAndBookId]: ", err.Error())
		http.Error(w, "can't find the chapter", http.StatusNotFound)
		return
	}

	ctx := context.WithValue(r.Context(), middleWareChapter, chapter)
	next(r.WithContext(ctx))
}

func (ctrl *DefaultController) ValidateChapterId(w http.ResponseWriter,r *http.Request) {
	var request shared_types.ValidateChapterIdRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		log.Println("ERROR [ValidateChapterId - Decode ValidateChapterIdRequest]: ", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if !request.IsValid() {
		log.Println("ERROR [ValidateChapterId - Validate ValidateChapterIdRequest]: ", "Invalid request")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	chapter, receivingUserId, err := ctrl.chapterRepository.ValidateChapterId(request.ChapterId)
	if err != nil {
		log.Println("ERROR [ValidateChapterId - ValidateChapterId]: ", err.Error())
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if *receivingUserId == request.UserId {
		log.Println("ERROR [ValidateChapterId - receivingUserId == request.UserId]: ", "Author and buyer are the same")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(shared_types.ValidateChapterIdResponse{
		ChapterId: 	 chapter.ID,
		BookId:      chapter.BookID,
		ReceivingUserId: *receivingUserId,
		Amount:      chapter.Price,
	})
}
