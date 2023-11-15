package shared_types

type ValidateChapterIdResponse struct {
	ChapterId       uint64 `json:"chapterId"`
	BookId          uint64 `json:"bookId"`
	ReceivingUserId uint64 `json:"receivingUserId"`
	Amount          uint64 `json:"amount"`
}
