package model

type Book struct {
	ID          uint64 `json:"id"`
	Name        string `json:"name"`
	AuthorID    uint64 `json:"authorId"`
	Description string `json:"description"`
}

type BookPatch struct {
	ID          *uint64
	Name        *string
	AuthorID    *uint64
	Description *string
}
