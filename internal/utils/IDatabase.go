package utils

import "database/sql"

type DB interface {
	Begin() (*sql.Tx, error)
	Execute(string, ...any) (sql.Result, error)
	Query(string, ...any) (*sql.Rows, error)
	QueryRow(string, ...any) *sql.Row
	Close() error
}
