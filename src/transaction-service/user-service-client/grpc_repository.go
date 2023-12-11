package user_service_client

import (
	"context"
	"errors"

	"github.com/akatranlp/hsfl-master-ai-cloud-engineering/lib/grpc/user-service/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GRPCRepository struct {
	client proto.UserServiceClient
}

func NewGRPCRepository(client proto.UserServiceClient) *GRPCRepository {
	return &GRPCRepository{
		client: client,
	}
}

func (repo *GRPCRepository) MoveBalance(userId uint64, receivingUserId uint64, amount int64) error {
	req := &proto.MoveUserAmountRequest{
		UserId:          userId,
		ReceivingUserId: receivingUserId,
		Amount:          amount,
	}

	res, err := repo.client.MoveUserAmount(context.Background(), req)
	if err != nil {
		st, ok := status.FromError(err)
		if ok && st.Code() == codes.InvalidArgument {
			return errors.New("you cannot buy this book")
		}
		return err
	}

	if !res.Success {
		return errors.New("an unknown error")
	}

	return nil
}
