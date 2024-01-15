package repository

import (
	"context"
	"database/sql"
	"testing"

	"github.com/akatranlp/hsfl-master-ai-cloud-engineering/lib/containerhelpers"
	"github.com/akatranlp/hsfl-master-ai-cloud-engineering/lib/crypto"
	"github.com/akatranlp/hsfl-master-ai-cloud-engineering/lib/database"
	"github.com/akatranlp/hsfl-master-ai-cloud-engineering/test-data-service/config"
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

	password := "password"

	testDataConfig := config.TestDataConfig{
		UserPassword: password,
		File:         getSqlString(),
	}

	hasher := crypto.NewBcryptHasher()

	repository, err := NewPsqlRepository(database.PsqlConfig{
		Host:     "0.0.0.0",
		Port:     port.Int(),
		Username: "postgres",
		Password: "postgres",
		Database: "postgres",
	}, testDataConfig, hasher)

	if err != nil {
		t.Fatalf("could not create user repository: %s", err.Error())
	}

	t.Run("ResetDatabase", func(t *testing.T) {
		t.Run("should insert testdata at empty database", func(t *testing.T) {
			// given
			user := getUserFromDatabase(t, repository.db, "test@test.com")
			assert.Nil(t, user)

			// when
			err := repository.ResetDatabase()

			// then
			assert.NoError(t, err)
			user = getUserFromDatabase(t, repository.db, "test@test.com")
			assert.True(t, hasher.Validate([]byte(password), []byte(user.password)))
		})

		t.Run("should reset database", func(t *testing.T) {
			// given
			user := getUserFromDatabase(t, repository.db, "test@test.com")
			assert.NotNil(t, user)
			updateUser(t, repository.db, "test@test.com", "newPassword")
			user = getUserFromDatabase(t, repository.db, "test@test.com")
			assert.NotNil(t, user)
			assert.Equal(t, "newPassword", user.password)

			// when
			err := repository.ResetDatabase()

			// then
			assert.NoError(t, err)
			user = getUserFromDatabase(t, repository.db, "test@test.com")
			assert.True(t, hasher.Validate([]byte(password), []byte(user.password)))
		})
	})
}

func getSqlString() string {
	return `
	create table if not exists users (
		id				serial primary key, 
		email			varchar(100) not null unique,
		password 		bytea not null
	);

	INSERT INTO users (email, password) VALUES ('test@test.com', $1);
	`
}

type user struct {
	id       int
	email    string
	password string
}

func getUserFromDatabase(t *testing.T, db *sql.DB, email string) *user {
	row := db.QueryRow(`select id, email, password from users where email = $1`, email)

	var user user
	if err := row.Scan(&user.id, &user.email, &user.password); err != nil {
		return nil
	}

	return &user
}

func updateUser(t *testing.T, db *sql.DB, email string, password string) error {
	_, err := db.Exec(`update users set password = $1 where email = $2`, password, email)
	return err
}
