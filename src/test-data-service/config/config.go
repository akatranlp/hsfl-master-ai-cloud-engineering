package config

type Config interface {
	GetSqlString() (string, error)
	GetUserPassword() string
}
