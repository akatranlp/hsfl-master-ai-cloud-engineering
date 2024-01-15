package model

type DbUser struct {
	ID           uint64
	Email        string
	Password     []byte
	ProfileName  string
	Balance      int64
	TokenVersion uint64
}

type DbUserPatch struct {
	ID           *uint64
	Email        *string
	Password     *[]byte
	ProfileName  *string
	Balance      *int64
	TokenVersion *uint64
}

type UserDTO struct {
	ID          uint64 `json:"id"`
	Email       string `json:"email"`
	ProfileName string `json:"profileName"`
	Balance     int64  `json:"balance"`
}

func (user *DbUser) ToDto() UserDTO {
	return UserDTO{
		user.ID,
		user.Email,
		user.ProfileName,
		user.Balance,
	}
}
