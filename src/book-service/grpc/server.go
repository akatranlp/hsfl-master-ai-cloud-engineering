package grpc

import (
	"context"
	"log"

	"github.com/akatranlp/hsfl-master-ai-cloud-engineering/book-service/service"
	"github.com/akatranlp/hsfl-master-ai-cloud-engineering/lib/grpc/book-service/proto"
	"google.golang.org/grpc/status"
)

type server struct {
	proto.UnimplementedBookServiceServer
	service service.Service
}

func NewServer(service service.Service) proto.BookServiceServer {
	return &server{service: service}
}

func (s *server) ValidateChapterId(ctx context.Context, req *proto.ValidateChapterIdRequest) (*proto.ValidateChapterIdResponse, error) {
	result, statusCode, err := s.service.ValidateChapterId(req.UserId, req.ChapterId, req.BookId)
	if err != nil {
		log.Println("ERROR [ValidateChapterId - Execute ValidateChapterId]: ", err.Error())
		return nil, status.Error(statusCode.ToGRPCStatusCode(), err.Error())
	}

	return &proto.ValidateChapterIdResponse{
		ChapterId:       result.ChapterId,
		BookId:          result.BookId,
		ReceivingUserId: result.ReceivingUserId,
		Amount:          result.Amount,
	}, nil
}
