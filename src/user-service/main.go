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
	Database     database.PsqlConfig `envPrefix:"POSTGRES_"`
	Jwt          auth.JwtConfig      `envPrefix:"JWT_"`
	AuthIsActive bool                `env:"AUTH_IS_ACTIVE" envDefault:"false"`
	PORT         uint16              `env:"PORT" envDefault:"8080"`
}

func main() {
	log.Println("Starting Server...")
	log.Println("Loading env...")

	godotenv.Load()

	log.Println("Parsing Config...")

	config := ApplicationConfig{}
	if err := env.Parse(&config); err != nil {
		log.Fatalf("Couldn't parse environment %s", err.Error())
	}

	log.Println("Create token Generator...")

	tokenGenerator, err := auth.NewJwtTokenGenerator(config.Jwt)
	if err != nil {
		log.Fatalf("could not create JWT token generator: %s", err.Error())
	}
	log.Println("Create user Repo - Connection to db...")

	userRepository, err := user.NewPsqlRepository(config.Database)
	if err != nil {
		log.Fatalf("could not create user repository: %s", err.Error())
	}
	log.Println("Starting DB Migration...")

	if err := userRepository.Migrate(); err != nil {
		log.Fatalf("could not migrate: %s", err.Error())
	}

	log.Println("Create Hasher...")

	hasher := crypto.NewBcryptHasher()
	log.Println("Create Health Controller...")

	healthController := health.NewDefaultController()
	log.Println("Create Defualt Controller...")
	controller := user.NewDefaultController(userRepository, hasher, tokenGenerator, config.AuthIsActive)

	log.Println("Create Router...")
	handler := router.New(controller, healthController)

	log.Println("Server Started!")

	addr := fmt.Sprintf("0.0.0.0:%d", config.PORT)
	if err := http.ListenAndServe(addr, handler); err != nil {
		log.Fatalf("error while listen and serve: %s", err.Error())
	}
}
