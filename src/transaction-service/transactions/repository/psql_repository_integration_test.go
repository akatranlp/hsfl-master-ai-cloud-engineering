package transactions_repository

import (
	"context"
	"database/sql"
	"testing"

	"github.com/akatranlp/hsfl-master-ai-cloud-engineering/lib/containerhelpers"
	"github.com/akatranlp/hsfl-master-ai-cloud-engineering/lib/database"
	"github.com/akatranlp/hsfl-master-ai-cloud-engineering/transaction-service/transactions/model"
	"github.com/stretchr/testify/assert"
)

func TestIntegrationPsqlRepository(t *testing.T) {
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

	repository, err := NewPsqlRepository(database.PsqlConfig{
		Host:     "0.0.0.0",
		Port:     port.Int(),
		Username: "postgres",
		Password: "postgres",
		Database: "postgres",
	})
	if err != nil {
		t.Fatalf("could not create transaction repository: %s", err.Error())
	}
	t.Cleanup(clearTables(t, repository.db))

	t.Run("Migrate", func(t *testing.T) {
		t.Run("should create transactions table", func(t *testing.T) {
			t.Cleanup(clearTables(t, repository.db))

			// given
			createUserBookAndChapterTable(t, repository.db)
			// when
			err := repository.Migrate()

			// then
			assert.NoError(t, err)
			assertTableExists(t, repository.db, "transactions", []string{"id"})
		})
	})

	t.Run("Create", func(t *testing.T) {
		t.Run("should insert transactions in batches", func(t *testing.T) {
			t.Cleanup(clearTables(t, repository.db))

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
			// when
			err := repository.Create(transactions)

			// then
			assert.NoError(t, err)
			assert.Equal(t, transactions[0], getTransactionFromDatabase(t, repository.db, 1))
			assert.Equal(t, transactions[1], getTransactionFromDatabase(t, repository.db, 2))
		})
	})

	t.Run("FindAll", func(t *testing.T) {
		t.Run("should return all products", func(t *testing.T) {
			t.Cleanup(clearTables(t, repository.db))

			// given
			transactions := []*model.Transaction{
				{
					ID:              3,
					ChapterID:       1,
					BookID:          1,
					ReceivingUserID: 2,
					PayingUserID:    1,
					Amount:          0,
				},
				{
					ID:              4,
					ChapterID:       2,
					BookID:          1,
					ReceivingUserID: 2,
					PayingUserID:    1,
					Amount:          100,
				},
			}

			for _, transaction := range transactions {
				insertTransaction(t, repository.db, transaction)
			}

			// when
			transactions, err := repository.FindAll()

			// then
			assert.NoError(t, err)
			assert.Len(t, transactions, 2)
		})
	})

	t.Run("FindAllForUserId", func(t *testing.T) {
		t.Run("should return transaction", func(t *testing.T) {
			t.Cleanup(clearTables(t, repository.db))

			// given
			transactions := []*model.Transaction{
				{
					ID:              5,
					ChapterID:       1,
					BookID:          1,
					ReceivingUserID: 2,
					PayingUserID:    1,
					Amount:          0,
				},
				{
					ID:              6,
					ChapterID:       2,
					BookID:          1,
					ReceivingUserID: 2,
					PayingUserID:    1,
					Amount:          100,
				},
			}

			for _, transaction := range transactions {
				insertTransaction(t, repository.db, transaction)
			}

			// when
			transactions, err := repository.FindAllForUserId(1)

			// then
			assert.NoError(t, err)
			assert.Equal(t, transactions[1], getTransactionFromDatabase(t, repository.db, 6))
		})
	})

	t.Run("FindAllForReceivingUserId", func(t *testing.T) {
		t.Run("should return transaction", func(t *testing.T) {
			t.Cleanup(clearTables(t, repository.db))

			// given
			transactions := []*model.Transaction{
				{
					ID:              7,
					ChapterID:       1,
					BookID:          1,
					ReceivingUserID: 1,
					PayingUserID:    2,
					Amount:          0,
				},
				{
					ID:              8,
					ChapterID:       1,
					BookID:          2,
					ReceivingUserID: 2,
					PayingUserID:    1,
					Amount:          100,
				},
			}

			for _, transaction := range transactions {
				insertTransaction(t, repository.db, transaction)
			}

			// when
			transactions, err := repository.FindAllForReceivingUserId(1)

			// then
			assert.NoError(t, err)
			assert.Equal(t, transactions[0], getTransactionFromDatabase(t, repository.db, 7))
		})
	})

	t.Run("FindById", func(t *testing.T) {
		t.Run("should return transaction", func(t *testing.T) {
			t.Cleanup(clearTables(t, repository.db))

			// given
			insertTransaction(t, repository.db, &model.Transaction{
				ID:              9,
				ChapterID:       2,
				BookID:          1,
				ReceivingUserID: 2,
				PayingUserID:    1,
				Amount:          100,
			})

			// when
			transaction, err := repository.FindById(9)

			// then
			assert.NoError(t, err)
			assert.NotNil(t, transaction)
		})
	})
}

func createUserBookAndChapterTable(t *testing.T, db *sql.DB) {
	db.Exec(`
create table if not exists users
	(
		id           serial primary key,
		email        varchar(100) not null unique,
		username     varchar(16)  not null unique,
		password     bytea        not null,
		profile_name varchar(100) not null,
		balance      int          not null default 0
	)
`)

	db.Exec(`
create table if not exists books
	(
	id          serial primary key,
	name        varchar(100) not null,
	authorId    int          not null,
	description text         not null,
	foreign key (authorId) REFERENCES users (id)
	)
`)
	//TODO: ERROR HERE: pq: there is no unique constraint matching given keys for referenced table "books"
	db.Exec(`
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
`)

	db.Exec(`
insert into users (email, username, password, profile_name, balance)
	values ('test@test.com', 'test', '$2a$10$uRS0zZPBGERH5Z3o2SgJhekhtR1z4orHigzpGNKZNTDGf0DcAhlRa', 'Toni Tester', 1000),
	('max.mustermann@gmail.com', 'mustermax', '$2a$10$saLKixUXtrauUIeLBoZeNuW/kacwWfLkPZHZGmP04xWAdxMn5uwty',
	'Max Mustermann', 1000)
`)

	db.Exec(`
insert into books (name, authorId, description)
	values ('Book One', 1, 'A good book'),
	('Book Two', 2, 'A bad book'),
	('Book Three', 1, 'A mid book')
`)

	db.Exec(`
insert into chapters (id, bookId, name, price, content)
	values (1, 1, 'The beginning', 0, 'Lorem Ipsum'),
	(2, 1, 'The beginning 2: Electric Boogaloo', 100, 'Lorem Ipsum 2'),
	(3, 1, 'The beginning 3: My Enemy', 100, 'Lorem Ipsum 3'),
	(1, 2, 'A different book chapter 1', 0, 'LorIp 4'),
	(2, 2, 'What came after', 100, 'Lorem Ipsum 5')
`)
}

func getTransactionFromDatabase(t *testing.T, db *sql.DB, id int) *model.Transaction {
	row := db.QueryRow(`select id, bookid, chapterid, receivinguserid, payinguserid, amount from transactions where id = $1`, id)

	var transaction model.Transaction
	if err := row.Scan(&transaction.ID, &transaction.BookID, &transaction.ChapterID, &transaction.ReceivingUserID, &transaction.PayingUserID, &transaction.Amount); err != nil {
		return nil
	}

	return &transaction
}

func insertTransaction(t *testing.T, db *sql.DB, transaction *model.Transaction) {
	_, err := db.Exec(`insert into transactions (id, bookid, chapterid, receivinguserid, payinguserid, amount) values ($1, $2, $3, $4, $5, $6)`, transaction.ID, transaction.BookID, transaction.ChapterID, transaction.ReceivingUserID, transaction.PayingUserID, transaction.Amount)
	if err != nil {
		t.Logf("could not insert transaction: %s", err.Error())
		t.FailNow()
	}
}

func clearTables(t *testing.T, db *sql.DB) func() {
	return func() {
		if _, err := db.Exec("delete from transactions"); err != nil {
			t.Logf("could not delete rows from transactions: %s", err.Error())
			t.FailNow()
		}
	}
}

func assertTableExists(t *testing.T, db *sql.DB, name string, columns []string) {
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
