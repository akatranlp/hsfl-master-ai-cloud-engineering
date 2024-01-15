package containerhelpers

import (
	"context"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func StartBookService(postgresHost string, postgresPort string, authIsEnabled string, grpcEnabled string) (testcontainers.Container, error) {
	req := testcontainers.ContainerRequest{
		Image:        "akatranlp/book-service:latest",
		ExposedPorts: []string{"8080/tcp", "8081/tcp"},
		Env: map[string]string{
			"PORT":                         "8080",
			"GRPC_PORT":                    "8081",
			"GRPC_COMMUNICATION":           grpcEnabled,
			"AUTH_SERVICE_ENDPOINT":        "http://user:8081",
			"TRANSACTION_SERVICE_ENDPOINT": "http://transaction:8081",
			"POSTGRES_HOST":                postgresHost,
			"POSTGRES_PORT":                postgresPort,
			"POSTGRES_USER":                "postgres",
			"POSTGRES_PASSWORD":            "postgres",
			"POSTGRES_DB":                  "postgres",
			"AUTH_IS_ACTIVE":               authIsEnabled,
		},
		WaitingFor: wait.ForAll(
			wait.ForListeningPort("8081/tcp"),
			wait.ForListeningPort("8080/tcp"),
			wait.ForLog(".*GRPC-Server started on Port:.*").
				AsRegexp().WithStartupTimeout(60+time.Second),
			wait.ForLog(".*REST-Server started on Port:.*").
				AsRegexp().WithStartupTimeout(60+time.Second),
		).WithStartupTimeoutDefault(60 * time.Second),
	}

	return testcontainers.GenericContainer(context.Background(), testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
}

func GetBookServiceData(container testcontainers.Container, err error) (
	bookService testcontainers.Container,
	bookServiceRESTPort int,
	bookServiceGRPCPort int,
	bookServiceHost string,
	Error error,
) {
	bookService = container
	bookServiceRESTPort = -1
	bookServiceGRPCPort = -1
	bookServiceHost = ""
	Error = err

	if err != nil {
		return
	}

	bookServiceREST, err := bookService.MappedPort(context.Background(), "8080")
	if err != nil {
		Error = err
		return
	}
	bookServiceRESTPort = bookServiceREST.Int()

	bookServiceGRPC, err := bookService.MappedPort(context.Background(), "8081")
	if err != nil {
		Error = err
		return
	}
	bookServiceGRPCPort = bookServiceGRPC.Int()

	bookServiceHost, err = bookService.Host(context.Background())
	if err != nil {
		Error = err
		return
	}

	return
}
