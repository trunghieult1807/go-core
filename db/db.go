package db

import "database/sql"

// DBHelper is helper of DB
type DBHelper interface {
	Open() *sql.DB
	Close() error
	Begin() (*sql.Tx, error)
	Commit(tx *sql.Tx) error
	Rollback(tx *sql.Tx) error
	QueryRowsPaging(statement string, offset, limit uint32, agruments []interface{}) (*sql.Rows, error)
	QueryRowPaging(statement string, offset, limit uint32, agruments []interface{}) *sql.Row
}
