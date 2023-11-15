package book_service_client

import shared_types "github.com/akatranlp/hsfl-master-ai-cloud-engineering/lib/shared-types"

type Repository interface {
	ValidateChapterId(userId uint64, chapterId uint64) (*shared_types.ValidateChapterIdResponse, error)
}
