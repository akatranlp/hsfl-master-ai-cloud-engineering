package chapters_repository

import "github.com/akatranlp/hsfl-master-ai-cloud-engineering/book-service/chapters/model"

type Repository interface {
	Migrate() error
	Create([]*model.Chapter) error
	Update(id uint64, bookId uint64, updateChapter *model.ChapterPatch) error
	FindAllPreviewsByBookId(bookId uint64) ([]*model.ChapterPreview, error)
	FindByIdAndBookId(id uint64, bookId uint64) (*model.Chapter, error)
	ValidateChapterId(id uint64, bookId uint64) (*model.Chapter, *uint64, error)
	Delete([]*model.Chapter) error
}
