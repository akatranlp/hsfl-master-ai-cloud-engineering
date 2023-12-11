package transaction_service_client

import (
	"context"
	"errors"

	"github.com/akatranlp/hsfl-master-ai-cloud-engineering/lib/grpc/transaction-service/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GRPCRepository struct {
	client proto.TransactionServiceClient
}

func NewGRPCRepository(client proto.TransactionServiceClient) *GRPCRepository {
	return &GRPCRepository{
		client: client,
	}
}

func (repo *GRPCRepository) CheckChapterBought(userId uint64, chapterId uint64, bookId uint64) error {
	req := &proto.CheckChapterBoughtRequest{
		UserId:    userId,
		ChapterId: chapterId,
		BookId:    bookId,
	}

	res, err := repo.client.CheckChapterBought(context.Background(), req)
	if err != nil {
		st, ok := status.FromError(err)
		if ok && st.Code() == codes.NotFound {
			return errors.New("you haven't bought this book")
		}
		return err
	}

	if !res.Success {
		return errors.New("an unknown error")
	}

	return nil
}
