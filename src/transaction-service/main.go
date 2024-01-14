package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"

	auth_middleware "github.com/akatranlp/hsfl-master-ai-cloud-engineering/lib/auth-middleware"
	"github.com/akatranlp/hsfl-master-ai-cloud-engineering/lib/database"
	bproto "github.com/akatranlp/hsfl-master-ai-cloud-engineering/lib/grpc/book-service/proto"
	tproto "github.com/akatranlp/hsfl-master-ai-cloud-engineering/lib/grpc/transaction-service/proto"
	uproto "github.com/akatranlp/hsfl-master-ai-cloud-engineering/lib/grpc/user-service/proto"
	"github.com/akatranlp/hsfl-master-ai-cloud-engineering/lib/health"
	"github.com/akatranlp/hsfl-master-ai-cloud-engineering/transaction-service/api/router"
	book_service_client "github.com/akatranlp/hsfl-master-ai-cloud-engineering/transaction-service/book-service-client"
	grpc_server "github.com/akatranlp/hsfl-master-ai-cloud-engineering/transaction-service/grpc"
	"github.com/akatranlp/hsfl-master-ai-cloud-engineering/transaction-service/service"
	transactions_controller "github.com/akatranlp/hsfl-master-ai-cloud-engineering/transaction-service/transactions/controller"
	transactions_repository "github.com/akatranlp/hsfl-master-ai-cloud-engineering/transaction-service/transactions/repository"
	user_service_client "github.com/akatranlp/hsfl-master-ai-cloud-engineering/transaction-service/user-service-client"
	"github.com/caarlos0/env/v10"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

type ApplicationConfig struct {
	Database            database.PsqlConfig `envPrefix:"POSTGRES_"`
	Port                uint16              `env:"PORT" envDefault:"8080"`
	GrpcPort            uint16              `env:"GRPC_PORT" envDefault:"8081"`
	GrpcCommunication   bool                `env:"GRPC_COMMUNICATION" envDefault:"true"`
	AuthIsActive        bool                `env:"AUTH_IS_ACTIVE" envDefault:"false"`
	AuthServiceEndpoint url.URL             `env:"AUTH_SERVICE_ENDPOINT,notEmpty"`
	BookServiceEndpoint url.URL             `env:"BOOK_SERVICE_ENDPOINT,notEmpty"`
	UserServiceEndpoint url.URL             `env:"USER_SERVICE_ENDPOINT,notEmpty"`
}

func main() {
	godotenv.Load()

	config := ApplicationConfig{}
	if err := env.Parse(&config); err != nil {
		log.Fatalf("Couldn't parse environment %s", err.Error())
	}

	transactionRepository, err := transactions_repository.NewPsqlRepository(config.Database)
	if err != nil {
		log.Fatalf("could not create user repository: %s", err.Error())
	}

	var authRepository auth_middleware.Repository
	var bookServiceClientRepository book_service_client.Repository
	var userServiceClientRepository user_service_client.Repository

	if config.GrpcCommunication {
		userConn, err := grpc.Dial(config.AuthServiceEndpoint.Host, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Fatalf("could not connect: %v", err)
		}
		defer userConn.Close()

		bookConn, err := grpc.Dial(config.BookServiceEndpoint.Host, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Fatalf("could not connect: %v", err)
		}
		defer bookConn.Close()

		userGrpcClient := uproto.NewUserServiceClient(userConn)
		bookGrpcClient := bproto.NewBookServiceClient(bookConn)

		authRepository = auth_middleware.NewGRPCRepository(userGrpcClient)
		bookServiceClientRepository = book_service_client.NewGRPCRepository(bookGrpcClient)
		userServiceClientRepository = user_service_client.NewGRPCRepository(userGrpcClient)
	} else {
		authRepository = auth_middleware.NewHTTPRepository(&config.AuthServiceEndpoint, http.DefaultClient)
		bookServiceClientRepository = book_service_client.NewHTTPRepository(&config.BookServiceEndpoint, http.DefaultClient)
		userServiceClientRepository = user_service_client.NewHTTPRepository(&config.UserServiceEndpoint, http.DefaultClient)
	}

	service := service.NewDefaultService(transactionRepository)

	authController := auth_middleware.NewDefaultController(authRepository, config.AuthIsActive)
	healthController := health.NewDefaultController()

	controller := transactions_controller.NewDefaultController(transactionRepository, bookServiceClientRepository, userServiceClientRepository, service)

	handler := router.New(controller, authController, healthController)

	if err := transactionRepository.Migrate(); err != nil {
		log.Fatalf("could not migrate: %s", err.Error())
	}

	if config.GrpcCommunication {
		grpcAddr := fmt.Sprintf("0.0.0.0:%d", config.GrpcPort)
		listener, err := net.Listen("tcp", grpcAddr)
		if err != nil {
			log.Fatalf("could not listen: %v", err)
		}

		srv := grpc.NewServer()
		reflection.Register(srv)
		grpcServer := grpc_server.NewServer(service)
		tproto.RegisterTransactionServiceServer(srv, grpcServer)

		go func() {
			log.Printf("GRPC-Server started on Port: %d\n", config.GrpcPort)
			if err := srv.Serve(listener); err != nil {
				log.Fatalf("could not serve: %v", err)
			}
		}()
	}

	log.Printf("REST-Server started on Port: %d\n", config.Port)
	addr := fmt.Sprintf("0.0.0.0:%d", config.Port)
	if err := http.ListenAndServe(addr, handler); err != nil {
		log.Fatalf("error while listen and serve: %s", err.Error())
	}
}
