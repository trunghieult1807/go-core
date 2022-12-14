package db

import (
	"database/sql"
	"fmt"
	"strings"

	log "go-core/log"

	"go.uber.org/zap"
)

type postgresDBHelper struct {
	db *sql.DB
}

// NewDBHelper creates an instance
func NewPostgresDBHelper(host string, port int, username, password, database string) DBHelper {
	db, err := initPostgres(host, port, username, password, database)
	if err != nil {
		log.Logger.Panic("Failed to init postgres", zap.Error(err))
	}
	return &postgresDBHelper{
		db: db,
	}
}

func (h *postgresDBHelper) QueryRowsPaging(statement string, offset, limit uint32,
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

func (h *postgresDBHelper) QueryRowPaging(statement string, offset, limit uint32,
	agruments []interface{}) (row *sql.Row) {
	db := h.Open()
	queryBuilder := strings.Builder{}
	queryBuilder.WriteString(statement)
	queryBuilder.WriteString(fmt.Sprintf(" offset :%d rows fetch next :%d rows only", len(agruments)+1, len(agruments)+2))
	agruments = append(agruments, offset, limit)
	row = db.QueryRow(queryBuilder.String(), agruments...)
	return row
}

func (h *postgresDBHelper) Open() *sql.DB {
	return h.db
}

func (h *postgresDBHelper) Close() error {
	return h.db.Close()
}

func (h *postgresDBHelper) Begin() (*sql.Tx, error) {
	return h.db.Begin()
}

func (h *postgresDBHelper) Commit(tx *sql.Tx) error {
	return tx.Commit()
}

func (h *postgresDBHelper) Rollback(tx *sql.Tx) error {
	return tx.Rollback()
}

func initPostgres(host string, port int, username, password, database string) (*sql.DB, error) {
	connectionString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, username, password, database)

	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
