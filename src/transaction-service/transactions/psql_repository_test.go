package transactions

import (
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"testing"

	"github.com/akatranlp/hsfl-master-ai-cloud-engineering/transaction-service/transactions/model"
	"github.com/stretchr/testify/assert"
)

func TestPsqlRepository(t *testing.T) {
	db, dbmock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}

	repository := PsqlRepository{db}

	t.Run("Create", func(t *testing.T) {
		t.Run("should return error if executing query failed", func(t *testing.T) {
			// given
			transactions := []*model.Transaction{{
				ID:           1,
				ChapterID:    1,
				PayingUserID: 1,
				Amount:       0,
			}}

			dbmock.
				ExpectExec(`insert into transactions`).
				WillReturnError(errors.New("database error"))

			// when
			err := repository.Create(transactions)

			// then
			assert.Error(t, err)
		})

		t.Run("should insert transactions in batches", func(t *testing.T) {
			// given
			transactions := []*model.Transaction{
				{
					ID:              1,
					ChapterID:       1,
					BookID:          1,
					ReceivingUserID: 2,
					PayingUserID:    1,
					Amount:          0,
				},
				{
					ID:              2,
					ChapterID:       2,
					BookID:          1,
					ReceivingUserID: 2,
					PayingUserID:    1,
					Amount:          100,
				},
			}

			dbmock.
				ExpectExec(`insert into transactions \(bookid, chapterid, receivinguserid, payinguserid, amount\) values \(\$1,\$2,\$3,\$4,\$5\),\(\$6,\$7,\$8,\$9,\$10\)`).
				WithArgs(1, 1, 2, 1, 0, 1, 2, 2, 1, 100).
				WillReturnResult(sqlmock.NewResult(0, 2))

			// when
			err := repository.Create(transactions)

			// then
			assert.NoError(t, err)
		})
	})

	t.Run("FindAllForUserId", func(t *testing.T) {
		t.Run("should return error if executing query failed", func(t *testing.T) {
			// given
			dbmock.ExpectQuery(`select id, bookid, chapterid, receivinguserid, payinguserid, amount from transactions where payinguserid = \$1`).
				WithArgs(2).
				WillReturnError(errors.New("database error"))

			// when
			transactions, err := repository.FindAllForUserId(uint64(2))

			// then
			assert.Error(t, err)
			assert.Nil(t, transactions)
			assert.NoError(t, dbmock.ExpectationsWereMet())
		})

		t.Run("should return all transactions", func(t *testing.T) {
			// given
			dbmock.ExpectQuery(`select id, bookid, chapterid, receivinguserid, payinguserid, amount from transactions where payinguserid = \$1`).
				WithArgs(2).
				WillReturnRows(sqlmock.NewRows([]string{"id", "bookid", "chapterid", "receivinguserid", "payinguserid", "amount"}).
					AddRow(1, 1, 1, 1, 2, 100).
					AddRow(2, 1, 2, 1, 2, 100))

			// when
			transactions, err := repository.FindAllForUserId(uint64(2))

			// then
			assert.NoError(t, err)
			assert.NoError(t, dbmock.ExpectationsWereMet())
			assert.Len(t, transactions, 2)
			assert.Equal(t, uint64(1), transactions[0].ID)
			assert.Equal(t, uint64(2), transactions[1].ID)
		})
	})

	t.Run("FindAllForReceivingUserId", func(t *testing.T) {
		t.Run("should return error if executing query failed", func(t *testing.T) {
			// given
			dbmock.ExpectQuery(`select id, bookid, chapterid, receivinguserid, payinguserid, amount from transactions where receivinguserid = \$1`).
				WithArgs(1).
				WillReturnError(errors.New("database error"))

			// when
			transactions, err := repository.FindAllForReceivingUserId(uint64(1))

			// then
			assert.Error(t, err)
			assert.Nil(t, transactions)
			assert.NoError(t, dbmock.ExpectationsWereMet())
		})

		t.Run("should return all transactions", func(t *testing.T) {
			// given
			dbmock.ExpectQuery(`select id, bookid, chapterid, receivinguserid, payinguserid, amount from transactions where receivinguserid = \$1`).
				WithArgs(1).
				WillReturnRows(sqlmock.NewRows([]string{"id", "bookid", "chapterid", "receivinguserid", "payinguserid", "amount"}).
					AddRow(1, 1, 1, 1, 2, 100).
					AddRow(2, 1, 2, 1, 2, 100))

			// when
			transactions, err := repository.FindAllForReceivingUserId(uint64(1))

			// then
			assert.NoError(t, err)
			assert.NoError(t, dbmock.ExpectationsWereMet())
			assert.Len(t, transactions, 2)
			assert.Equal(t, uint64(1), transactions[0].ID)
			assert.Equal(t, uint64(2), transactions[1].ID)
		})
	})
}
