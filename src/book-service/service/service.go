package service

import shared_types "github.com/akatranlp/hsfl-master-ai-cloud-engineering/lib/shared-types"

type Service interface {
	ValidateChapterId(userId uint64, chapterId uint64, bookId uint64) (*shared_types.ValidateChapterIdResponse, shared_types.Code, error)
}
