package database

import (
	"log"

	"github.com/jmoiron/sqlx"
	// Provide mysql driver for sql.Open
	_ "github.com/go-sql-driver/mysql"
)

// Connect - Creates a database handle instance using the given DSN
func Connect(dsn string) *sqlx.DB {
	// sqlx's Connect handles both opening the connection and Ping
	dbHandle, err := sqlx.Connect("mysql", dsn)
	if err != nil {
		log.Println(err)
		return nil
	}

	return dbHandle
}
