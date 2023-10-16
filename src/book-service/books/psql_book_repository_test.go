package books

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/akatranlp/hsfl-master-ai-cloud-engineering/book-service/books/model"
	"github.com/stretchr/testify/assert"
)

func TestPsqlBookRepository(t *testing.T) {
	db, dbmock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}

	repository := PsqlBookRepository{db}

	t.Run("Create", func(t *testing.T) {
		t.Run("should return error if executing query failed", func(t *testing.T) {
			// given
			books := []*model.Book{{
				ID:          1,
				Name:        "Book One",
				AuthorID:    1,
				Description: "An ok book",
			}}

			dbmock.
				ExpectExec(`insert into books`).
				WillReturnError(errors.New("database error"))

			// when
			err := repository.Create(books)

			// then
			assert.Error(t, err)
		})

		t.Run("should insert books in batches", func(t *testing.T) {
			// given
			books := []*model.Book{
				{
					ID:          2,
					Name:        "Book Two",
					AuthorID:    1,
					Description: "A good book",
				},
				{
					ID:          3,
					Name:        "Book Three",
					AuthorID:    1,
					Description: "A bad book",
				},
			}

			dbmock.
				ExpectExec(`insert into books \(name, authorId, description\) values \(\$1,\$2,\$3\),\(\$4,\$5,\$6\)`).
				WithArgs("Book Two", 1, "A good book", "Book Three", 1, "A bad book").
				WillReturnResult(sqlmock.NewResult(0, 2))

			// when
			err := repository.Create(books)

			// then
			assert.NoError(t, err)
		})
	})

	t.Run("FindById", func(t *testing.T) {
		t.Run("should return error if executing query failed", func(t *testing.T) {
			// given
			var id uint64 = 1

			dbmock.
				ExpectQuery(`select id, name, authorId, description from books where id = \$1`).
				WillReturnError(errors.New("database error"))

			// when
			users, err := repository.FindById(id)

			// then
			assert.Error(t, err)
			assert.Nil(t, users)
		})

		t.Run("should return books by id", func(t *testing.T) {
			// given
			var id uint64 = 1

			dbmock.
				ExpectQuery(`select id, name, authorId, description from books where id = \$1 limit 1`).
				WithArgs(1).
				WillReturnRows(sqlmock.NewRows([]string{"id", "name", "authorId", "description"}).
					AddRow(1, "Updated Book", 1, "An updated description"))

			// when
			book, err := repository.FindById(id)

			// then
			assert.NoError(t, err)
			assert.Equal(t, book.ID, uint64(1))
		})
	})

	t.Run("Delete", func(t *testing.T) {
		t.Run("should return error if executing query failed", func(t *testing.T) {
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

			dbmock.
				ExpectExec(`delete from books`).
				WillReturnError(errors.New("database error"))

			// when
			err := repository.Delete(books)

			// then
			assert.Error(t, err)
		})

		t.Run("should delete books in batches", func(t *testing.T) {
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

			dbmock.
				ExpectExec(`delete from books where id in \(\$1,\$2\)`).
				WithArgs(1, 2).
				WillReturnResult(sqlmock.NewResult(0, 2))

			// when
			err := repository.Delete(books)

			// then
			assert.NoError(t, err)
		})
	})
	t.Run("Update", func(t *testing.T) {
		t.Run("should return error if executing query failed", func(t *testing.T) {
			// given
			newBookData := &model.UpdateBook{
				Name:        "Updated Book",
				Description: "An updated description",
			}

			dbmock.
				ExpectExec(`update books set name = \$1, description = \$2 where id = \$3`).
				WillReturnError(errors.New("database error"))

			// when
			err := repository.Update(1, newBookData)

			// then
			assert.Error(t, err)
		})

		t.Run("should update book", func(t *testing.T) {
			// given
			newBookData := &model.UpdateBook{
				Name:        "Updated Book",
				Description: "An updated description",
			}

			dbmock.
				ExpectExec(`update books set name = \$1, description = \$2 where id = \$3`).
				WithArgs("Updated Book", "An updated description", 2).
				WillReturnResult(sqlmock.NewResult(0, 1))

			// when
			err := repository.Update(2, newBookData)

			// then
			assert.NoError(t, err)
		})
	})
}
