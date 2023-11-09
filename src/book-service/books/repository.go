package books

import "github.com/akatranlp/hsfl-master-ai-cloud-engineering/book-service/books/model"

type Repository interface {
	Migrate() error
	Create([]*model.Book) error
	Update(id uint64, updateBook *model.BookPatch) error
	FindAll() ([]*model.Book, error)
	FindAllByUserId(id uint64) ([]*model.Book, error)
	FindById(id uint64) (*model.Book, error)
	Delete([]*model.Book) error
}
