package config

import "os"

type TestDataConfig struct {
	UserPassword string `env:"USER_PASSWORD,notEmpty"`
	FilePath     string `env:"FILE_PATH"`
	File         string `env:"FILE"`
}

func (config TestDataConfig) GetSqlString() (string, error) {
	var sqlString string
	if config.File != "" {
		sqlString = config.File
	} else {
		data, err := os.ReadFile(config.FilePath)
		if err != nil {
			return "", err
		}
		sqlString = string(data)
	}
	return sqlString, nil
}

func (config TestDataConfig) GetUserPassword() string {
	return config.UserPassword
}
