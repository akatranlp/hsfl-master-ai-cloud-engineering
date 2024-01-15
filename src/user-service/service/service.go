package service

import (
	shared_types "github.com/akatranlp/hsfl-master-ai-cloud-engineering/lib/shared-types"
	"github.com/akatranlp/hsfl-master-ai-cloud-engineering/user-service/model"
)

type Service interface {
	ValidateToken(token string) (*model.DbUser, shared_types.Code, error)
	MoveUserAmount(payingUserId uint64, receivingUserId uint64, amount int64) (shared_types.Code, error)
}
