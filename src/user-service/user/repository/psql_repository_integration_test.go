package user_repository

import (
	"context"
	"database/sql"
	"testing"

	"github.com/akatranlp/hsfl-master-ai-cloud-engineering/lib/containerhelpers"
	"github.com/akatranlp/hsfl-master-ai-cloud-engineering/lib/database"
	"github.com/akatranlp/hsfl-master-ai-cloud-engineering/user-service/user/model"
	"github.com/stretchr/testify/assert"
)

func TestIntegrationPsqlRepository(t *testing.T) {
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
		t.Fatalf("could not create user repository: %s", err.Error())
	}
	t.Cleanup(clearTables(t, repository.db))

	t.Run("Migrate", func(t *testing.T) {
		t.Run("should create users table", func(t *testing.T) {
			t.Cleanup(clearTables(t, repository.db))

			// given
			// when
			err := repository.Migrate()

			// then
			assert.NoError(t, err)
			assertTableExists(t, repository.db, "users", []string{"email", "password"})
		})
	})

	t.Run("Create", func(t *testing.T) {
		t.Run("should insert users in batches", func(t *testing.T) {
			t.Cleanup(clearTables(t, repository.db))

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

			// when
			err := repository.Create(users)

			// then
			assert.NoError(t, err)
			assert.Equal(t, users[0], getUserFromDatabase(t, repository.db, "test@test.com"))
			assert.Equal(t, users[1], getUserFromDatabase(t, repository.db, "abc@abc.com"))
		})
	})

	t.Run("Update", func(t *testing.T) {
		t.Run("should update user", func(t *testing.T) {
			t.Cleanup(clearTables(t, repository.db))

			// given
			insertUser(t, repository.db, &model.DbUser{
				Email:       "test@test.com",
				Password:    []byte("hash"),
				ProfileName: "Toni Tester",
				Balance:     0,
			})
			newName := "Tino Taster"
			newBalance := int64(1000)

			user := &model.DbUserPatch{
				ProfileName: &newName,
				Balance:     &newBalance,
			}

			// when
			err := repository.Update(3, user)

			// then
			assert.NoError(t, err)
			assert.Equal(t, &model.DbUser{
				ID:          3,
				Email:       "test@test.com",
				Password:    []byte("hash"),
				ProfileName: "Tino Taster",
				Balance:     1000,
			}, getUserFromDatabase(t, repository.db, "test@test.com"))
		})
	})

	t.Run("FindAll", func(t *testing.T) {
		t.Run("should return all products", func(t *testing.T) {
			t.Cleanup(clearTables(t, repository.db))

			// given
			users := []*model.DbUser{
				{
					Email:       "test@test.com",
					Password:    []byte("hash"),
					ProfileName: "Toni Tester",
					Balance:     0,
				},
				{
					Email:       "abc@abc.com",
					Password:    []byte("hash"),
					ProfileName: "ABC ABC",
					Balance:     0,
				},
			}

			for _, user := range users {
				insertUser(t, repository.db, user)
			}

			// when
			users, err := repository.FindAll()

			// then
			assert.NoError(t, err)
			assert.Len(t, users, 2)
		})
	})

	t.Run("FindByEmail", func(t *testing.T) {
		t.Run("should return user", func(t *testing.T) {
			t.Cleanup(clearTables(t, repository.db))

			// given
			insertUser(t, repository.db, &model.DbUser{
				Email:       "test@test.com",
				Password:    []byte("hash"),
				ProfileName: "Toni Tester",
				Balance:     0,
			})

			// when
			user, err := repository.FindByEmail("test@test.com")

			// then
			assert.NoError(t, err)
			assert.NotNil(t, user)
		})
	})

	t.Run("FindById", func(t *testing.T) {
		t.Run("should return user", func(t *testing.T) {
			t.Cleanup(clearTables(t, repository.db))

			// given
			insertUser(t, repository.db, &model.DbUser{
				ID:          7,
				Email:       "test@test.com",
				Password:    []byte("hash"),
				ProfileName: "Toni Tester",
				Balance:     0,
			})

			// when
			user, err := repository.FindById(7)

			// then
			assert.NoError(t, err)
			assert.NotNil(t, user)
		})
	})

	t.Run("Delete", func(t *testing.T) {
		t.Run("should delete provided users", func(t *testing.T) {
			t.Cleanup(clearTables(t, repository.db))

			// given
			users := []*model.DbUser{
				{
					ID:          8,
					Email:       "test@test.com",
					Password:    []byte("hash"),
					ProfileName: "Toni Tester",
					Balance:     0,
				},
				{
					ID:          9,
					Email:       "abc@abc.com",
					Password:    []byte("hash"),
					ProfileName: "ABC ABC",
					Balance:     0,
				},
			}

			for _, user := range users {
				insertUser(t, repository.db, user)
			}

			// when
			err := repository.Delete([]*model.DbUser{users[1]})

			// then
			assert.NoError(t, err)
			assert.Equal(t, users[0], getUserFromDatabase(t, repository.db, "test@test.com"))
			assert.Nil(t, getUserFromDatabase(t, repository.db, "abc@abc.com"))
		})
	})
}

func getUserFromDatabase(t *testing.T, db *sql.DB, email string) *model.DbUser {
	row := db.QueryRow(`select id, email, password, profile_name, balance from users where email = $1`, email)

	var user model.DbUser
	if err := row.Scan(&user.ID, &user.Email, &user.Password, &user.ProfileName, &user.Balance); err != nil {
		return nil
	}

	return &user
}

func insertUser(t *testing.T, db *sql.DB, user *model.DbUser) {
	_, err := db.Exec(`insert into users (email, password, profile_name) values ($1, $2, $3)`, user.Email, user.Password, user.ProfileName)
	if err != nil {
		t.Logf("could not insert user: %s", err.Error())
		t.FailNow()
	}
}

func clearTables(t *testing.T, db *sql.DB) func() {
	return func() {
		if _, err := db.Exec("delete from users"); err != nil {
			t.Logf("could not delete rows from users: %s", err.Error())
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
