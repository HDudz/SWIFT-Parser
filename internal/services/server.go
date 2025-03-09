package services

import (
	"database/sql"
	"fmt"
	"github.com/HDudz/SWIFT-Parser/internal/database"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func StartServer(router *chi.Mux, db *sql.DB) error {

	server := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	database.ImportDataIfNeeded(db, "data/SWIFT_Code.csv")

	fmt.Printf("Server running on port %s\n", server.Addr)
	err := server.ListenAndServe()

	if err != nil {
		return fmt.Errorf("failed to listen and serve. error: %w", err)
	}

	return nil
}

func LoadDB() *sql.DB {
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
