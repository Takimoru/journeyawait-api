package golang_database

import (
	"database/sql"
	"time"
)

func getConnection() *sql.DB {
	db, err := sql.Open("mysql", "root:1234@tcp(localhost:3306)golangdb")
	if err != nil {
		panic(err)
	}

	db.SetMaxIdleConns(10)
	db.SetMaxOpenConns(100)
	db.SetConnMaxIdleTime(5 * time.Minute)
	db.SetConnMaxLifetime(60 * time.Minute)
	return db
}
