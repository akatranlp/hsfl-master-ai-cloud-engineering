package transaction_service_client

type Repository interface {
	CheckChapterBought(userId uint64, chapterId uint64) error
}
