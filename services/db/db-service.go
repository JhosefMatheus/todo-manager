package dbservice

import (
	"database/sql"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

func GetDbConnection() (*sql.DB, error) {
	dbms := os.Getenv("DBMS")
	dbConnection := os.Getenv("DB_CONNECTION")

	db, err := sql.Open(dbms, dbConnection)

	if err != nil {
		return nil, err
	}

	return db, nil
}

func CloseDbConnection(db *sql.DB) error {
	err := db.Close()

	return err
}
