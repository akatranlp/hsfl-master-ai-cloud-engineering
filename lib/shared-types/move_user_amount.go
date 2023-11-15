package shared_types

type MoveBalanceRequest struct {
	UserId          uint64 `json:"userId"`
	ReceivingUserId uint64 `json:"receivingUserId"`
	Amount          int64  `json:"amount"`
}

func (r *MoveBalanceRequest) IsValid() bool {
	return r.UserId != 0 && r.ReceivingUserId != 0
}

type MoveBalanceResponse struct {
	Success bool `json:"success"`
}
