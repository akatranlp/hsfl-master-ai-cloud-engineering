package transactions_repository

import "github.com/akatranlp/hsfl-master-ai-cloud-engineering/transaction-service/transactions/model"

type Repository interface {
	Migrate() error
	Create([]*model.Transaction) error
	FindForUserIdAndChapterId(userId uint64, chapterId uint64, bookId uint64) (*model.Transaction, error)
	FindAllForUserId(userId uint64) ([]*model.Transaction, error)
	FindAllForReceivingUserId(userId uint64) ([]*model.Transaction, error)
}
