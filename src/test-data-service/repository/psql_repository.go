package repository

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/akatranlp/hsfl-master-ai-cloud-engineering/lib/crypto"
	"github.com/akatranlp/hsfl-master-ai-cloud-engineering/lib/database"
	"github.com/akatranlp/hsfl-master-ai-cloud-engineering/test-data-service/config"
	_ "github.com/lib/pq"
)

type PsqlRepository struct {
	db             *sql.DB
	testDataConfig config.Config
	hasher         crypto.Hasher
}

func NewPsqlRepository(
	config database.Config,
	testDataConfig config.Config,
	hasher crypto.Hasher,
) (*PsqlRepository, error) {
	dsn := config.Dsn()
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	return &PsqlRepository{db, testDataConfig, hasher}, nil
}

const resetDataBaseQuery = `
drop table if exists transactions;
drop table if exists chapters;
drop table if exists books;
drop table if exists users;
`

func (r *PsqlRepository) ResetDatabase() error {
	data, err := r.testDataConfig.GetSqlString()
	if err != nil {
		return err
	}

	command := fmt.Sprintf("BEGIN;\n%s\n%s\nCOMMIT;", resetDataBaseQuery, data)

	hashedPassword, err := r.hasher.Hash([]byte(r.testDataConfig.GetUserPassword()))
	if err != nil {
		return err
	}

	command = strings.ReplaceAll(command, "$1", fmt.Sprintf("'%s'", hashedPassword))

	_, err = r.db.Exec(command)
	if err != nil {
		return err
	}

	return nil
}
