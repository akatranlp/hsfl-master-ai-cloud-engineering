package service

import shared_types "github.com/akatranlp/hsfl-master-ai-cloud-engineering/lib/shared-types"

type Service interface {
	CheckChapterBought(userId uint64, chapterId uint64, bookId uint64) (bool, shared_types.Code, error)
}
