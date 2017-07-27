package core

import (
	"crypto/sha1"
	"encoding/hex"

	"github.com/crazyprograms/sud/corebase"
)

const salt = "alsHF;adsgh;aiygh;dashgdasbnvgn"

type User struct {
	Login        string
	HashPassword string
	Access       map[string]bool
}

func GenHashPassword(Password string) string {
	b := ([]byte)(salt + Password)
	b1 := sha1.Sum(b)
	return hex.EncodeToString(b1[:])
}
func (u *User) GetLogin() string { return u.Login }
func (u *User) GetCheckPassword(Password string) bool {
	return GenHashPassword(Password) == u.HashPassword
}
func (u *User) CheckAccess(Access string) bool {
	A, ok := u.Access[Access]
	if ok {
		return A
	}
	return false
}
func (u *User) Users() []corebase.IUser {
	return []corebase.IUser{u}
}

var _u corebase.IUser = (corebase.IUser)(&User{})
