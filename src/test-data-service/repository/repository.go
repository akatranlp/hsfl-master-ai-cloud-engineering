package repository

type Repository interface {
	ResetDatabase() error
}
