package golang_database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"testing"

	// Make sure to replace the driver import with the one that matches your database
	_ "github.com/go-sql-driver/mysql" // This is an example for MySQL; replace with your specific driver if different
)

// getConnection function to open a new database connection
func getConnections() *sql.DB {
	dbname := "root:1234@tcp(localhost:3306)/golangdb"

	db, err := sql.Open("mysql", dbname)
	if err != nil {
		log.Fatalf("failed open database: %v", err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatalf("failed to ping database: %v", err)
	}
	return db
}

func TestExecsql(t *testing.T) {
	db := getConnections()
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Fatalf("Failed to close database connection: %v", err)
		}
	}(db)

	ctx := context.Background()

	script := "INSERT INTO customer(ID, NAME) VALUES(3, 'dini')"
	_, err := db.ExecContext(ctx, script)
	if err != nil {
		panic(fmt.Sprintf("Failed to execute SQL script: %v", err))
	}

	fmt.Println("SQL script executed successfully")
}
