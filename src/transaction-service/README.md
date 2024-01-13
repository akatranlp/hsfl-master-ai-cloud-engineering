### Service for financial transaction management

Purchase currency: Our platform allows you to acquire our native currency, VV-Coins, using real-world currency, giving you the flexibility to engage in various services.

Buy Chapters: With the acquired currency, you can buy and access a diverse range of texts and content.

Earn Money: Every time someone buys a published work, you get the VV-Coins directly to your account, which you can convert back to your preferred currency

## How to use the Transaction-Service

- Firstly the transaction-service needs the user-service and the database to be running before you can start it.
- Secondly the book-service must also be running for full functionality
- Then the following environment variables must be correctly set to point to the other services

```bash
AUTH_IS_ACTIVE=true
AUTH_SERVICE_ENDPOINT=<user-service-grpc-endpoint>
USER_SERVICE_ENDPOINT=<user-service-grpc-endpoint>
BOOK_SERVICE_ENDPOINT=<transaction-service-grpc-endpoint>
PORT=<HTTP-Port>
GRPC_PORT=<GRPC-Port>
POSTGRES_HOST=<db-hostname>
POSTGRES_PORT=<db-port>
POSTGRES_USER=<db-username>
POSTGRES_PASSWORD=<db-passwort>
POSTGRES_DB=<db-name>
```

- At last the service can be run with the following command: `go run main.go`

### Create Docker-Image

If you want to use an docker-image instead, the following commands must be executed from the root of this project:

```bash
docker build -t transaction-service -f ./src/transaction-service/Dockerfile .
docker run -dit -p port:port -p grpcPort:grpcPort --envFile envFile transaction-service
```
