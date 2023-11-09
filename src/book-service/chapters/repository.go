package chapters

import "github.com/akatranlp/hsfl-master-ai-cloud-engineering/book-service/chapters/model"

type Repository interface {
	Migrate() error
	Create([]*model.Chapter) error
	Update(id uint64, updateChapter *model.ChapterPatch) error
	FindAllPreviewsByBookId(bookId uint64) ([]*model.ChapterPreview, error)
	FindById(id uint64) (*model.Chapter, error)
	FindByIdAndBookId(id uint64, bookId uint64) (*model.Chapter, error)
	Delete([]*model.Chapter) error
}
