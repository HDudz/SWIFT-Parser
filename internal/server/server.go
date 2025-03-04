package server

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func ConnectDB() *sql.DB {
	dsn := "root:rootpswd@tcp(swift-db:3306)/swiftcodes?parseTime=true"
	var db *sql.DB
	var err error

	for i := 0; i < 15; i++ {
		db, err = sql.Open("mysql", dsn)
		if err == nil && db.Ping() == nil {
			fmt.Println("Connected to database successfully!")
			return db
		}
		fmt.Println("Waiting for database to come up...")
		time.Sleep(2 * time.Second)
	}

	log.Fatal("Failed to connect with database: ", err, " Ping: ", db.Ping())
	return nil
}
