package transaction_service_client

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"testing"

	"github.com/akatranlp/hsfl-master-ai-cloud-engineering/lib/containerhelpers"
	"github.com/akatranlp/hsfl-master-ai-cloud-engineering/lib/grpc/transaction-service/proto"
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

	transactionService,
		_,
		transactionServiceGRPCPort,
		transactionServiceHost,
		err := containerhelpers.GetTransactionServiceData(containerhelpers.StartTransactionService(postgresHost, "5432", "false", "true"))
	if err != nil {
		t.Fatalf("could not start transaction service container: %s", err.Error())
	}

	t.Cleanup(func() {
		transactionService.Terminate(context.Background())
	})

	transactionConn, err := grpc.Dial(fmt.Sprintf("%s:%d", transactionServiceHost, transactionServiceGRPCPort), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("could not connect to transaction_service: %v", err)
	}
	defer transactionConn.Close()

	transactionServiceClient := proto.NewTransactionServiceClient(transactionConn)

	repository := NewGRPCRepository(transactionServiceClient)

	t.Cleanup(resetDatabase(testDataServicePort))

	t.Run("CheckChapterBought", func(t *testing.T) {
		t.Run("should return error when chapter is not found", func(t *testing.T) {
			t.Cleanup(resetDatabase(testDataServicePort))

			// given
			// when
			err := repository.CheckChapterBought(1, 1000, 1)

			// then
			assert.Error(t, err)
		})

		t.Run("should return error when book is not found", func(t *testing.T) {
			t.Cleanup(resetDatabase(testDataServicePort))

			// given
			// when
			err := repository.CheckChapterBought(1, 1, 1000)

			// then
			assert.Error(t, err)
		})

		t.Run("should return error when chapter is not bought", func(t *testing.T) {
			t.Cleanup(resetDatabase(testDataServicePort))

			// given
			// when
			err := repository.CheckChapterBought(2, 3, 1)

			// then
			assert.Error(t, err)
		})

		t.Run("should return no error when chapter is already bought", func(t *testing.T) {
			t.Cleanup(resetDatabase(testDataServicePort))

			// given
			// when
			err := repository.CheckChapterBought(2, 1, 1)

			// then
			assert.NoError(t, err)
		})
	})
}

func resetDatabase(testDataServicePort int) func() {
	return func() {
		http.Post(fmt.Sprintf("http://localhost:%d/api/v1/reset", testDataServicePort), "application/json", nil)
	}
}
