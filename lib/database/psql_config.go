package database

import (
	"fmt"
)

type PsqlConfig struct {
	Host     string `env:"HOST,notEmpty"`
	Port     int    `env:"PORT,notEmpty"`
	Username string `env:"USER,notEmpty"`
	Password string `env:"PASSWORD,notEmpty"`
	Database string `env:"DB,notEmpty"`
}

func (config PsqlConfig) Dsn() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		config.Host, config.Port, config.Username, config.Password, config.Database)
}
