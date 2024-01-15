package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/akatranlp/hsfl-master-ai-cloud-engineering/lib/crypto"
	"github.com/akatranlp/hsfl-master-ai-cloud-engineering/lib/database"
	"github.com/akatranlp/hsfl-master-ai-cloud-engineering/lib/health"
	router "github.com/akatranlp/hsfl-master-ai-cloud-engineering/test-data-service/api"
	"github.com/akatranlp/hsfl-master-ai-cloud-engineering/test-data-service/config"
	"github.com/akatranlp/hsfl-master-ai-cloud-engineering/test-data-service/controller"
	"github.com/akatranlp/hsfl-master-ai-cloud-engineering/test-data-service/repository"
	"github.com/caarlos0/env/v10"
	"github.com/joho/godotenv"
)

type ApplicationConfig struct {
	Database    database.PsqlConfig   `envPrefix:"POSTGRES_"`
	Port        uint16                `env:"PORT" envDefault:"8080"`
	TestData    config.TestDataConfig `envPrefix:"TEST_DATA_"`
	ResetOnInit bool                  `env:"RESET_ON_INIT" envDefault:"false"`
}

func main() {
	godotenv.Load()

	config := ApplicationConfig{}
	if err := env.Parse(&config); err != nil {
		log.Fatalf("Couldn't parse environment %s", err.Error())
	}

	hasher := crypto.NewBcryptHasher()

	repository, err := repository.NewPsqlRepository(config.Database, config.TestData, hasher)
	if err != nil {
		log.Fatalf("could not create repository: %s", err.Error())
	}

	if config.ResetOnInit {
		if err := repository.ResetDatabase(); err != nil {
			log.Fatalf("could not reset database: %s", err.Error())
		}
	}

	controller := controller.NewDefaultController(repository)

	healthController := health.NewDefaultController()

	handler := router.New(controller, healthController)

	log.Printf("REST-Server started on Port: %d\n", config.Port)
	addr := fmt.Sprintf("0.0.0.0:%d", config.Port)
	if err := http.ListenAndServe(addr, handler); err != nil {
		log.Fatalf("error while listen and serve: %s", err.Error())
	}
}
