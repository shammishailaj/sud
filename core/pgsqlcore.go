package core

import (
	"database/sql"

	"github.com/crazyprograms/sud/corebase"
)

func getUsers(connection *sql.DB) (map[string]corebase.IUser, error) {
	var err error
	var rows1 *sql.Rows
	if rows1, err = connection.Query(`SELECT "User"."Login", "User"."PasswordHash" FROM "User" `); err != nil {
		return nil, err
	}
	defer rows1.Close()
	Users := map[string]*User{}
	Result := map[string]corebase.IUser{}
	for rows1.Next() {
		var Login, PasswordHash string
		if err = rows1.Scan(&Login, &PasswordHash); err != nil {
			return nil, err
		}
		var user = &User{Login: Login, HashPassword: PasswordHash, Access: map[string]bool{}}
		Users[Login] = user
		Result[Login] = user
	}
	var rows2 *sql.Rows
	if rows2, err = connection.Query(`SELECT "Access"."Login", "Access"."Access" FROM "Access" `); err != nil {
		return nil, err
	}
	defer rows2.Close()
	for rows2.Next() {
		var Login, Access string
		if err = rows2.Scan(&Login, &Access); err != nil {
			return nil, err
		}
		Users[Login].Access[Access] = true
	}
	return Result, nil
}
