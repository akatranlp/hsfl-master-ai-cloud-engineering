package containerhelpers

import (
	"context"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func StartPostgres(withTestData bool) (testcontainers.Container, error) {
	var image string
	if withTestData {
		image = "akatranlp/postgres:latest"
	} else {
		image = "postgres:latest"
	}
	req := testcontainers.ContainerRequest{
		Image:        image,
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_USER":     "postgres",
			"POSTGRES_PASSWORD": "postgres",
			"POSTGRES_DB":       "postgres",
		},
		WaitingFor: wait.ForAll(wait.ForListeningPort("5432/tcp"), wait.ForLog(".*database system is ready to accept connections.*").AsRegexp().WithStartupTimeout(60+time.Second)).WithStartupTimeoutDefault(60 * time.Second),
	}

	return testcontainers.GenericContainer(context.Background(), testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
}
