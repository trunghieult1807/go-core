package db

import (
	"database/sql"
	"fmt"
	"strings"

	"go.uber.org/zap"
)

type sqlServerDBHelper struct {
	db *sql.DB
}

// NewSQLServerDBHelper creates an instance
func NewSQLServerDBHelper(host string, port int, username, password, database string) DBHelper {
	dbConnection, err := initSQLServer(host, port, username, password, database)
	if err != nil {
		zap.S().Panic("Failed to init SQL Server", zap.Error(err))
	}
	return &sqlServerDBHelper{
		db: dbConnection,
	}
}

func (h *sqlServerDBHelper) Open() *sql.DB {
	return h.db
}

func (h *sqlServerDBHelper) Close() error {
	return h.db.Close()
}

func (h *sqlServerDBHelper) Begin() (*sql.Tx, error) {
	return h.db.Begin()
}

func (h *sqlServerDBHelper) Commit(tx *sql.Tx) error {
	return tx.Commit()
}

func (h *sqlServerDBHelper) Rollback(tx *sql.Tx) error {
	return tx.Rollback()
}

func (h *sqlServerDBHelper) QueryRowPaging(statement string, offset, limit uint32,
	agruments []interface{}) (row *sql.Row) {
	db := h.Open()
	queryBuilder := strings.Builder{}
	queryBuilder.WriteString(statement)
	queryBuilder.WriteString(fmt.Sprintf(" offset :%d rows fetch next :%d rows only", len(agruments)+1, len(agruments)+2))
	agruments = append(agruments, offset, limit)
	return db.QueryRow(queryBuilder.String(), agruments...)
}

func (h *sqlServerDBHelper) QueryRowsPaging(statement string, offset, limit uint32,
	agruments []interface{}) (rows *sql.Rows, errQuery error) {
	db := h.Open()
	queryBuilder := strings.Builder{}
	queryBuilder.WriteString(statement)
	queryBuilder.WriteString(fmt.Sprintf(" offset :%d rows fetch next :%d rows only", len(agruments)+1, len(agruments)+2))
	agruments = append(agruments, offset, limit)
	rows, errQuery = db.Query(queryBuilder.String(), agruments...)
	if errQuery != nil {
		return nil, errQuery
	}
	return rows, nil
}

func initSQLServer(host string, port int, username, password, database string) (*sql.DB, error) {
	dbConnectionString := fmt.Sprintf("server=%s;port=%d;user id=%s;password=%s;database=%s;",
		host, port, username, password, database)
	dbConnection, err := sql.Open("sqlserver", dbConnectionString)
	if err != nil {
		return nil, err
	}
	if err := dbConnection.Ping(); err != nil {
		return nil, err
	}
	return dbConnection, nil
}
