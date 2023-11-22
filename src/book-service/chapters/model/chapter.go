package model

type Status int

const (
	Draft Status = iota
	Published
)

type Chapter struct {
	ID      uint64 `json:"id"`
	BookID  uint64 `json:"bookid"`
	Name    string `json:"name"`
	Price   uint64 `json:"price"`
	Content string `json:"content"`
	Status  Status `json:"status"`
}

type ChapterPreview struct {
	ID     uint64 `json:"id"`
	BookID uint64 `json:"bookid"`
	Name   string `json:"name"`
	Price  uint64 `json:"price"`
	Status  Status `json:"status"`
}

type ChapterPatch struct {
	ID      *uint64
	BookID  *uint64
	Name    *string
	Price   *uint64
	Content *string
	Status  *Status
}
