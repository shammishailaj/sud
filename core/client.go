package core

import (
	"database/sql"

	_ "github.com/lib/pq"
)

type IQuery interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	Prepare(query string) (*sql.Stmt, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
}
type transaction struct {
	core *Core
	tx   *sql.Tx
}

func (t *transaction) Commit() error {
	return t.tx.Commit()
}
func (t *transaction) Rollback() error {
	return t.tx.Rollback()
}
func (t *transaction) Exec(query string, args ...interface{}) (sql.Result, error) {
	return t.tx.Exec(query, args...)
}
func (t *transaction) Prepare(query string) (*sql.Stmt, error) {
	return t.tx.Prepare(query)
}
func (t *transaction) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return t.tx.Query(query, args...)
}
func (t *transaction) QueryRow(query string, args ...interface{}) *sql.Row {
	return t.tx.QueryRow(query, args...)
}

type Client struct {
	core              *Core
	user              IUser
	configurationName string
}

func (core *Core) NewClient(Login string, Password string, ConfigurationName string) *Client {
	user := core.getUser(Login)
	if !user.GetCheckPassword(Password) {
		return nil
	}
	//configuration := core.LoadConfiguration(ConfigurationName)
	return &Client{user: user, configurationName: ConfigurationName, core: core}
}

/**/
