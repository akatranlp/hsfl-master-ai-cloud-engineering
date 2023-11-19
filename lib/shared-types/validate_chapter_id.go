package shared_types

type ValidateChapterIdRequest struct {
	UserId    uint64 `json:"userId"`
	ChapterId uint64 `json:"chapterId"`
}

func (r *ValidateChapterIdRequest) IsValid() bool {
	return r.UserId != 0 && r.ChapterId != 0
}


type ValidateChapterIdResponse struct {
	ChapterId       uint64 `json:"chapterId"`
	BookId          uint64 `json:"bookId"`
	ReceivingUserId uint64 `json:"receivingUserId"`
	Amount          uint64 `json:"amount"`
}
