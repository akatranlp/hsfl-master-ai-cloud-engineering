package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"

	auth_middleware "github.com/akatranlp/hsfl-master-ai-cloud-engineering/lib/auth-middleware"
	"github.com/akatranlp/hsfl-master-ai-cloud-engineering/lib/database"
	"github.com/akatranlp/hsfl-master-ai-cloud-engineering/lib/health"
	"github.com/akatranlp/hsfl-master-ai-cloud-engineering/transaction-service/api/router"
	book_service_client "github.com/akatranlp/hsfl-master-ai-cloud-engineering/transaction-service/book-service-client"
	"github.com/akatranlp/hsfl-master-ai-cloud-engineering/transaction-service/transactions"
	user_service_client "github.com/akatranlp/hsfl-master-ai-cloud-engineering/transaction-service/user-service-client"
	"github.com/caarlos0/env/v10"
	"github.com/joho/godotenv"
)

type ApplicationConfig struct {
	Database            database.PsqlConfig `envPrefix:"POSTGRES_"`
	Port                uint16              `env:"PORT" envDefault:"8080"`
	AuthUrlEndpoint     url.URL             `env:"AUTH_URL_ENDPOINT,notEmpty"`
	BookServiceEndpoint url.URL             `env:"BOOK_SERVICE_ENDPOINT,notEmpty"`
	UserServiceEndpoint url.URL             `env:"USER_SERVICE_ENDPOINT,notEmpty"`
}

func main() {
	godotenv.Load()

	config := ApplicationConfig{}
	if err := env.Parse(&config); err != nil {
		log.Fatalf("Couldn't parse environment %s", err.Error())
	}

	transactionRepository, err := transactions.NewPsqlRepository(config.Database)
	if err != nil {
		log.Fatalf("could not create user repository: %s", err.Error())
	}

	authRepository := auth_middleware.NewHTTPRepository(&config.AuthUrlEndpoint, http.DefaultClient)
	authController := auth_middleware.NewDefaultController(authRepository)
	healthController := health.NewDefaultController()

	bookServiceClientRepository := book_service_client.NewHTTPRepository(&config.BookServiceEndpoint, http.DefaultClient)
	userServiceClientRepository := user_service_client.NewHTTPRepository(&config.UserServiceEndpoint, http.DefaultClient)

	controller := transactions.NewDefaultController(transactionRepository, bookServiceClientRepository, userServiceClientRepository)

	handler := router.New(controller, authController, healthController)

	if err := transactionRepository.Migrate(); err != nil {
		log.Fatalf("could not migrate: %s", err.Error())
	}

	log.Println("Server Started!")

	addr := fmt.Sprintf("0.0.0.0:%d", config.Port)
	if err := http.ListenAndServe(addr, handler); err != nil {
		log.Fatalf("error while listen and serve: %s", err.Error())
	}
}
