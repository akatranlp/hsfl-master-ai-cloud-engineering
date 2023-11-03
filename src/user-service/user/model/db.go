package model

type UpdateUser struct {
	ProfileName string `json:"profileName"`
	Balance     int64  `json:"balance"`
}

type DbUser struct {
	ID          uint64 `json:"id"`
	Email       string `json:"email"`
	Password    []byte `json:"password"`
	Username    string `json:"username"`
	ProfileName string `json:"profileName"`
	Balance     int64  `json:"balance"`
}
