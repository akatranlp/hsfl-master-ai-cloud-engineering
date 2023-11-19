package transactions

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/akatranlp/hsfl-master-ai-cloud-engineering/lib/database"
	"github.com/akatranlp/hsfl-master-ai-cloud-engineering/transaction-service/transactions/model"
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

const createTransactionsTable = `
create table if not exists transactions (
	id					serial primary key,
	bookid				int not null,
	chapterid    		int not null,
	receivinguserid		int not null,
	payinguserid 		int not null,
	amount 				int not null,
	foreign key (chapterid) references chapters(id),
	foreign key (bookid) references books (id),
	foreign key (payinguserid) references users(id),
	foreign key (receivinguserid) references users(id)
)
`

func (repo *PsqlRepository) Migrate() error {
	_, err := repo.db.Exec(createTransactionsTable)
	return err
}

const createTransactionsBatchQuery = `
insert into transactions (bookid, chapterid, receivinguserid, payinguserid, amount) values %s
`

func (repo *PsqlRepository) Create(transactions []*model.Transaction) error {
	placeholders := make([]string, len(transactions))
	values := make([]interface{}, len(transactions)*5)

	for i := 0; i < len(transactions); i++ {
		placeholders[i] = fmt.Sprintf("($%d,$%d,$%d,$%d,$%d)", i*5+1, i*5+2, i*5+3, i*5+4, i*5+5)
		values[i*5+0] = transactions[i].BookID
		values[i*5+1] = transactions[i].ChapterID
		values[i*5+2] = transactions[i].ReceivingUserID
		values[i*5+3] = transactions[i].PayingUserID
		values[i*5+4] = transactions[i].Amount
	}

	query := fmt.Sprintf(createTransactionsBatchQuery, strings.Join(placeholders, ","))
	_, err := repo.db.Exec(query, values...)
	return err
}

const findAllTransactionsQuery = `
select id, bookid, chapterid, receivinguserid, payinguserid, amount from transactions
`

func (repo *PsqlRepository) FindAll() ([]*model.Transaction, error) {
	rows, err := repo.db.Query(findAllTransactionsQuery)
	if err != nil {
		return nil, err
	}

	transactions := make([]*model.Transaction, 0)
	for rows.Next() {
		transaction := model.Transaction{}
		if err := rows.Scan(&transaction.ID, &transaction.BookID, &transaction.ChapterID, &transaction.ReceivingUserID, &transaction.PayingUserID, &transaction.Amount); err != nil {
			return nil, err
		}

		transactions = append(transactions, &transaction)
	}

	return transactions, nil
}

const findAllTransactionsForUserQuery = `
select id, bookid, chapterid, receivinguserid, payinguserid, amount from transactions where payinguserid = $1
`

func (repo *PsqlRepository) FindAllForUserId(userId uint64) ([]*model.Transaction, error) {
	rows, err := repo.db.Query(findAllTransactionsForUserQuery, userId)
	if err != nil {
		return nil, err
	}

	transactions := make([]*model.Transaction, 0)
	for rows.Next() {
		transaction := model.Transaction{}
		if err := rows.Scan(&transaction.ID, &transaction.BookID, &transaction.ChapterID, &transaction.ReceivingUserID, &transaction.PayingUserID, &transaction.Amount); err != nil {
			return nil, err
		}

		transactions = append(transactions, &transaction)
	}

	return transactions, nil
}

const findAllTransactionsForReceivingUserQuery = `
select id, bookid, chapterid, receivinguserid, payinguserid, amount from transactions where receivinguserid = $1
`

func (repo *PsqlRepository) FindAllForReceivingUserId(userId uint64) ([]*model.Transaction, error) {
	rows, err := repo.db.Query(findAllTransactionsForReceivingUserQuery, userId)
	if err != nil {
		return nil, err
	}

	transactions := make([]*model.Transaction, 0)
	for rows.Next() {
		transaction := model.Transaction{}
		if err := rows.Scan(&transaction.ID, &transaction.BookID, &transaction.ChapterID, &transaction.ReceivingUserID, &transaction.PayingUserID, &transaction.Amount); err != nil {
			return nil, err
		}

		transactions = append(transactions, &transaction)
	}

	return transactions, nil
}

const findTransactionbyIDQuery = `
select id, bookid, chapterid, receivinguserid, payinguserid, amount from transactions where id = $1
`

func (repo *PsqlRepository) FindById(id uint64) (*model.Transaction, error) {
	row := repo.db.QueryRow(findTransactionbyIDQuery, id)
	transaction := &model.Transaction{}
	if err := row.Scan(&transaction.ID, &transaction.BookID, &transaction.ChapterID, &transaction.ReceivingUserID, &transaction.PayingUserID, &transaction.Amount); err != nil {
		return nil, err
	}

	return transaction, nil
}

const findTransactionByUserIdAndChapterIdQuery = `
select id, bookid, chapterid, receivinguserid, payinguserid, amount from transactions where chapterid = $1 and payinguserid = $2
`

func (repo *PsqlRepository) FindForUserIdAndChapterId(userId uint64, chapterId uint64) (*model.Transaction, error) {
	row := repo.db.QueryRow(findTransactionByUserIdAndChapterIdQuery, chapterId, userId)
	transaction := &model.Transaction{}
	if err := row.Scan(&transaction.ID, &transaction.BookID, &transaction.ChapterID, &transaction.ReceivingUserID, &transaction.PayingUserID, &transaction.Amount); err != nil {
		return nil, err
	}

	return transaction, nil
}
