### Service for everything related to chapters and their content

Write: The Text Service provides a platform for users to compose and create written content, facilitating the process of generating original texts.

Edit: Users can make revisions and modifications to their texts using the editing features offered by the service, ensuring content accuracy and quality.

Draft: This service allows users to save their work as drafts, providing a safe and convenient space to work on chapters before publishing them.

Publish: Users can release their completed chapters to a wider audience by publishing them on the platform, making their content accessible to readers.

Read: The Text Service also functions as a reading platform, enabling users to explore and enjoy a diverse range of texts authored by others.

## How to use the Book-Service

- Firstly the book-service needs the user-service and the database to be running before you can start it.
- Secondly the transaction-service must also be running for full functionality
- Then the following environment variables must be correctly set to point to the other services

```bash
AUTH_IS_ACTIVE=true
AUTH_SERVICE_ENDPOINT=<user-service-grpc-endpoint>
TRANSACTION_SERVICE_ENDPOINT=<transaction-service-grpc-endpoint>
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
docker build -t book-service -f ./src/book-service/Dockerfile .
docker run -dit -p port:port -p grpcPort:grpcPort --envFile envFile book-service
```
