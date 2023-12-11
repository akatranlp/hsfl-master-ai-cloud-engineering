package grpc

import (
	"context"
	"errors"
	"log"

	"github.com/akatranlp/hsfl-master-ai-cloud-engineering/lib/grpc/user-service/proto"
	"github.com/akatranlp/hsfl-master-ai-cloud-engineering/user-service/auth"
	"github.com/akatranlp/hsfl-master-ai-cloud-engineering/user-service/user"
	"github.com/akatranlp/hsfl-master-ai-cloud-engineering/user-service/user/model"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type server struct {
	proto.UnimplementedUserServiceServer
	userRepository user.Repository
	tokenGenerator auth.TokenGenerator
	authIsActive   bool
}

func NewServer(
	userRepository user.Repository,
	tokenGenerator auth.TokenGenerator,
	authIsActive bool,
) proto.UserServiceServer {
	return &server{
		userRepository: userRepository,
		authIsActive:   authIsActive,
		tokenGenerator: tokenGenerator,
	}
}

func (s *server) tokenVerification(token string) (*model.DbUser, codes.Code, error) {
	claims, err := s.tokenGenerator.VerifyToken(token)
	if err != nil {
		log.Println("ERROR [tokenVerification - VerifyToken]: ", err.Error())
		return nil, codes.Unauthenticated, errors.New("token couldn't be verified")
	}

	email, ok := claims["email"].(string)
	if !ok {
		log.Println("ERROR [tokenVerification - get email claim]: ", "There is no email claim in your token")
		return nil, codes.Unauthenticated, errors.New("there is no email claim in your token")
	}

	tokenV, ok := claims["token_version"].(float64)
	if !ok {
		log.Println("ERROR [tokenVerification - get token_version claim]: ", "There is no token_version claim in your token")
		return nil, codes.Unauthenticated, errors.New("there is no token_version claim in your token")
	}
	tokenVersion := uint64(tokenV)

	users, err := s.userRepository.FindByEmail(email)
	if err != nil {
		log.Println("ERROR [tokenVerification - FindByEmail]: ", err.Error())
		return nil, codes.Internal, errors.New("internal server error")
	}

	if len(users) < 1 {
		log.Println("ERROR [tokenVerification - len(users) < 1]: ", "Couldn't find user by email")
		return nil, codes.Unauthenticated, errors.New("couldn't find user by email")
	}

	if users[0].TokenVersion != tokenVersion {
		log.Println("ERROR [tokenVerification - token version]: ", "The token version is not valid")
		return nil, codes.Unauthenticated, errors.New("the token version is not valid")
	}

	return users[0], codes.OK, nil
}

func (s *server) ValidateToken(ctx context.Context, req *proto.ValidateTokenRequest) (*proto.ValidateTokenResponse, error) {
	if !s.authIsActive {
		response := &proto.ValidateTokenResponse{
			Success: true,
			UserId:  uint64(1),
		}
		return response, nil
	}

	user, statusCode, err := s.tokenVerification(req.Token)
	if user == nil {
		return nil, status.Error(statusCode, err.Error())
	}

	response := &proto.ValidateTokenResponse{
		Success: true,
		UserId:  user.ID,
	}
	return response, nil
}

func (s *server) MoveUserAmount(ctx context.Context, req *proto.MoveUserAmountRequest) (*proto.MoveUserAmountResponse, error) {
	// Fully implement this if we need Authentication ????

	payingUser, err := s.userRepository.FindById(req.UserId)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	receivingUser, err := s.userRepository.FindById(req.ReceivingUserId)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	payingUserBalance := payingUser.Balance - req.Amount
	receivingUserBalance := receivingUser.Balance + req.Amount

	userPatch := &model.DbUserPatch{Balance: &payingUserBalance}
	err = s.userRepository.Update(payingUser.ID, userPatch)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	userPatch = &model.DbUserPatch{Balance: &receivingUserBalance}
	err = s.userRepository.Update(receivingUser.ID, userPatch)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	response := &proto.MoveUserAmountResponse{
		Success: true,
	}
	return response, nil
}
