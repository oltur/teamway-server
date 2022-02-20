package model

import (
	"errors"
	"github.com/oltur/teamway-server/types"
	"github.com/oltur/teamway-server/utils"
	"time"

	"github.com/rs/xid"
)

type User struct {
	ID           types.Id    `json:"id"`
	UserName     string      `json:"userName" example:"user_name"`
	PasswordHash string      `json:"passwordHash"`
	TestsTaken   []TestTaken `json:"testsTaken"`
	Token        string      `json:"token"`
	TokenExpires int64       `json:"tokenExpires"`
}

//func UsersAll() (res []*User, err error) {
//	allUsers := GetMapValuesForUsers(usersByIds)
//	res = allUsers
//	return
//}

func UserOne(id types.Id) (res *User, err error) {
	for k := range usersByIds {
		if id == k {
			res = usersByIds[k]
			return
		}
	}
	return nil, ErrNotFound
}

func UserInsert(req *AddUserRequest) (res *User, err error) {
	if err = req.Validation(); err != nil {
		return
	}
	user := &User{
		ID:           types.Id(xid.New().String()),
		UserName:     req.UserName,
		PasswordHash: utils.Hash(req.Password),
		Token:        "",
		TokenExpires: 0,
	}

	_, err = UserOne(user.ID)
	if err != nil {
		if !errors.Is(err, ErrNotFound) {
			return
		}
	} else {
		err = ErrIdExists
		return
	}

	isFree, err := IsUserNameFree(req.UserName)
	if err != nil {
		return
	}
	if !isFree {
		err = ErrUserNameExists
		return
	}

	usersByIds[user.ID] = user
	res = user
	return
}

// UserUpdate part of CRUD
func UserUpdate(req *UpdateUserRequest) (err error) {
	user := usersByIds[req.ID]

	user.PasswordHash = utils.Hash(req.Password)

	err = UserSave(user)
	if err != nil {
		return
	}
	return
}

// UserSave Internal use only
func UserSave(req *User) (err error) {
	usersByIds[req.ID] = req
	return
}

func UserDelete(id types.Id) (err error) {
	if _, ok := usersByIds[id]; ok {
		err = ErrNotFound
		return
	}
	delete(usersByIds, id)
	return
}

var usersByIds = make(map[types.Id]*User)

//func GetMapValuesForUsers(m map[types.Id]*User) (res []*User) {
//	res = make([]*User, len(m))
//	i := 0
//	for _, v := range m {
//		res[i] = v
//		i++
//	}
//	return res
//}

func IsUserNameFree(userName string) (res bool, err error) {
	for k := range usersByIds {
		if usersByIds[k].UserName == userName {
			res = false
			return
		}
	}
	res = true
	return
}

func GetUserByCredentials(userName string, password string) (res *User, err error) {
	passwordHash := utils.Hash(password)
	for k := range usersByIds {
		if usersByIds[k].UserName == userName && usersByIds[k].PasswordHash == passwordHash {
			res = usersByIds[k]
			return
		}
	}
	err = ErrNotFound
	return
}

func UserLogout(id types.Id) (err error) {
	for k := range usersByIds {
		if usersByIds[k].ID == id {
			usersByIds[k].Token = ""
			usersByIds[k].TokenExpires = 0
			return
		}
	}
	err = ErrNotFound
	return
}

func VerifyToken(userId string, token string, expires int64) (res bool, err error) {
	if expires < time.Now().UnixMilli() {
		res = false
		return
	}
	for k := range usersByIds {
		if usersByIds[k].ID == types.Id(userId) && usersByIds[k].Token == token && usersByIds[k].TokenExpires == expires {
			res = true
			return
		}
	}
	err = ErrNotFound
	return
}
