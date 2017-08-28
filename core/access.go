package core

import (
	"github.com/crazyprograms/sud/corebase"
)

type AccessInfo struct {
	Name  string
	Title string
}

func (ai *AccessInfo) GetName() string  { return ai.Name }
func (ai *AccessInfo) GetTitle() string { return ai.Title }

type accessUnion struct {
	Access []corebase.IAccess
}

// AccessUnion - Пересечение всех доступов
func AccessUnion(Access ...corebase.IAccess) corebase.IAccess {
	return &accessUnion{Access: Access}
}
func (au *accessUnion) CheckAccess(Access string) bool {
	for i := 0; i < len(au.Access); i++ {
		if !au.Access[i].CheckAccess(Access) {
			return false
		}
	}
	return true
}
func (au *accessUnion) Users() []corebase.IUser {
	list := make([]corebase.IUser, 0, 10)
	for i := 0; i < len(au.Access); i++ {
		list = append(list, au.Access[i].Users()...)
	}
	return list
}
