package user_service_client

import (
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
		err := containerhelpers.GetTestDataServiceData(containerhelpers.StartTestDataService(postgresHost, "5432", testUserPassword, true, "../../test-data-service/init.sql"))
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
		err := containerhelpers.GetUserServiceData(containerhelpers.StartUserService(postgresHost, "5432", "false", "true"))
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

	t.Cleanup(resetDatabase(testDataServicePort))

	t.Run("MoveBalance", func(t *testing.T) {

		t.Run("Should return error when user1 does not exist", func(t *testing.T) {
			t.Cleanup(resetDatabase(testDataServicePort))

			// given
			shouldUsers := getAllUsers(t, userServiceRESTPort)

			// when
			err := repository.MoveBalance(1000, 2, 100)

			// then
			assert.Error(t, err)
			users := getAllUsers(t, userServiceRESTPort)
			assert.Equal(t, shouldUsers, users)
		})

		t.Run("Should return error when user2 does not exist", func(t *testing.T) {
			t.Cleanup(resetDatabase(testDataServicePort))

			// given
			shouldUsers := getAllUsers(t, userServiceRESTPort)

			// when
			err := repository.MoveBalance(1, 1000, 100)

			// then
			assert.Error(t, err)
			users := getAllUsers(t, userServiceRESTPort)
			assert.Equal(t, shouldUsers, users)
		})

		t.Run("Should return error when balance is not enough", func(t *testing.T) {
			t.Skip("TODO: implement this test")
			t.Cleanup(resetDatabase(testDataServicePort))

			// given
			shouldUsers := getAllUsers(t, userServiceRESTPort)

			// when
			err := repository.MoveBalance(1, 2, 100000)

			// then
			assert.Error(t, err)
			users := getAllUsers(t, userServiceRESTPort)
			assert.Equal(t, shouldUsers, users)
		})

		t.Run("Should move balance", func(t *testing.T) {
			t.Cleanup(resetDatabase(testDataServicePort))

			// given
			shouldUsers := getAllUsers(t, userServiceRESTPort)

			// when
			err := repository.MoveBalance(1, 2, 1000)

			// then
			assert.NoError(t, err)
			users := getAllUsers(t, userServiceRESTPort)
			assert.Equal(t, shouldUsers[0]["balance"].(float64)-1000, users[0]["balance"].(float64))
			assert.Equal(t, shouldUsers[1]["balance"].(float64)+1000, users[1]["balance"].(float64))
		})
	})
}

func getAllUsers(t *testing.T, port int) []map[string]interface{} {
	res, err := http.Get(fmt.Sprintf("http://localhost:%d/api/v1/users", port))
	assert.NoError(t, err)

	if res.StatusCode != http.StatusOK {
		assert.True(t, false)
	}

	var resBody []map[string]interface{}
	err = json.NewDecoder(res.Body).Decode(&resBody)

	assert.NoError(t, err)
	assert.Len(t, resBody, 2)
	assert.Equal(t, uint64(1), uint64(resBody[0]["id"].(float64)))
	assert.Equal(t, uint64(2), uint64(resBody[1]["id"].(float64)))

	return resBody
}

func resetDatabase(testDataServicePort int) func() {
	return func() {
		http.Post(fmt.Sprintf("http://localhost:%d/api/v1/reset", testDataServicePort), "application/json", nil)
	}
}
