package server

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

func ImportData(db *sql.DB) {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM swiftTable").Scan(&count)
	if err != nil {
		log.Fatal("Failed to check database:", err)
	}
	if count == 0 {
		err = database.ImportFromCSV(db, "data/SWIFT_Code.csv")
		if err != nil {
			log.Fatal("Failed to import data from CSV:", err)
		}
		fmt.Println("Data has been imported")
	} else {
		fmt.Println("Data already exists, skipping import")
	}
}

func StartServer() {
	r := chi.NewRouter()
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	port := "8080"
	fmt.Printf("Server running on port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
