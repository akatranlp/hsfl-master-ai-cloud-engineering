package books

import (
	"context"
	"database/sql"
	"testing"

	"github.com/akatranlp/hsfl-master-ai-cloud-engineering/book-service/books/model"
	"github.com/akatranlp/hsfl-master-ai-cloud-engineering/lib/containerhelpers"
	"github.com/akatranlp/hsfl-master-ai-cloud-engineering/lib/database"
	"github.com/stretchr/testify/assert"
)

func TestIntegrationPsqlBookRepository(t *testing.T) {
	postgres, err := containerhelpers.StartPostgres()
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

	dbConfig := database.PsqlConfig{
		Host:     "0.0.0.0",
		Port:     port.Int(),
		Username: "postgres",
		Password: "postgres",
		Database: "postgres",
	}

	repository, err := NewPsqlRepository(dbConfig)
	if err != nil {
		t.Fatalf("could not create book repository: %s", err.Error())
	}
	t.Cleanup(clearBookTables(t, repository.db))

	t.Run("Migrate", func(t *testing.T) {
		t.Run("should create books table", func(t *testing.T) {
			t.Cleanup(clearBookTables(t, repository.db))

			// given
			createUserTable(t, repository.db)

			// when
			err := repository.Migrate()
			// then
			assert.NoError(t, err)
			assertBooksTableExists(t, repository.db, "books", []string{"id", "name", "authorid", "description"})
		})
	})

	t.Run("Create", func(t *testing.T) {
		t.Run("should insert books in batches", func(t *testing.T) {
			t.Cleanup(clearBookTables(t, repository.db))

			// given
			books := []*model.Book{
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

			// when
			err := repository.Create(books)

			// then
			assert.NoError(t, err)
			assert.Equal(t, books[0], getBookFromDatabase(t, repository.db, 1))
			assert.Equal(t, books[1], getBookFromDatabase(t, repository.db, 2))
		})
	})
	t.Run("Update", func(t *testing.T) {
		t.Run("should update book", func(t *testing.T) {
			t.Cleanup(clearBookTables(t, repository.db))

			// given
			insertBook(t, repository.db, &model.Book{
				ID:          3,
				Name:        "Book One",
				AuthorID:    1,
				Description: "A good book",
			})
			name := "Updated Book"
			desc := "An updated description"
			newBookData := &model.BookPatch{
				Name:        &name,
				Description: &desc,
			}

			// when
			err := repository.Update(3, newBookData)

			// then
			assert.NoError(t, err)
			assert.Equal(t, &model.Book{
				ID:          3,
				Name:        "Updated Book",
				AuthorID:    1,
				Description: "An updated description",
			}, getBookFromDatabase(t, repository.db, 3))
		})
	})
	t.Run("FindById", func(t *testing.T) {
		t.Run("should return book", func(t *testing.T) {
			t.Cleanup(clearBookTables(t, repository.db))

			// given
			insertBook(t, repository.db, &model.Book{
				ID:          4,
				Name:        "Book One",
				AuthorID:    1,
				Description: "A good book",
			})

			// when
			book, err := repository.FindById(4)

			// then
			assert.NoError(t, err)
			assert.NotNil(t, book)
		})
	})

	t.Run("Delete", func(t *testing.T) {
		t.Run("should delete provided books", func(t *testing.T) {
			t.Cleanup(clearBookTables(t, repository.db))

			// given
			books := []*model.Book{
				{
					ID:          5,
					Name:        "Book One",
					AuthorID:    1,
					Description: "A good book",
				},
				{
					ID:          6,
					Name:        "Book Two",
					AuthorID:    1,
					Description: "A bad book",
				},
			}

			for _, book := range books {
				insertBook(t, repository.db, book)
			}

			// when
			err := repository.Delete([]*model.Book{books[1]})

			// then
			assert.NoError(t, err)
			assert.Equal(t, books[0], getBookFromDatabase(t, repository.db, 5))
			assert.Nil(t, getBookFromDatabase(t, repository.db, 6))
		})
	})

	t.Run("FindAll", func(t *testing.T) {
		t.Run("should get all books", func(t *testing.T) {
			t.Cleanup(clearBookTables(t, repository.db))

			// given
			books := []*model.Book{
				{
					ID:          7,
					Name:        "Book One",
					AuthorID:    1,
					Description: "A good book",
				},
				{
					ID:          8,
					Name:        "Book Two",
					AuthorID:    1,
					Description: "A bad book",
				},
				{
					ID:          9,
					Name:        "Book Three",
					AuthorID:    2,
					Description: "A middle book",
				},
			}

			for _, book := range books {
				insertBook(t, repository.db, book)
			}

			// when
			result, err := repository.FindAll()

			// then
			assert.NoError(t, err)
			assert.Equal(t, books, result)
		})
	})

	t.Run("FindAll By UserId", func(t *testing.T) {
		t.Run("should get all books from author 1", func(t *testing.T) {
			t.Cleanup(clearBookTables(t, repository.db))

			// given
			books := []*model.Book{
				{
					ID:          10,
					Name:        "Book One",
					AuthorID:    1,
					Description: "A good book",
				},
				{
					ID:          11,
					Name:        "Book Two",
					AuthorID:    1,
					Description: "A bad book",
				},
				{
					ID:          12,
					Name:        "Book Three",
					AuthorID:    2,
					Description: "A middle book",
				},
			}

			for _, book := range books {
				insertBook(t, repository.db, book)
			}

			// when
			result, err := repository.FindAllByUserId(1)

			// then
			assert.NoError(t, err)
			assert.Equal(t, []*model.Book{books[0], books[1]}, result)
		})
	})

}

func createUserTable(t *testing.T, db *sql.DB) {
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
	db.Exec(`insert into users (email, username, password, profile_name) values ($1,$2,$3,$4)`, "1a@mail.com", "tester1", []byte("pw"), "Peter")
	db.Exec(`insert into users (email, username, password, profile_name) values ($1,$2,$3,$4)`, "2a@mail.com", "tester2", []byte("pw"), "Ursula")
}

func getBookFromDatabase(t *testing.T, db *sql.DB, id int) *model.Book {
	row := db.QueryRow(`select id, name, authorId, description from books where id = $1`, id)

	var book model.Book
	if err := row.Scan(&book.ID, &book.Name, &book.AuthorID, &book.Description); err != nil {
		return nil
	}

	return &book
}

func insertBook(t *testing.T, db *sql.DB, book *model.Book) {
	_, err := db.Exec(`insert into books (name, authorId, description) values ($1, $2, $3)`, book.Name, book.AuthorID, book.Description)
	if err != nil {
		t.Logf("could not insert book: %s", err.Error())
		t.FailNow()
	}
}

func clearBookTables(t *testing.T, db *sql.DB) func() {
	return func() {
		if _, err := db.Exec("delete from books"); err != nil {
			t.Logf("could not delete rows from books: %s", err.Error())
			t.FailNow()
		}
	}
}

func assertBooksTableExists(t *testing.T, db *sql.DB, name string, columns []string) {
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
