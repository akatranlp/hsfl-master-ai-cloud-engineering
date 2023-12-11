package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"

	"github.com/akatranlp/hsfl-master-ai-cloud-engineering/book-service/api/router"
	"github.com/akatranlp/hsfl-master-ai-cloud-engineering/book-service/books"
	"github.com/akatranlp/hsfl-master-ai-cloud-engineering/book-service/chapters"
	grpc_server "github.com/akatranlp/hsfl-master-ai-cloud-engineering/book-service/grpc"
	transaction_service_client "github.com/akatranlp/hsfl-master-ai-cloud-engineering/book-service/transaction-service-client"
	authMiddleware "github.com/akatranlp/hsfl-master-ai-cloud-engineering/lib/auth-middleware"
	"github.com/akatranlp/hsfl-master-ai-cloud-engineering/lib/database"
	"github.com/akatranlp/hsfl-master-ai-cloud-engineering/lib/grpc/book-service/proto"
	"github.com/akatranlp/hsfl-master-ai-cloud-engineering/lib/health"
	"github.com/caarlos0/env/v10"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type ApplicationConfig struct {
	Database                  database.PsqlConfig `envPrefix:"POSTGRES_"`
	Port                      uint16              `env:"PORT" envDefault:"8080"`
	GrpcPort                  uint16              `env:"GRPC_PORT" envDefault:"8081"`
	AuthUrlEndpoint           url.URL             `env:"AUTH_URL_ENDPOINT,notEmpty"`
	TransactionServiceBaseUrl url.URL             `env:"TRANSACTION_SERVICE_ENDPOINT,notEmpty"`
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
	transactionServiceClient := transaction_service_client.NewHTTPRepository(&config.TransactionServiceBaseUrl, http.DefaultClient)
	chapterController := chapters.NewDefaultController(chapterRepository, transactionServiceClient)
	authController := authMiddleware.NewDefaultController(authRepository)
	healthController := health.NewDefaultController()

	handler := router.New(authController, bookController, chapterController, healthController)

	if err := bookRepository.Migrate(); err != nil {
		log.Fatalf("could not migrate: %s", err.Error())
	}
	if err := chapterRepository.Migrate(); err != nil {
		log.Fatalf("could not migrate: %s", err.Error())
	}

	grpcAddr := fmt.Sprintf("0.0.0.0:%d", config.GrpcPort)
	listener, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		log.Fatalf("could not listen: %v", err)
	}

	srv := grpc.NewServer()
	reflection.Register(srv)
	grpcServer := grpc_server.NewServer(bookRepository, chapterRepository)
	proto.RegisterBookServiceServer(srv, grpcServer)

	go func() {
		log.Printf("GRPC-Server started on Port: %d\n", config.GrpcPort)
		if err := srv.Serve(listener); err != nil {
			log.Fatalf("could not serve: %v", err)
		}
	}()

	log.Printf("REST-Server started on Port: %d\n", config.Port)
	addr := fmt.Sprintf("0.0.0.0:%d", config.Port)
	if err := http.ListenAndServe(addr, handler); err != nil {
		log.Fatalf("error while listen and serve: %s", err.Error())
	}
}
