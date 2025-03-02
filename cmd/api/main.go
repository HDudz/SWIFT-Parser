package main

import (
	"fmt"
	"github.com/HDudz/SWIFT-Parser/internal/server"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
)

func main() {

	r := chi.NewRouter()

	db := server.ConnectDB()

	server.ImportDataIfNeeded(db)

	r.Get("/swift/{code}", server.GetCodeHandler(db))

	port := "8080"
	fmt.Printf("Server running on port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
