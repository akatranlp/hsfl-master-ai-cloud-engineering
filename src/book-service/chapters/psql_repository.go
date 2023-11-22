package chapters

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/akatranlp/hsfl-master-ai-cloud-engineering/book-service/chapters/model"
	"github.com/akatranlp/hsfl-master-ai-cloud-engineering/lib/database"
	_ "github.com/lib/pq"
)

type PsqlRepository struct {
	db *sql.DB
}

func NewPsqlRepository(config database.Config) (*PsqlRepository, error) {
	dsn := config.Dsn()
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	return &PsqlRepository{db}, nil
}

const createChaptersTable = `
create table if not exists chapters (
    id			serial primary key,
    bookId		int not null,
	name    	varchar(100) not null,
	price		int not null,
	content 	text not null,
	status		int not null default 0,
   	foreign key (bookId) REFERENCES books(id)
)
`

func (repo *PsqlRepository) Migrate() error {
	_, err := repo.db.Exec(createChaptersTable)
	return err
}

const createChaptersBatchQuery = `
insert into chapters (bookId, name, price, content) values %s
`

func (repo *PsqlRepository) Create(chapters []*model.Chapter) error {
	placeholders := make([]string, len(chapters))
	values := make([]interface{}, len(chapters)*4)

	for i := 0; i < len(chapters); i++ {
		placeholders[i] = fmt.Sprintf("($%d,$%d,$%d,$%d)", i*4+1, i*4+2, i*4+3, i*4+4)
		values[i*4+0] = chapters[i].BookID
		values[i*4+1] = chapters[i].Name
		values[i*4+2] = chapters[i].Price
		values[i*4+3] = chapters[i].Content
	}

	query := fmt.Sprintf(createChaptersBatchQuery, strings.Join(placeholders, ","))
	_, err := repo.db.Exec(query, values...)
	return err
}

const updateChapterBatchQuery = `
update chapters set name = $1, price = $2, content = $3, status = $4 where id = $5
`

func (repo *PsqlRepository) Update(id uint64, updateChapter *model.ChapterPatch) error {
	dbChapter, err := repo.FindById(id)
	if err != nil {
		return err
	}
	if updateChapter.Name != nil {
		dbChapter.Name = *updateChapter.Name
	}
	if updateChapter.Price != nil {
		dbChapter.Price = *updateChapter.Price
	}
	if updateChapter.Content != nil {
		dbChapter.Content = *updateChapter.Content
	}
	if updateChapter.Status != nil {
		dbChapter.Status = *updateChapter.Status
	}

	_, err = repo.db.Exec(updateChapterBatchQuery, dbChapter.Name, dbChapter.Price, dbChapter.Content, dbChapter.Status , id)
	return err
}
const findAllChaptersIdByBookIdQuery = `
select id, bookId, name, price, status from chapters where bookId = $1
`

func (repo *PsqlRepository) FindAllPreviewsByBookId(bookId uint64) ([]*model.ChapterPreview, error) {
	rows, err := repo.db.Query(findAllChaptersIdByBookIdQuery, bookId)
	if err != nil {
		return nil, err
	}

	chapters := make([]*model.ChapterPreview, 0)
	for rows.Next() {
		chapter := model.ChapterPreview{}
		if err := rows.Scan(&chapter.ID, &chapter.BookID, &chapter.Name, &chapter.Price, &chapter.Status); err != nil {
			return nil, err
		}
		chapters = append(chapters, &chapter)
	}

	return chapters, nil
}

const findChapterByIDQuery = `
select id, bookId, name, price, content, status from chapters where id = $1
`

func (repo *PsqlRepository) FindById(id uint64) (*model.Chapter, error) {
	row := repo.db.QueryRow(findChapterByIDQuery, id)

	var chapter model.Chapter
	if err := row.Scan(&chapter.ID, &chapter.BookID, &chapter.Name, &chapter.Price, &chapter.Content, &chapter.Status); err != nil {
		return nil, err
	}

	return &chapter, nil
}

const findChapterByIdAndBookIdQuery = `
select id, bookId, name, price, content, status from chapters where id = $1 and bookId = $2
`

func (repo *PsqlRepository) FindByIdAndBookId(id uint64, bookId uint64) (*model.Chapter, error) {
	row := repo.db.QueryRow(findChapterByIdAndBookIdQuery, id, bookId)

	var chapter model.Chapter
	if err := row.Scan(&chapter.ID, &chapter.BookID, &chapter.Name, &chapter.Price, &chapter.Content, &chapter.Status); err != nil {
		return nil, err
	}

	return &chapter, nil
}

const deleteChaptersBatchQuery = `
delete from chapters where id in (%s)
`

func (repo *PsqlRepository) Delete(chapters []*model.Chapter) error {
	placeholders := make([]string, len(chapters))
	ids := make([]interface{}, len(chapters))

	for i := 0; i < len(chapters); i++ {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
		ids[i] = chapters[i].ID
	}

	query := fmt.Sprintf(deleteChaptersBatchQuery, strings.Join(placeholders, ","))
	_, err := repo.db.Exec(query, ids...)
	return err
}

const validateChapterIdQuery = `
select c.id, c.bookId, c.name, c.price, c.content, c.status, b.authorId from chapters as c inner join books as b on c.bookId = b.id where c.id = $1
`

func (repo *PsqlRepository) ValidateChapterId(id uint64) (*model.Chapter, *uint64, error) {
	row := repo.db.QueryRow(validateChapterIdQuery, id)

	var chapter model.Chapter
	var receivingUserId uint64
	if err := row.Scan(&chapter.ID, &chapter.BookID, &chapter.Name, &chapter.Price, &chapter.Content, &chapter.Status, &receivingUserId); err != nil {
		return nil, nil, err
	}

	return &chapter, &receivingUserId, nil
}
