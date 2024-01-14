package user_repository

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/akatranlp/hsfl-master-ai-cloud-engineering/lib/database"
	"github.com/akatranlp/hsfl-master-ai-cloud-engineering/user-service/user/model"

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

const createUsersTable = `
create table if not exists users (
    id				serial primary key, 
	email			varchar(100) not null unique,
	password 		bytea not null,
	profile_name 	varchar(100) not null,
	balance 		int not null default 0,
	token_version 	bigint not null default 0
)
`

func (repo *PsqlRepository) Migrate() error {
	_, err := repo.db.Exec(createUsersTable)
	return err
}

const createUsersBatchQuery = `
insert into users (email, password, profile_name) values %s
`

func (repo *PsqlRepository) Create(users []*model.DbUser) error {
	placeholders := make([]string, len(users))
	values := make([]interface{}, len(users)*3)

	for i := 0; i < len(users); i++ {
		placeholders[i] = fmt.Sprintf("($%d,$%d,$%d)", i*3+1, i*3+2, i*3+3)
		values[i*3+0] = users[i].Email
		values[i*3+1] = users[i].Password
		values[i*3+2] = users[i].ProfileName
	}

	query := fmt.Sprintf(createUsersBatchQuery, strings.Join(placeholders, ","))
	_, err := repo.db.Exec(query, values...)
	return err
}

const updateUserQuery = `
update users set profile_name = $1, password = $2, balance = $3, token_version = $4 where id = $5 returning id
`

func (repo *PsqlRepository) Update(id uint64, user *model.DbUserPatch) error {
	dbUser, err := repo.FindById(id)
	if err != nil {
		return err
	}
	if user.ProfileName != nil {
		dbUser.ProfileName = *user.ProfileName
	}
	if user.Password != nil {
		dbUser.Password = *user.Password
	}
	if user.Balance != nil {
		dbUser.Balance = *user.Balance
	}
	if user.TokenVersion != nil {
		dbUser.TokenVersion = *user.TokenVersion
	}

	_, err = repo.db.Exec(updateUserQuery, dbUser.ProfileName, dbUser.Password, dbUser.Balance, dbUser.TokenVersion, dbUser.ID)
	return err
}

const findAllUsersQuery = `
select id, email, password, profile_name, balance, token_version from users
`

func (repo *PsqlRepository) FindAll() ([]*model.DbUser, error) {
	rows, err := repo.db.Query(findAllUsersQuery)
	if err != nil {
		return nil, err
	}

	users := make([]*model.DbUser, 0)
	for rows.Next() {
		user := model.DbUser{}
		if err := rows.Scan(&user.ID, &user.Email, &user.Password, &user.ProfileName, &user.Balance, &user.TokenVersion); err != nil {
			return nil, err
		}

		users = append(users, &user)
	}
	return users, nil
}

const findUsersByEmailQuery = `
select id, email, password, profile_name, balance, token_version from users where email = $1
`

func (repo *PsqlRepository) FindByEmail(email string) ([]*model.DbUser, error) {
	rows, err := repo.db.Query(findUsersByEmailQuery, email)
	if err != nil {
		return nil, err
	}

	var users []*model.DbUser
	for rows.Next() {
		user := model.DbUser{}
		if err := rows.Scan(&user.ID, &user.Email, &user.Password, &user.ProfileName, &user.Balance, &user.TokenVersion); err != nil {
			return nil, err
		}

		users = append(users, &user)
	}

	return users, nil
}

const findUsersByIdQuery = `
select id, email, password, profile_name, balance, token_version from users where id = $1 LIMIT 1
`

func (repo *PsqlRepository) FindById(id uint64) (*model.DbUser, error) {
	row := repo.db.QueryRow(findUsersByIdQuery, id)
	user := model.DbUser{}
	if err := row.Scan(&user.ID, &user.Email, &user.Password, &user.ProfileName, &user.Balance, &user.TokenVersion); err != nil {
		return nil, err
	}
	return &user, nil
}

const deleteUsersBatchQuery = `
delete from users where id in (%s)
`

func (repo *PsqlRepository) Delete(users []*model.DbUser) error {
	placeholders := make([]string, len(users))
	ids := make([]interface{}, len(users))

	for i := 0; i < len(users); i++ {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
		ids[i] = users[i].ID
	}

	query := fmt.Sprintf(deleteUsersBatchQuery, strings.Join(placeholders, ","))
	_, err := repo.db.Exec(query, ids...)
	return err
}
