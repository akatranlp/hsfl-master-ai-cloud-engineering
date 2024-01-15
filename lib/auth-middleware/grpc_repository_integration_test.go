package auth_middleware

import (
	"bytes"
	"context"
	"encoding/json"
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
	testUserPassword := "test_password"

	postgres,
		postgresHost,
		err := containerhelpers.GetPostgresData(containerhelpers.StartPostgres(true))
	if err != nil {
		t.Fatalf("could not start postgres container: %s", err.Error())
	}

	t.Cleanup(func() {
		postgres.Terminate(context.Background())
	})

	testDataService,
		testDataServicePort,
		err := containerhelpers.GetTestDataServiceData(containerhelpers.StartTestDataService(postgresHost, "5432", testUserPassword, true, "../../src/test-data-service/init.sql"))
	if err != nil {
		t.Fatalf("could not start test-data-service container: %s", err.Error())
	}

	t.Cleanup(func() {
		testDataService.Terminate(context.Background())
	})

	userService,
		userServiceRESTPort,
		userServiceGRPCPort,
		userServiceHost,
		err := containerhelpers.GetUserServiceData(containerhelpers.StartUserService(postgresHost, "5432", "true", "true"))
	if err != nil {
		t.Fatalf("could not start user service container: %s", err.Error())
	}

	t.Cleanup(func() {
		userService.Terminate(context.Background())
	})

	userConn, err := grpc.Dial(fmt.Sprintf("%s:%d", userServiceHost, userServiceGRPCPort), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("could not connect to user_service: %v", err)
	}
	defer userConn.Close()

	client := proto.NewUserServiceClient(userConn)

	repository := NewGRPCRepository(client)
	valid_token, err := generateValidToken(userServiceRESTPort, testUserPassword)
	if err != nil {
		t.Fatalf("could not generate valid token: %s", err.Error())
	}

	t.Cleanup(resetDatabase(testDataServicePort))

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

		t.Run("should return userId if token is valid", func(t *testing.T) {
			// when
			userId, err := repository.VerifyToken(valid_token)

			// then
			assert.Equal(t, uint64(2), userId)
			assert.NoError(t, err)
		})
	})
}

func generateValidToken(port int, testUserPassword string) (string, error) {
	body := map[string]interface{}{
		"email":    "test",
		"password": testUserPassword,
	}

	reqBody, err := json.Marshal(body)
	if err != nil {
		return "", err
	}

	res, err := http.Post(fmt.Sprintf("http://localhost:%d/api/v1/login", port), "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return "", err
	}

	if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf("could not login: %s", res.Status)
	}

	var resBody map[string]interface{}
	json.NewDecoder(res.Body).Decode(&resBody)

	token := resBody["access_token"].(string)
	return token, nil
}

func resetDatabase(testDataServicePort int) func() {
	return func() {
		http.Post(fmt.Sprintf("http://localhost:%d/api/v1/reset", testDataServicePort), "application/json", nil)
	}
}
