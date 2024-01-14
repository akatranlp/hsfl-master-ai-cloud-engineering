package grpc

import (
	"context"
	"log"

	"github.com/akatranlp/hsfl-master-ai-cloud-engineering/lib/grpc/transaction-service/proto"
	"github.com/akatranlp/hsfl-master-ai-cloud-engineering/transaction-service/service"
	"google.golang.org/grpc/status"
)

type server struct {
	proto.UnimplementedTransactionServiceServer
	service service.Service
}

func NewServer(
	service service.Service,
) proto.TransactionServiceServer {
	return &server{service: service}
}

func (s *server) CheckChapterBought(ctx context.Context, req *proto.CheckChapterBoughtRequest) (*proto.CheckChapterBoughtResponse, error) {
	success, statuCode, err := s.service.CheckChapterBought(req.UserId, req.ChapterId, req.BookId)
	if err != nil {
		log.Println("ERROR [CheckChapterBought - FindForUserIdAndChapterId]: ", err.Error())
		return nil, status.Error(statuCode.ToGRPCStatusCode(), err.Error())
	}

	response := &proto.CheckChapterBoughtResponse{
		Success: success,
	}
	return response, nil
}
