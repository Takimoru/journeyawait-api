package golang_database

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	//driver wajib di import agar terdetect tergistrasi databasenya jika tidak maka akan err
	"testing"
)

func TestEmpty(t *testing.T) {

}

func TestOpenConnection(t *testing.T) {
	db, err := sql.Open("mysql", "root:1234@tcp(localhost:3306)/golangdb")
	if err != nil {
		panic(err)
	}
	defer db.Close()
}
