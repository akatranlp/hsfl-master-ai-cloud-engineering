package containerhelpers

import (
	"context"
	"time"

	"github.com/akatranlp/hsfl-master-ai-cloud-engineering/lib/utils"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func StartUserService(postgresHost string, postgresPort string, authIsEnabled string, grpcEnabled string) (testcontainers.Container, error) {
	accessPrivateKey, accessPublicKey := utils.GenerateRSAKeyPairPem()
	refreshPrivateKey, refreshPublicKey := utils.GenerateRSAKeyPairPem()

	req := testcontainers.ContainerRequest{
		Image:        "akatranlp/user-service:latest",
		ExposedPorts: []string{"8080/tcp", "8081/tcp"},
		Env: map[string]string{
			"PORT":                         "8080",
			"GRPC_PORT":                    "8081",
			"GRPC_COMMUNICATION":           grpcEnabled,
			"POSTGRES_HOST":                postgresHost,
			"POSTGRES_PORT":                postgresPort,
			"POSTGRES_USER":                "postgres",
			"POSTGRES_PASSWORD":            "postgres",
			"POSTGRES_DB":                  "postgres",
			"AUTH_IS_ACTIVE":               authIsEnabled,
			"JWT_ACCESS_PRIVATE_KEY":       accessPrivateKey,
			"JWT_ACCESS_PUBLIC_KEY":        accessPublicKey,
			"JWT_REFRESH_PRIVATE_KEY":      refreshPrivateKey,
			"JWT_REFRESH_PUBLIC_KEY":       refreshPublicKey,
			"JWT_ACCESS_TOKEN_EXPIRATION":  "15m",
			"JWT_REFRESH_TOKEN_EXPIRATION": "168h",
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

func GetUserServiceData(container testcontainers.Container, err error) (
	userService testcontainers.Container,
	userServiceRESTPort int,
	userServiceGRPCPort int,
	userServiceHost string,
	Error error,
) {
	userService = container
	userServiceRESTPort = -1
	userServiceGRPCPort = -1
	userServiceHost = ""
	Error = err

	if err != nil {
		return
	}

	userServiceREST, err := userService.MappedPort(context.Background(), "8080")
	if err != nil {
		Error = err
		return
	}
	userServiceRESTPort = userServiceREST.Int()

	userServiceGRPC, err := userService.MappedPort(context.Background(), "8081")
	if err != nil {
		Error = err
		return
	}
	userServiceGRPCPort = userServiceGRPC.Int()

	userServiceHost, err = userService.Host(context.Background())
	if err != nil {
		Error = err
		return
	}

	return
}
