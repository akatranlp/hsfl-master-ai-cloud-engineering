package service

import (
	"errors"
	"log"

	shared_types "github.com/akatranlp/hsfl-master-ai-cloud-engineering/lib/shared-types"
	"github.com/akatranlp/hsfl-master-ai-cloud-engineering/user-service/auth"
	"github.com/akatranlp/hsfl-master-ai-cloud-engineering/user-service/model"
	"github.com/akatranlp/hsfl-master-ai-cloud-engineering/user-service/repository"
)

type DefaultService struct {
	repository            repository.Repository
	accessTokenGenerator  auth.TokenGenerator
	refreshTokenGenerator auth.TokenGenerator
	authIsActive          bool
}

func NewDefaultService(
	repository repository.Repository,
	accessTokenGenerator auth.TokenGenerator,
	refreshTokenGenerator auth.TokenGenerator,
	authIsActive bool,
) *DefaultService {
	return &DefaultService{
		repository:            repository,
		accessTokenGenerator:  accessTokenGenerator,
		refreshTokenGenerator: refreshTokenGenerator,
		authIsActive:          authIsActive,
	}
}

func (s *DefaultService) validateToken(token string, tokenGenerator auth.TokenGenerator) (*model.DbUser, shared_types.Code, error) {
	if !s.authIsActive {
		user, err := s.repository.FindById(1)
		if err != nil {
			log.Println("ERROR [tokenVerification - FindById]: ", err.Error())
			return nil, shared_types.NotFound, errors.New("user not found")
		}
		return user, shared_types.OK, nil
	}

	claims, err := tokenGenerator.VerifyToken(token)
	if err != nil {
		log.Println("ERROR [tokenVerification - VerifyToken]: ", err.Error())
		return nil, shared_types.Unauthenticated, errors.New("token couldn't be verified")
	}

	email, ok := claims["email"].(string)
	if !ok {
		log.Println("ERROR [tokenVerification - get email claim]: ", "There is no email claim in your token")
		return nil, shared_types.Unauthenticated, errors.New("there is no email claim in your token")
	}

	tokenV, ok := claims["token_version"].(float64)
	if !ok {
		log.Println("ERROR [tokenVerification - get token_version claim]: ", "There is no token_version claim in your token")
		return nil, shared_types.Unauthenticated, errors.New("there is no token_version claim in your token")
	}
	tokenVersion := uint64(tokenV)

	users, err := s.repository.FindByEmail(email)
	if err != nil {
		log.Println("ERROR [tokenVerification - FindByEmail]: ", err.Error())
		return nil, shared_types.Internal, errors.New("internal server error")
	}

	if len(users) < 1 {
		log.Println("ERROR [tokenVerification - len(users) < 1]: ", "Couldn't find user by email")
		return nil, shared_types.Unauthenticated, errors.New("couldn't find user by email")
	}

	if users[0].TokenVersion != tokenVersion {
		log.Println("ERROR [tokenVerification - token version]: ", "The token version is not valid")
		return nil, shared_types.Unauthenticated, errors.New("the token version is not valid")
	}

	return users[0], shared_types.OK, nil
}

func (s *DefaultService) ValidateAccessToken(token string) (*model.DbUser, shared_types.Code, error) {
	return s.validateToken(token, s.accessTokenGenerator)
}

func (s *DefaultService) ValidateRefreshToken(token string) (*model.DbUser, shared_types.Code, error) {
	return s.validateToken(token, s.refreshTokenGenerator)
}

func (s *DefaultService) MoveUserAmount(payingUserId uint64, receivingUserId uint64, amount int64) (shared_types.Code, error) {
	payingUser, err := s.repository.FindById(payingUserId)
	if err != nil {
		log.Println("ERROR [MoveUserAmount - FindById - Paying]: ", err.Error())
		return shared_types.NotFound, errors.New("payingUser not found")
	}

	receivingUser, err := s.repository.FindById(receivingUserId)
	if err != nil {
		log.Println("ERROR [MoveUserAmount - FindById - Receiving]: ", err.Error())
		return shared_types.NotFound, errors.New("receivingUser not found")
	}

	payingUserBalance := payingUser.Balance - amount
	receivingUserBalance := receivingUser.Balance + amount

	userPatch := &model.DbUserPatch{Balance: &payingUserBalance}
	err = s.repository.Update(payingUser.ID, userPatch)
	if err != nil {
		log.Println("ERROR [MoveUserAmount - Update - Paying]: ", err.Error())
		return shared_types.Internal, errors.New("internal server error")
	}

	userPatch = &model.DbUserPatch{Balance: &receivingUserBalance}
	err = s.repository.Update(receivingUser.ID, userPatch)
	if err != nil {
		log.Println("ERROR [MoveUserAmount - Update - Receiving]: ", err.Error())
		return shared_types.Internal, errors.New("internal server error")
	}

	return shared_types.OK, nil
}
