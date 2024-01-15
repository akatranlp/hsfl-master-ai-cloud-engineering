package repository

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/akatranlp/hsfl-master-ai-cloud-engineering/user-service/model"
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
			users := []*model.DbUser{{
				ID:          1,
				Email:       "doesnt matter",
				Password:    []byte("doesnt matter"),
				ProfileName: "doesnt matter",
				Balance:     0,
			}}

			dbmock.
				ExpectExec(`insert into users`).
				WillReturnError(errors.New("database error"))

			// when
			err := repository.Create(users)

			// then
			assert.Error(t, err)
		})

		t.Run("should insert users in batches", func(t *testing.T) {
			// given
			users := []*model.DbUser{
				{
					ID:          1,
					Email:       "test@test.com",
					Password:    []byte("test"),
					ProfileName: "Toni Tester",
					Balance:     0,
				},
				{
					ID:          2,
					Email:       "abc@abc.com",
					Password:    []byte("abc"),
					ProfileName: "ABC ABC",
					Balance:     0,
				},
			}

			dbmock.
				ExpectExec(`insert into users \(email, password, profile_name\) values \(\$1,\$2,\$3\),\(\$4,\$5,\$6\)`).
				WithArgs("test@test.com", []byte("test"), "Toni Tester", "abc@abc.com", []byte("abc"), "ABC ABC").
				WillReturnResult(sqlmock.NewResult(0, 2))

			// when
			err := repository.Create(users)

			// then
			assert.NoError(t, err)
			assert.NoError(t, dbmock.ExpectationsWereMet())
		})
	})

	t.Run("Update", func(t *testing.T) {
		t.Run("should return error if executing query failed", func(t *testing.T) {
			// given
			name := "doesn't matter"
			balance := int64(1000)
			user := &model.DbUserPatch{
				ProfileName: &name,
				Balance:     &balance,
			}

			dbmock.
				ExpectQuery(`select id, email, password, profile_name, balance, token_version from users where id = \$1 LIMIT 1`).
				WithArgs(1).
				WillReturnError(errors.New("database error"))

			// when
			err := repository.Update(1, user)

			// then
			assert.Error(t, err)
			assert.NoError(t, dbmock.ExpectationsWereMet())
		})

		t.Run("should return error if executing query failed", func(t *testing.T) {
			// given
			name := "doesn't matter"
			balance := int64(1000)
			password := []byte("hash2")
			tokenVersion := uint64(1)
			user := &model.DbUserPatch{
				ProfileName:  &name,
				Balance:      &balance,
				Password:     &password,
				TokenVersion: &tokenVersion,
			}

			dbmock.
				ExpectQuery(`select id, email, password, profile_name, balance, token_version from users where id = \$1 LIMIT 1`).
				WithArgs(1).
				WillReturnRows(sqlmock.NewRows([]string{"id", "email", "password", "profile_name", "balance", "token_version"}).
					AddRow(1, "test@test.com", []byte("hash"), "Toni Tester", 0, 0))

			dbmock.
				ExpectExec("").
				WillReturnError(errors.New("database error"))

			// when
			err := repository.Update(1, user)

			// then
			assert.Error(t, err)
			assert.NoError(t, dbmock.ExpectationsWereMet())
		})

		t.Run("should update user", func(t *testing.T) {
			// given
			name := "doesn't matter"
			balance := int64(1000)
			password := []byte("hash2")
			tokenVersion := uint64(1)
			user := &model.DbUserPatch{
				ProfileName:  &name,
				Balance:      &balance,
				Password:     &password,
				TokenVersion: &tokenVersion,
			}

			dbmock.
				ExpectQuery(`select id, email, password, profile_name, balance, token_version from users where id = \$1 LIMIT 1`).
				WithArgs(1).
				WillReturnRows(sqlmock.NewRows([]string{"id", "email", "password", "profile_name", "balance", "token_version"}).
					AddRow(1, "test@test.com", []byte("hash"), "Toni Tester", 0, 0))

			dbmock.
				ExpectExec(`update users set profile_name = \$1, password = \$2, balance = \$3, token_version = \$4 where id = \$5 returning id`).
				WithArgs("doesn't matter", []byte("hash2"), 1000, 1, 1).
				WillReturnResult(sqlmock.NewResult(1, 1))

			// when
			err := repository.Update(1, user)

			// then
			assert.NoError(t, err)
			assert.NoError(t, dbmock.ExpectationsWereMet())
		})
	})

	t.Run("FindAll", func(t *testing.T) {
		t.Run("should return error if executing query failed", func(t *testing.T) {
			// given
			dbmock.ExpectQuery(`select id, email, password, profile_name, balance, token_version from users`).
				WillReturnError(errors.New("database error"))

			// when
			users, err := repository.FindAll()

			// then
			assert.Error(t, err)
			assert.Nil(t, users)
		})
		t.Run("should return all users", func(t *testing.T) {
			// given
			dbmock.ExpectQuery(`select id, email, password, profile_name, balance, token_version from users`).
				WillReturnRows(sqlmock.NewRows([]string{"id", "email", "password", "profile_name", "balance", "token_version"}).
					AddRow(1, "test@test.com", []byte("hash"), "Toni Tester", 0, 0).
					AddRow(2, "abc@abc.com", []byte("hash"), "ABC ABC", 0, 0))

			// when
			users, err := repository.FindAll()

			// then
			assert.NoError(t, err)
			assert.NoError(t, dbmock.ExpectationsWereMet())
			assert.Len(t, users, 2)
			assert.Equal(t, "Toni Tester", users[0].ProfileName)
			assert.Equal(t, "ABC ABC", users[1].ProfileName)
		})
	})

	t.Run("FindByEmail", func(t *testing.T) {
		t.Run("should return error if executing query failed", func(t *testing.T) {
			// given
			email := "test@test.com"

			dbmock.
				ExpectQuery(`select id, email, password, profile_name, balance, token_version from users where email = \$1`).
				WillReturnError(errors.New("database error"))

			// when
			users, err := repository.FindByEmail(email)

			// then
			assert.Error(t, err)
			assert.Nil(t, users)
		})

		t.Run("should return users by email", func(t *testing.T) {
			// given
			email := "test@test.com"

			dbmock.
				ExpectQuery(`select id, email, password, profile_name, balance, token_version from users where email = \$1`).
				WillReturnRows(sqlmock.NewRows([]string{"id", "email", "password", "profile_name", "balance", "token_version"}).
					AddRow(1, "test@test.com", []byte("hash"), "Toni Tester", 0, 0))

			// when
			users, err := repository.FindByEmail(email)

			// then
			assert.NoError(t, err)
			assert.Len(t, users, 1)
		})
	})

	t.Run("FindById", func(t *testing.T) {
		t.Run("should return error if executing query failed", func(t *testing.T) {
			// given
			id := uint64(1)

			dbmock.
				ExpectQuery(`select id, email, password, profile_name, balance, token_version from users where id = \$1`).
				WillReturnError(errors.New("database error"))

			// when
			users, err := repository.FindById(id)

			// then
			assert.Error(t, err)
			assert.Nil(t, users)
		})

		t.Run("should return users by id", func(t *testing.T) {
			// given
			id := uint64(1)

			dbmock.
				ExpectQuery(`select id, email, password, profile_name, balance, token_version from users where id = \$1`).
				WillReturnRows(sqlmock.NewRows([]string{"id", "email", "password", "profile_name", "balance", "token_version"}).
					AddRow(1, "test@test.com", []byte("hash"), "Toni Tester", 0, 0))

			// when
			user, err := repository.FindById(id)

			// then
			assert.NoError(t, err)
			assert.Equal(t, "Toni Tester", user.ProfileName)
		})
	})

	t.Run("Delete", func(t *testing.T) {
		t.Run("should return error if executing query failed", func(t *testing.T) {
			// given
			users := []*model.DbUser{
				{
					ID:          1,
					Email:       "test@test.com",
					Password:    []byte("hash"),
					ProfileName: "Toni Tester",
					Balance:     0,
				},
				{
					ID:          2,
					Email:       "abc@abc.com",
					Password:    []byte("hash"),
					ProfileName: "ABC ABC",
					Balance:     0,
				},
			}

			dbmock.
				ExpectExec(`delete from users`).
				WillReturnError(errors.New("database error"))

			// when
			err := repository.Delete(users)

			// then
			assert.Error(t, err)
		})

		t.Run("should delete users in batches", func(t *testing.T) {
			// given
			users := []*model.DbUser{
				{
					ID:          1,
					Email:       "test@test.com",
					Password:    []byte("hash"),
					ProfileName: "Toni Tester",
					Balance:     0,
				},
				{
					ID:          2,
					Email:       "abc@abc.com",
					Password:    []byte("hash"),
					ProfileName: "ABC ABC",
					Balance:     0,
				},
			}

			dbmock.
				ExpectExec(`delete from users where id in \(\$1,\$2\)`).
				WithArgs(1, 2).
				WillReturnResult(sqlmock.NewResult(0, 2))

			// when
			err := repository.Delete(users)

			// then
			assert.NoError(t, err)
		})
	})
}
