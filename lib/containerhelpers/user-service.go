package containerhelpers

import (
	"context"
	"time"

	"github.com/akatranlp/hsfl-master-ai-cloud-engineering/lib/utils"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func StartUserService(postgresHost string, postgresPort string) (testcontainers.Container, error) {
	privateKey, publicKey := utils.GenerateRSAKeyPairPem()

	req := testcontainers.ContainerRequest{
		Image:        "akatranlp/user-service:latest",
		ExposedPorts: []string{"8080/tcp", "8081/tcp"},
		Env: map[string]string{
			"PORT":               "8080",
			"GRPC_PORT":          "8081",
			"GRPC_COMMUNICATION": "true",
			"POSTGRES_HOST":      postgresHost,
			"POSTGRES_PORT":      postgresPort,
			"POSTGRES_USER":      "postgres",
			"POSTGRES_PASSWORD":  "postgres",
			"POSTGRES_DB":        "postgres",
			"AUTH_IS_ACTIVE":     "true",
			"JWT_PRIVATE_KEY":    privateKey,
			"JWT_PUBLIC_KEY":     publicKey,
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
