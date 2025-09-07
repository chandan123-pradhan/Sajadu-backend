package config

import (
	"database/sql"
	"log"
	"os"
    _ "github.com/go-sql-driver/mysql"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func InitDB() {
	var err error
	connStr := os.Getenv("DB_CONN_STR")
	DB, err = sql.Open("mysql", connStr)
	if err != nil {
		log.Fatal(err)
	}
	if err = DB.Ping(); err != nil {
		log.Fatal("Cannot connect to DB: ", err)
	}
	log.Println("Connected to the database")
}
