package auth_middleware

import (
	"context"
	"errors"

	"github.com/akatranlp/hsfl-master-ai-cloud-engineering/lib/grpc/user-service/proto"
)

type GRPCRepository struct {
	client proto.UserServiceClient
}

func NewGRPCRepository(client proto.UserServiceClient) *GRPCRepository {
	return &GRPCRepository{
		client: client,
	}
}

func (repo *GRPCRepository) VerifyToken(token string) (uint64, error) {
	req := &proto.ValidateTokenRequest{
		Token: token,
	}

	res, err := repo.client.ValidateToken(context.Background(), req)
	if err != nil {
		return 0, err
	}

	if !res.Success {
		return 0, errors.New("an unknown error")
	}

	return res.UserId, nil
}
