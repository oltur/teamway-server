package model

type LoginResponse struct {
	Token        string `json:"token"`
	TokenExpires int64  `json:"tokenExpires"`
	UserId       string `json:"userId"`
	UserName     string `json:"userName"`
}
