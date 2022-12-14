package db

import (
	"database/sql"
	"fmt"
	"strings"

	log "go-core/log"

	"go.uber.org/zap"
)

type oracleDBHelper struct {
	db *sql.DB
}

// NewDBHelper creates an instance
func NewOracleDBHelper(host string, port int, username, password, database string) DBHelper {
	db, err := initOracle(host, port, username, password, database)
	if err != nil {
		log.Logger.Panic("Failed to init oracle", zap.Error(err))
	}
	return &oracleDBHelper{
		db: db,
	}
}

func (h *oracleDBHelper) QueryRowsPaging(statement string, offset, limit uint32,
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

func (h *oracleDBHelper) QueryRowPaging(statement string, offset, limit uint32,
	agruments []interface{}) (row *sql.Row) {
	db := h.Open()
	queryBuilder := strings.Builder{}
	queryBuilder.WriteString(statement)
	queryBuilder.WriteString(fmt.Sprintf(" offset :%d rows fetch next :%d rows only", len(agruments)+1, len(agruments)+2))
	agruments = append(agruments, offset, limit)
	row = db.QueryRow(queryBuilder.String(), agruments...)
	return row
}

func (h *oracleDBHelper) Open() *sql.DB {
	return h.db
}

func (h *oracleDBHelper) Close() error {
	return h.db.Close()
}

func (h *oracleDBHelper) Begin() (*sql.Tx, error) {
	return h.db.Begin()
}

func (h *oracleDBHelper) Commit(tx *sql.Tx) error {
	return tx.Commit()
}

func (h *oracleDBHelper) Rollback(tx *sql.Tx) error {
	return tx.Rollback()
}

func initOracle(host string, port int, username, password, database string) (*sql.DB, error) {
	connectionString := fmt.Sprintf("%v/%v@%v:%v/%v", username, password, host, port, database)

	db, err := sql.Open("oci8", connectionString)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
