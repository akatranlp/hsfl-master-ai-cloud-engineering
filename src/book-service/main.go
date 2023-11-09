package main

import (
	"fmt"
	"github.com/akatranlp/hsfl-master-ai-cloud-engineering/book-service/api/router"
	"github.com/akatranlp/hsfl-master-ai-cloud-engineering/book-service/books"
	"github.com/akatranlp/hsfl-master-ai-cloud-engineering/book-service/chapters"
	authMiddleware "github.com/akatranlp/hsfl-master-ai-cloud-engineering/lib/auth-middleware"
	"github.com/akatranlp/hsfl-master-ai-cloud-engineering/lib/database"
	"github.com/caarlos0/env/v10"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"net/url"
)

type ApplicationConfig struct {
	Database        database.PsqlConfig `envPrefix:"POSTGRES_"`
	Port            uint16              `env:"PORT" envDefault:"8080"`
	AuthUrlEndpoint url.URL             `env:"AUTH_URL_ENDPOINT,notEmpty"`
}

func main() {
	godotenv.Load()

	config := ApplicationConfig{}
	if err := env.Parse(&config); err != nil {
		log.Fatalf("Couldn't parse environment %s", err.Error())
	}

	bookRepository, err := books.NewPsqlRepository(config.Database)
	if err != nil {
		log.Fatalf("could not instanciate bookRepo: %s", err.Error())
	}
	chapterRepository, err := chapters.NewPsqlRepository(config.Database)
	if err != nil {
		log.Fatalf("could not instanciate chapterRepo: %s", err.Error())
	}

	authRepository := authMiddleware.NewHTTPRepository(&config.AuthUrlEndpoint, http.DefaultClient)
	bookController := books.NewDefaultController(bookRepository)
	chapterController := chapters.NewDefaultController(chapterRepository)
	authController := authMiddleware.NewDefaultController(authRepository)

	handler := router.New(authController, bookController, chapterController)

	if err := bookRepository.Migrate(); err != nil {
		log.Fatalf("could not migrate: %s", err.Error())
	}
	if err := chapterRepository.Migrate(); err != nil {
		log.Fatalf("could not migrate: %s", err.Error())
	}

	log.Println("Server Started!")

	addr := fmt.Sprintf("0.0.0.0:%d", config.Port)
	if err := http.ListenAndServe(addr, handler); err != nil {
		log.Fatalf("error while listen and serve: %s", err.Error())
	}
}
