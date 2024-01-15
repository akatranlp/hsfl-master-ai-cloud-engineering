package containerhelpers

import (
	"context"
	"os"
	"path"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func StartTestDataService(postgresHost string, postgresPort string, testUserPassword string, resetOnInit bool, relativPathToInitSQL string) (testcontainers.Container, error) {
	var resetOnInitString string
	if resetOnInit {
		resetOnInitString = "true"
	} else {
		resetOnInitString = "false"
	}

	dir, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	sqlFile := path.Join(dir, relativPathToInitSQL)

	req := testcontainers.ContainerRequest{
		Image:        "akatranlp/test-data-service:latest",
		ExposedPorts: []string{"8080/tcp"},
		Env: map[string]string{
			"RESET_ON_INIT":           resetOnInitString,
			"PORT":                    "8080",
			"POSTGRES_HOST":           postgresHost,
			"POSTGRES_PORT":           postgresPort,
			"TEST_DATA_USER_PASSWORD": testUserPassword,
			"TEST_DATA_FILE_PATH":     "/init.sql",
			"POSTGRES_USER":           "postgres",
			"POSTGRES_PASSWORD":       "postgres",
			"POSTGRES_DB":             "postgres",
		},
		Mounts: []testcontainers.ContainerMount{
			{
				Source:   testcontainers.GenericBindMountSource{HostPath: sqlFile},
				Target:   "/init.sql",
				ReadOnly: true,
			},
		},
		WaitingFor: wait.ForAll(wait.ForListeningPort("8080/tcp"), wait.ForLog(".*REST-Server started on Port:.*").AsRegexp().WithStartupTimeout(60+time.Second)).WithStartupTimeoutDefault(60 * time.Second),
	}

	return testcontainers.GenericContainer(context.Background(), testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
}

func GetTestDataServiceData(testDataService testcontainers.Container, err error) (testcontainers.Container, int, error) {
	if err != nil {
		return nil, -1, err
	}
	testDataServicePort, err := testDataService.MappedPort(context.Background(), "8080")
	if err != nil {
		return nil, -1, err
	}
	return testDataService, testDataServicePort.Int(), nil
}
