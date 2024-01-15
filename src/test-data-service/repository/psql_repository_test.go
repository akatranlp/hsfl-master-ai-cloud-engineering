package repository

import (
	"errors"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	crypto_mocks "github.com/akatranlp/hsfl-master-ai-cloud-engineering/lib/crypto/_mocks"
	mocks "github.com/akatranlp/hsfl-master-ai-cloud-engineering/test-data-service/_mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestPsqlRepository(t *testing.T) {
	ctrl := gomock.NewController(t)

	db, dbmock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}

	hasher := crypto_mocks.NewMockHasher(ctrl)

	testConfig := mocks.NewMockConfig(ctrl)

	repository := PsqlRepository{db, testConfig, hasher}

	t.Run("ResetDatabase", func(t *testing.T) {

		t.Run("should error if config errors", func(t *testing.T) {
			// given
			testConfig.
				EXPECT().
				GetSqlString().
				Return("", errors.New("error"))

			// when
			err := repository.ResetDatabase()

			// then
			assert.Error(t, err)
		})

		t.Run("should error if hasher errors", func(t *testing.T) {
			// given
			password := "password"

			testConfig.
				EXPECT().
				GetSqlString().
				Return("", nil)

			testConfig.
				EXPECT().
				GetUserPassword().
				Return(password)

			hasher.
				EXPECT().
				Hash([]byte(password)).
				Return(nil, errors.New("error"))

			// when
			err := repository.ResetDatabase()

			// then
			assert.Error(t, err)
		})

		t.Run("should error if database error accured", func(t *testing.T) {
			// given
			password := "password"
			hashedPassword := "hashedPassword"

			testConfig.
				EXPECT().
				GetSqlString().
				Return("", nil)

			testConfig.
				EXPECT().
				GetUserPassword().
				Return(password)

			hasher.
				EXPECT().
				Hash([]byte(password)).
				Return([]byte(hashedPassword), nil)

			dbmock.
				ExpectExec("").
				WillReturnError(errors.New("error"))

			// when
			err := repository.ResetDatabase()

			// then
			assert.Error(t, err)
		})

		t.Run("should exec without password if string is wrong", func(t *testing.T) {
			// given
			password := "password"
			hashedPassword := "hashedPassword"

			sql := "DROP TABLE IF EXISTS users;"

			testConfig.
				EXPECT().
				GetSqlString().
				Return(sql, nil)

			testConfig.
				EXPECT().
				GetUserPassword().
				Return(password)

			hasher.
				EXPECT().
				Hash([]byte(password)).
				Return([]byte(hashedPassword), nil)

			command := fmt.Sprintf("BEGIN;\n%s\n%s\nCOMMIT;", resetDataBaseQuery, sql)

			dbmock.
				ExpectExec(command).
				WillReturnResult(sqlmock.NewResult(0, 0))

			// when
			err := repository.ResetDatabase()

			// then
			assert.NoError(t, err)
		})

		t.Run("should exec with password if string is correct", func(t *testing.T) {
			// given
			password := "password"
			hashedPassword := "hashedPassword"

			sql := "INSERT INTO users (username, password) VALUES ('test', $1);"

			testConfig.
				EXPECT().
				GetSqlString().
				Return(sql, nil)

			testConfig.
				EXPECT().
				GetUserPassword().
				Return(password)

			hasher.
				EXPECT().
				Hash([]byte(password)).
				Return([]byte(hashedPassword), nil)

			newSql := `INSERT INTO users \(username, password\) VALUES \('test', '` + hashedPassword + `'\);`
			command := fmt.Sprintf("BEGIN;\n%s\n%s\nCOMMIT;", resetDataBaseQuery, newSql)

			dbmock.
				ExpectExec(command).
				WillReturnResult(sqlmock.NewResult(0, 0))

			// when
			err := repository.ResetDatabase()

			// then
			assert.NoError(t, err)
		})
	})
}
