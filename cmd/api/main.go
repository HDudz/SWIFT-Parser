package main

import "github.com/HDudz/SWIFT-Parser/internal/server"

func main() {
	db := server.ConnectDB()
	server.ImportData(db)
	server.StartServer()
}
