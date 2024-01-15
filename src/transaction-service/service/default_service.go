package service

import (
	shared_types "github.com/akatranlp/hsfl-master-ai-cloud-engineering/lib/shared-types"
	repository "github.com/akatranlp/hsfl-master-ai-cloud-engineering/transaction-service/repository"
)

type DefaultService struct {
	transactionRepository repository.Repository
}

func NewDefaultService(transactionRepository repository.Repository) *DefaultService {
	return &DefaultService{
		transactionRepository: transactionRepository,
	}
}

func (s *DefaultService) CheckChapterBought(userId uint64, chapterId uint64, bookId uint64) (bool, shared_types.Code, error) {
	transaction, err := s.transactionRepository.FindForUserIdAndChapterId(userId, chapterId, bookId)
	if err != nil {
		return false, shared_types.NotFound, err
	}

	return transaction != nil, shared_types.OK, nil
}
