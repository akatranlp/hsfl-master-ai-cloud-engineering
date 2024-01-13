### Service for financial transaction management

Purchase Your Own Currency with Real Money: Our platform allows you to acquire our native currency using real-world currency, giving you the flexibility to engage in various services.

Buying Chapters: With the acquired currency, you can access a diverse range of texts and content, enriching your user experience.

Authors Receiving Support through Donations: Users have the option to express appreciation and support their favorite authors by donating currency directly to them, creating a vibrant and rewarding community for creators.

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
