package grpc

import (
	"context"
	"log"

	"github.com/akatranlp/hsfl-master-ai-cloud-engineering/book-service/books"
	"github.com/akatranlp/hsfl-master-ai-cloud-engineering/book-service/chapters"
	"github.com/akatranlp/hsfl-master-ai-cloud-engineering/book-service/chapters/model"
	"github.com/akatranlp/hsfl-master-ai-cloud-engineering/lib/grpc/book-service/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type server struct {
	proto.UnimplementedBookServiceServer
	bookRepository    books.Repository
	chapterRepository chapters.Repository
}

func NewServer(bookRepository books.Repository, chapterRepository chapters.Repository) proto.BookServiceServer {
	return &server{
		bookRepository:    bookRepository,
		chapterRepository: chapterRepository,
	}
}

func (s *server) ValidateChapterId(ctx context.Context, req *proto.ValidateChapterIdRequest) (*proto.ValidateChapterIdResponse, error) {

	chapter, receivingUserId, err := s.chapterRepository.ValidateChapterId(req.ChapterId, req.BookId)
	log.Println(req.ChapterId, req.BookId)
	if err != nil {
		log.Println("ERROR [ValidateChapterId - Execute ValidateChapterId]: ", err.Error())
		return nil, status.Error(codes.NotFound, err.Error())
	}

	if *receivingUserId == req.UserId {
		log.Println("ERROR [ValidateChapterId - receivingUserId == request.UserId]: ", "Author and buyer are the same")
		return nil, status.Error(codes.InvalidArgument, "Author and buyer are the same")
	}

	if chapter.Status != model.Published {
		log.Println("ERROR [ValidateChapterId - chapter.Status != model.Published]: ", "Chapter is not published")
		return nil, status.Error(codes.InvalidArgument, "Chapter is not published")
	}

	response := &proto.ValidateChapterIdResponse{
		ChapterId:       chapter.ID,
		BookId:          chapter.BookID,
		ReceivingUserId: *receivingUserId,
		Amount:          chapter.Price,
	}
	return response, nil
}
