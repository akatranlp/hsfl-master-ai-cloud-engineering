package auth_middleware

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"testing"

	"github.com/akatranlp/hsfl-master-ai-cloud-engineering/lib/containerhelpers"
	"github.com/akatranlp/hsfl-master-ai-cloud-engineering/lib/grpc/user-service/proto"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func TestIntegrationGRPCRepository(t *testing.T) {
	postgres, err := containerhelpers.StartPostgres(true)
	if err != nil {
		t.Fatalf("could not start postgres container: %s", err.Error())
	}

	t.Cleanup(func() {
		postgres.Terminate(context.Background())
	})

	postgresHost, err := postgres.ContainerIP(context.Background())
	if err != nil {
		t.Fatalf("could not get database container host: %s", err.Error())
	}

	user_service, err := containerhelpers.StartUserService(postgresHost, "5432")
	if err != nil {
		t.Fatalf("could not start user service container: %s", err.Error())
	}

	t.Cleanup(func() {
		user_service.Terminate(context.Background())
	})

	userServiceRESTPort, err := user_service.MappedPort(context.Background(), "8080")
	if err != nil {
		t.Fatalf("could not get user-service REST container port: %s", err.Error())
	}

	userServiceGRPCPort, err := user_service.MappedPort(context.Background(), "8081")
	if err != nil {
		t.Fatalf("could not get user-service gRPC container port: %s", err.Error())
	}

	userServiceHost, err := user_service.Host(context.Background())
	if err != nil {
		t.Fatalf("could not get database container host: %s", err.Error())
	}

	userConn, err := grpc.Dial(fmt.Sprintf("%s:%s", userServiceHost, userServiceGRPCPort.Port()), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("could not connect to user_service: %v", err)
	}
	defer userConn.Close()

	client := proto.NewUserServiceClient(userConn)

	repository := NewGRPCRepository(client)
	valid_token := generateValidToken(userServiceRESTPort.Int())
	_ = valid_token

	t.Run("VerifyToken", func(t *testing.T) {
		t.Run("should return error if token is invalid", func(t *testing.T) {
			// given
			token := "invalid_token"

			// when
			userId, err := repository.VerifyToken(token)

			// then
			assert.Equal(t, uint64(0), userId)
			assert.Error(t, err)
		})
	})
}

func generateValidToken(port int) string {
	data := map[string]interface{}{
		"username": "test",
		"password": "test",
	}

	_ = data

	http.Post(fmt.Sprintf("http://localhost:%d/users", port), "application/json", nil)
	return "valid"
}
