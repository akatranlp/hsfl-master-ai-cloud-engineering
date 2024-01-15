package model

type Transaction struct {
	ID              uint64 `json:"id"`
	ChapterID       uint64 `json:"chapterID"`
	BookID          uint64 `json:"bookID"`
	PayingUserID    uint64 `json:"payingUserID"`
	ReceivingUserID uint64 `json:"receivingUserID"`
	Amount          uint64 `json:"amount"`
}
