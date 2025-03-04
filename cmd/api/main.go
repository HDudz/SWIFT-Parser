package main

import (
	"fmt"
	"github.com/HDudz/SWIFT-Parser/internal/database"
	"github.com/HDudz/SWIFT-Parser/internal/server"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
)

func main() {

	r := chi.NewRouter()

	db := server.ConnectDB()

	database.ImportDataIfNeeded(db)

	server.SetupRoutes(r, db)

	port := "8080"
	fmt.Printf("Server running on port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
