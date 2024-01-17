package grpc

import (
	"context"

	"github.com/akatranlp/hsfl-master-ai-cloud-engineering/lib/grpc/user-service/proto"
	"github.com/akatranlp/hsfl-master-ai-cloud-engineering/user-service/service"
	"google.golang.org/grpc/status"
)

type server struct {
	proto.UnimplementedUserServiceServer
	service service.Service
}

func NewServer(
	service service.Service,
) proto.UserServiceServer {
	return &server{
		service: service,
	}
}

func (s *server) ValidateToken(ctx context.Context, req *proto.ValidateTokenRequest) (*proto.ValidateTokenResponse, error) {
	user, statusCode, err := s.service.ValidateAccessToken(req.Token)
	if user == nil {
		return nil, status.Error(statusCode.ToGRPCStatusCode(), err.Error())
	}

	response := &proto.ValidateTokenResponse{
		Success: true,
		UserId:  user.ID,
	}
	return response, nil
}

func (s *server) MoveUserAmount(ctx context.Context, req *proto.MoveUserAmountRequest) (*proto.MoveUserAmountResponse, error) {
	statusCode, err := s.service.MoveUserAmount(req.UserId, req.ReceivingUserId, req.Amount)
	if err != nil {
		return nil, status.Error(statusCode.ToGRPCStatusCode(), err.Error())
	}

	response := &proto.MoveUserAmountResponse{
		Success: true,
	}
	return response, nil
}
