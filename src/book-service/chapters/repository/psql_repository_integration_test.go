package chapters_repository

import (
	"context"
	"database/sql"
	"testing"

	booksModel "github.com/akatranlp/hsfl-master-ai-cloud-engineering/book-service/books/model"
	"github.com/akatranlp/hsfl-master-ai-cloud-engineering/book-service/chapters/model"
	"github.com/akatranlp/hsfl-master-ai-cloud-engineering/lib/containerhelpers"
	"github.com/akatranlp/hsfl-master-ai-cloud-engineering/lib/database"
	"github.com/stretchr/testify/assert"
)

func TestIntegrationPsqlChapterRepository(t *testing.T) {
	postgres, err := containerhelpers.StartPostgres(false)
	if err != nil {
		t.Fatalf("could not start postgres container: %s", err.Error())
	}

	t.Cleanup(func() {
		postgres.Terminate(context.Background())
	})

	port, err := postgres.MappedPort(context.Background(), "5432")
	if err != nil {
		t.Fatalf("could not get database container port: %s", err.Error())
	}

	repository, err := NewPsqlRepository(database.PsqlConfig{
		Host:     "0.0.0.0",
		Port:     port.Int(),
		Username: "postgres",
		Password: "postgres",
		Database: "postgres",
	})
	if err != nil {
		t.Fatalf("could not create chapter repository: %s", err.Error())
	}
	t.Cleanup(clearChapterTables(t, repository.db))

	t.Run("Migrate", func(t *testing.T) {
		t.Run("should create chapters table", func(t *testing.T) {
			t.Cleanup(clearChapterTables(t, repository.db))

			// given
			createUserAndBookTable(t, repository.db)
			// when
			err := repository.Migrate()

			// then
			assert.NoError(t, err)
			assertChapterTableExists(t, repository.db, "chapters", []string{"id", "bookid", "name", "price", "content", "status"})
		})
	})

	t.Run("Create", func(t *testing.T) {
		t.Run("should insert chapters in batches", func(t *testing.T) {
			t.Cleanup(clearChapterTables(t, repository.db))

			// given
			chapters := []*model.Chapter{
				{
					ID:      1,
					BookID:  1,
					Name:    "doesnt matter",
					Price:   0,
					Content: "doesnt matter",
				},
				{
					ID:      2,
					BookID:  1,
					Name:    "doesnt matter",
					Price:   0,
					Content: "doesnt matter",
				},
			}

			// when
			err := repository.Create(chapters)

			// then
			assert.NoError(t, err)
			assert.Equal(t, chapters[0], getChapterFromDatabase(t, repository.db, 1))
			assert.Equal(t, chapters[1], getChapterFromDatabase(t, repository.db, 2))
		})
	})
	t.Run("Update", func(t *testing.T) {
		t.Run("should update chapter", func(t *testing.T) {
			t.Cleanup(clearChapterTables(t, repository.db))

			// given
			insertChapter(t, repository.db, &model.Chapter{
				ID:      3,
				BookID:  1,
				Name:    "First Chapter",
				Price:   0,
				Content: "Chapter Content",
			})
			price := uint64(100)
			content := "Updated Chapter: Chapter Content but good now"
			newChapterData := &model.ChapterPatch{
				Price:   &price,
				Content: &content,
			}

			// when
			err := repository.Update(3, 1, newChapterData)

			// then
			assert.NoError(t, err)
			assert.Equal(t, &model.Chapter{
				ID:      3,
				BookID:  1,
				Name:    "First Chapter",
				Price:   100,
				Content: "Updated Chapter: Chapter Content but good now",
			}, getChapterFromDatabase(t, repository.db, 3))
		})
	})

	t.Run("Delete", func(t *testing.T) {
		t.Run("should delete provided chapters", func(t *testing.T) {
			t.Cleanup(clearChapterTables(t, repository.db))

			// given
			chapters := []*model.Chapter{
				{
					ID:      5,
					BookID:  1,
					Name:    "doesnt matter",
					Price:   0,
					Content: "doesnt matter",
				},
				{
					ID:      6,
					BookID:  1,
					Name:    "doesnt matter",
					Price:   0,
					Content: "doesnt matter",
				},
			}

			for _, chapter := range chapters {
				insertChapter(t, repository.db, chapter)
			}

			// when
			err := repository.Delete([]*model.Chapter{chapters[1]})

			// then
			assert.NoError(t, err)
			assert.Equal(t, chapters[0], getChapterFromDatabase(t, repository.db, 5))
			assert.Nil(t, getChapterFromDatabase(t, repository.db, 6))
		})
	})
}

func createUserAndBookTable(t *testing.T, db *sql.DB) {
	db.Exec(`
		create table if not exists users (
			id				serial primary key,
			email			varchar(100) not null unique,
			username    	varchar(16) not null unique,
			password 		bytea not null,
			profile_name 	varchar(100) not null,
			balance 		int not null default 0
		)
	`)

	db.Exec(`
		create table if not exists books (
			id			serial primary key,
			name    	varchar(100) not null,
			authorId	int not null,
			description text not null,
			foreign key (authorId) REFERENCES users(id)
		)
	`)
	db.Exec(`insert into users (email, username, password, profile_name) values ($1,$2,$3,$4)`, "1a@mail.com", "tester1", []byte("pw"), "Peter")
	db.Exec(`insert into users (email, username, password, profile_name) values ($1,$2,$3,$4)`, "2a@mail.com", "tester2", []byte("pw"), "Ursula")

	books := []*booksModel.Book{
		{
			ID:          1,
			Name:        "Book One",
			AuthorID:    1,
			Description: "A good book",
		},
		{
			ID:          2,
			Name:        "Book Two",
			AuthorID:    1,
			Description: "A bad book",
		},
	}
	for _, book := range books {
		insertBook(t, db, book)
	}

	assert.Equal(t, books[0], getBookFromDatabase(t, db, 1))
	assert.Equal(t, books[1], getBookFromDatabase(t, db, 2))
}

func getChapterFromDatabase(t *testing.T, db *sql.DB, id int) *model.Chapter {
	row := db.QueryRow(`select id,bookId,name,price,content from chapters where id = $1`, id)

	var chapter model.Chapter
	if err := row.Scan(&chapter.ID, &chapter.BookID, &chapter.Name, &chapter.Price, &chapter.Content); err != nil {
		return nil
	}

	return &chapter
}

func insertChapter(t *testing.T, db *sql.DB, chapter *model.Chapter) {
	_, err := db.Exec(`insert into chapters (id,bookId,name,price,content) values ($1, $2, $3, $4, $5)`, chapter.ID, chapter.BookID, chapter.Name, chapter.Price, chapter.Content)
	if err != nil {
		t.Logf("could not insert chapter: %s", err.Error())
		t.FailNow()
	}
}

func clearChapterTables(t *testing.T, db *sql.DB) func() {
	return func() {
		if _, err := db.Exec("delete from chapters"); err != nil {
			t.Logf("could not delete rows from chapters: %s", err.Error())
			t.FailNow()
		}
	}
}

func assertChapterTableExists(t *testing.T, db *sql.DB, name string, columns []string) {
	rows, err := db.Query(`select column_name from information_schema.columns where table_name = $1`, name)
	if err != nil {
		t.Fail()
		return
	}

	scannedCols := make(map[string]struct{})
	for rows.Next() {
		var column string
		if err := rows.Scan(&column); err != nil {
			t.Logf("expected")
			t.FailNow()
		}

		scannedCols[column] = struct{}{}
	}

	if len(scannedCols) == 0 {
		t.Logf("expected table '%s' to exist, but not found", name)
		t.FailNow()
	}

	for _, col := range columns {
		if _, ok := scannedCols[col]; !ok {
			t.Logf("expected table '%s' to have column '%s'", name, col)
			t.Fail()
		}
	}
}

//___________________________

func getBookFromDatabase(t *testing.T, db *sql.DB, id int) *booksModel.Book {
	row := db.QueryRow(`select id, name, authorId, description from books where id = $1`, id)

	var book booksModel.Book
	if err := row.Scan(&book.ID, &book.Name, &book.AuthorID, &book.Description); err != nil {
		return nil
	}

	return &book
}

func insertBook(t *testing.T, db *sql.DB, book *booksModel.Book) {
	_, err := db.Exec(`insert into books (name, authorId, description) values ($1, $2, $3)`, book.Name, book.AuthorID, book.Description)
	if err != nil {
		t.Logf("could not insert book: %s", err.Error())
		t.FailNow()
	}
}
