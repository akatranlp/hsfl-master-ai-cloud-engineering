package repository

import "github.com/akatranlp/hsfl-master-ai-cloud-engineering/user-service/model"

type Repository interface {
	Migrate() error
	Create([]*model.DbUser) error
	FindAll() ([]*model.DbUser, error)
	FindByEmail(email string) ([]*model.DbUser, error)
	FindById(id uint64) (*model.DbUser, error)
	Update(id uint64, user *model.DbUserPatch) error
	Delete([]*model.DbUser) error
}
