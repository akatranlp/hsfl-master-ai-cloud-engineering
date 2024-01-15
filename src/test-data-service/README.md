### Service to load our Test-Data to the database

The Test-Data-Service provides the `/api/v1/reset`-endpoint, when posted to resets the database to the initial testdata.

## How to use the Test-Data-Service

- The database needs to be running before you can start it.
- Then the following environment variables must be correctly set

```bash
TEST_DATA_USER_PASSWORD=<passwort>
TEST_DATA_FILE_PATH=<path to init.sql>
RESET_ON_INIT=<true | false>
PORT=<HTTP-Port>
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
docker build -t test-data-service -f ./src/test-data-service/Dockerfile .
docker run -dit -v <sql.init-file>:<sql.init-file> -p port:port --envFile envFile test-data-service
```
