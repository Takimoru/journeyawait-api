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
	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/go")
	if err != nil {
		panic(err)
	}
	defer db.Close()
}
