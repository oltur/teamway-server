package model

type AddUserRequest struct {
	UserName string `json:"userName" example:"user_name"`
	Password string `json:"password"`
}

func (a AddUserRequest) Validation() (err error) {
	if a.UserName == "" {
		err = ErrInvalidUserName
		return
	}
	if a.Password == "" {
		err = ErrInvalidPassword
		return
	}
	return
}
