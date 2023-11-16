package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/akatranlp/hsfl-master-ai-cloud-engineering/lib/database"
	"github.com/akatranlp/hsfl-master-ai-cloud-engineering/lib/health"
	"github.com/akatranlp/hsfl-master-ai-cloud-engineering/user-service/api/router"
	"github.com/akatranlp/hsfl-master-ai-cloud-engineering/user-service/auth"
	"github.com/akatranlp/hsfl-master-ai-cloud-engineering/user-service/crypto"
	"github.com/akatranlp/hsfl-master-ai-cloud-engineering/user-service/user"
	"github.com/caarlos0/env/v10"
	"github.com/joho/godotenv"
)

type ApplicationConfig struct {
	Database database.PsqlConfig `envPrefix:"POSTGRES_"`
	Jwt      auth.JwtConfig      `envPrefix:"JWT_"`
	PORT     uint16              `env:"PORT" envDefault:"8080"`
}

func main() {
	godotenv.Load()

	config := ApplicationConfig{}
	if err := env.Parse(&config); err != nil {
		log.Fatalf("Couldn't parse environment %s", err.Error())
	}

	tokenGenerator, err := auth.NewJwtTokenGenerator(config.Jwt)
	if err != nil {
		log.Fatalf("could not create JWT token generator: %s", err.Error())
	}

	userRepository, err := user.NewPsqlRepository(config.Database)
	if err != nil {
		log.Fatalf("could not create user repository: %s", err.Error())
	}

	if err := userRepository.Migrate(); err != nil {
		log.Fatalf("could not migrate: %s", err.Error())
	}

	hasher := crypto.NewBcryptHasher()

	healthController := health.NewDefaultController()
	controller := user.NewDefaultController(userRepository, hasher, tokenGenerator)

	handler := router.New(controller, healthController)

	log.Println("Server Started!")

	addr := fmt.Sprintf("0.0.0.0:%d", config.PORT)
	if err := http.ListenAndServe(addr, handler); err != nil {
		log.Fatalf("error while listen and serve: %s", err.Error())
	}
}
