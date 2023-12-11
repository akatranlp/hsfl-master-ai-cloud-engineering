package grpc

import (
	"context"
	"log"

	"github.com/akatranlp/hsfl-master-ai-cloud-engineering/lib/grpc/transaction-service/proto"
	"github.com/akatranlp/hsfl-master-ai-cloud-engineering/transaction-service/transactions"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type server struct {
	proto.UnimplementedTransactionServiceServer
	transactionRepository transactions.Repository
}

func NewServer(
	transactionRepository transactions.Repository,
) proto.TransactionServiceServer {
	return &server{
		transactionRepository: transactionRepository,
	}
}

func (s *server) CheckChapterBought(ctx context.Context, req *proto.CheckChapterBoughtRequest) (*proto.CheckChapterBoughtResponse, error) {
	transaction, err := s.transactionRepository.FindForUserIdAndChapterId(req.UserId, req.ChapterId, req.BookId)
	if err != nil {
		log.Println("ERROR [CheckChapterBought - FindForUserIdAndChapterId]: ", err.Error())
		return nil, status.Error(codes.NotFound, err.Error())
	}

	response := &proto.CheckChapterBoughtResponse{
		Success: transaction != nil,
	}
	return response, nil
}
