package db

import (
	"database/sql"
	"log"

	"github.com/go-sql-driver/mysql"
)

func NewMysqlDatabase(configuration mysql.Config) (*sql.DB, error) {
	db, err := sql.Open("mysql", configuration.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	return db, nil
}
