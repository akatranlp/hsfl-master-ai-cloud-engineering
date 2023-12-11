package book_service_client

import (
	"context"
	"errors"

	"github.com/akatranlp/hsfl-master-ai-cloud-engineering/lib/grpc/book-service/proto"
	shared_types "github.com/akatranlp/hsfl-master-ai-cloud-engineering/lib/shared-types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GRPCRepository struct {
	client proto.BookServiceClient
}

func NewGRPCRepository(client proto.BookServiceClient) *GRPCRepository {
	return &GRPCRepository{
		client: client,
	}
}

func (repo *GRPCRepository) ValidateChapterId(userId uint64, chapterId uint64, bookId uint64) (*shared_types.ValidateChapterIdResponse, error) {
	req := &proto.ValidateChapterIdRequest{
		UserId:    userId,
		ChapterId: chapterId,
		BookId:    bookId,
	}

	res, err := repo.client.ValidateChapterId(context.Background(), req)
	if err != nil {
		st, ok := status.FromError(err)
		if ok && st.Code() == codes.InvalidArgument {
			return nil, errors.New("you cannot buy this book")
		}
		return nil, err
	}

	return &shared_types.ValidateChapterIdResponse{
		ChapterId:       res.ChapterId,
		BookId:          res.BookId,
		ReceivingUserId: res.ReceivingUserId,
		Amount:          res.Amount,
	}, nil
}
