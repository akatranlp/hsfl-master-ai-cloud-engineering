package user_service_client

type Repository interface {
	MoveBalance(userId uint64, receivingUserId uint64, amount int64) error
}
