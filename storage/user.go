package storage

import "crypto/sha1"
import "encoding/hex"

const salt = "alsHF;adsgh;aiygh;dashgdasbnvgn"

type User struct {
	UserName     string
	HashPassword string
	Access       map[string]bool
}

func GenHashPassword(Password string) string {
	b := ([]byte)(salt + Password)
	b1 := sha1.Sum(b)
	return hex.EncodeToString(b1[:])
}
func (u *User) GetUserName() string { return u.UserName }
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

var _u IUser = (IUser)(&User{})
