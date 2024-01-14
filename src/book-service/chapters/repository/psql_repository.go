package chapters_repository

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
	id			int not null,
	bookId		int not null,
	name    	varchar(100) not null,
	price		int not null,
	content 	text not null,
	status		int not null default 0,
   	foreign key (bookId) REFERENCES books(id),
	primary key (id, bookId)
)
`

func (repo *PsqlRepository) Migrate() error {
	_, err := repo.db.Exec(createChaptersTable)
	return err
}

const createChaptersBatchQuery = `
insert into chapters (id, bookId, name, price, content) values %s
`

const createChaptersHighestIdQuery = `
select max(id) from chapters where bookId = $1
`

func (repo *PsqlRepository) Create(chapters []*model.Chapter) error {
	placeholders := make([]string, len(chapters))
	values := make([]interface{}, len(chapters)*5)

	row := repo.db.QueryRow(createChaptersHighestIdQuery, chapters[0].BookID)

	var id int64 = 0
	err := row.Scan(&id) // ignore error

	if err != nil {
		fmt.Println(err.Error())
	}

	for i := 0; i < len(chapters); i++ {
		id++
		placeholders[i] = fmt.Sprintf("($%d,$%d,$%d,$%d,$%d)", i*5+1, i*5+2, i*5+3, i*5+4, i*5+5)
		values[i*5+0] = id
		values[i*5+1] = chapters[i].BookID
		values[i*5+2] = chapters[i].Name
		values[i*5+3] = chapters[i].Price
		values[i*5+4] = chapters[i].Content
	}

	query := fmt.Sprintf(createChaptersBatchQuery, strings.Join(placeholders, ","))
	_, err = repo.db.Exec(query, values...)
	return err
}

const updateChapterBatchQuery = `
update chapters set name = $1, price = $2, content = $3, status = $4 where id = $5 and bookId = $6
`

func (repo *PsqlRepository) Update(id uint64, bookId uint64, updateChapter *model.ChapterPatch) error {
	dbChapter, err := repo.FindByIdAndBookId(id, bookId)
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

	_, err = repo.db.Exec(updateChapterBatchQuery, dbChapter.Name, dbChapter.Price, dbChapter.Content, dbChapter.Status, id, bookId)
	return err
}

const findAllChaptersIdByBookIdQuery = `
select id, bookId, name, price, status from chapters where bookId = $1 ORDER BY id ASC
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
delete from chapters where (id, bookId) in (%s)
`

func (repo *PsqlRepository) Delete(chapters []*model.Chapter) error {
	placeholders := make([]string, len(chapters))
	ids := make([]interface{}, len(chapters)*2)
	for i := 0; i < len(chapters); i++ {
		placeholders[i] = fmt.Sprintf("($%d,$%d)", i*2+1, i*2+2)
		ids[i*2] = chapters[i].ID
		ids[i*2+1] = chapters[i].BookID
	}
	query := fmt.Sprintf(deleteChaptersBatchQuery, strings.Join(placeholders, ","))
	_, err := repo.db.Exec(query, ids...)
	return err
}

// find the chapter with bookId and chapterId and return the chapter with the author from the bookId in chapter
// only place that uses a inner join in the whole project, without it we could decouple the database completely
const validateChapterIdQuery = `
select c.id, c.bookId, c.name, c.price, c.content, c.status, b.authorId from chapters c inner join books b on c.bookId = b.id where c.id = $1 and c.bookId = $2
`

func (repo *PsqlRepository) ValidateChapterId(id uint64, bookId uint64) (*model.Chapter, *uint64, error) {
	row := repo.db.QueryRow(validateChapterIdQuery, id, bookId)

	var chapter model.Chapter
	var receivingUserId uint64
	if err := row.Scan(&chapter.ID, &chapter.BookID, &chapter.Name, &chapter.Price, &chapter.Content, &chapter.Status, &receivingUserId); err != nil {
		return nil, nil, err
	}

	return &chapter, &receivingUserId, nil
}
