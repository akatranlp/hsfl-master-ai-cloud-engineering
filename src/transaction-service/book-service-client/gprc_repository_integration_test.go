package book_service_client

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"testing"

	"github.com/akatranlp/hsfl-master-ai-cloud-engineering/lib/containerhelpers"
	"github.com/akatranlp/hsfl-master-ai-cloud-engineering/lib/grpc/book-service/proto"
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

	bookService,
		bookServiceRESTPort,
		bookServiceGRPCPort,
		bookServiceHost,
		err := containerhelpers.GetBookServiceData(containerhelpers.StartBookService(postgresHost, "5432", "false", "true"))
	if err != nil {
		t.Fatalf("could not start book service container: %s", err.Error())
	}

	t.Cleanup(func() {
		bookService.Terminate(context.Background())
	})

	bookConn, err := grpc.Dial(fmt.Sprintf("%s:%d", bookServiceHost, bookServiceGRPCPort), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("could not connect to book_service: %v", err)
	}
	defer bookConn.Close()

	bookServiceClient := proto.NewBookServiceClient(bookConn)

	repository := NewGRPCRepository(bookServiceClient)

	t.Cleanup(resetDatabase(testDataServicePort))

	t.Run("ValidateChapterId", func(t *testing.T) {
		t.Run("should return error when chapter is not found", func(t *testing.T) {
			t.Cleanup(resetDatabase(testDataServicePort))

			// given
			// when
			res, err := repository.ValidateChapterId(1, 1000, 1)

			// then
			assert.Error(t, err)
			assert.Nil(t, res)
		})

		t.Run("should return error when book is not found", func(t *testing.T) {
			t.Cleanup(resetDatabase(testDataServicePort))

			// given
			// when
			res, err := repository.ValidateChapterId(1, 1, 1000)

			// then
			assert.Error(t, err)
			assert.Nil(t, res)
		})

		t.Run("should return error when book is yours", func(t *testing.T) {
			t.Cleanup(resetDatabase(testDataServicePort))

			// given
			// when
			res, err := repository.ValidateChapterId(1, 1, 1)

			// then
			assert.Error(t, err)
			assert.Nil(t, res)
		})

		t.Run("should return no error when book is not yours", func(t *testing.T) {
			t.Cleanup(resetDatabase(testDataServicePort))

			// given
			shouldBooks := getAllBooksAndChapters(t, bookServiceRESTPort)
			log.Println(shouldBooks)

			// when
			res, err := repository.ValidateChapterId(2, 1, 1)

			// then
			assert.NoError(t, err)
			assert.Equal(t, uint64(shouldBooks[0]["chapters"].([]map[string]interface{})[0]["id"].(float64)), res.ChapterId)
			assert.Equal(t, uint64(shouldBooks[0]["id"].(float64)), res.BookId)
			assert.Equal(t, uint64(shouldBooks[0]["chapters"].([]map[string]interface{})[0]["price"].(float64)), res.Amount)
			assert.Equal(t, uint64(shouldBooks[0]["authorId"].(float64)), res.ReceivingUserId)
		})
	})
}

func getAllBooksAndChapters(t *testing.T, port int) []map[string]interface{} {
	res, err := http.Get(fmt.Sprintf("http://localhost:%d/api/v1/books", port))
	assert.NoError(t, err)

	if res.StatusCode != http.StatusOK {
		assert.True(t, false)
	}

	var resBody []map[string]interface{}
	err = json.NewDecoder(res.Body).Decode(&resBody)

	assert.NoError(t, err)
	assert.Len(t, resBody, 3)

	for _, book := range resBody {
		res, err := http.Get(fmt.Sprintf("http://localhost:%d/api/v1/books/%d/chapters", port, int(book["id"].(float64))))
		assert.NoError(t, err)

		if res.StatusCode != http.StatusOK {
			assert.True(t, false)
		}

		var chapters []map[string]interface{}
		err = json.NewDecoder(res.Body).Decode(&chapters)

		assert.NoError(t, err)
		book["chapters"] = chapters
	}

	return resBody
}

func resetDatabase(testDataServicePort int) func() {
	return func() {
		http.Post(fmt.Sprintf("http://localhost:%d/api/v1/reset", testDataServicePort), "application/json", nil)
	}
}
