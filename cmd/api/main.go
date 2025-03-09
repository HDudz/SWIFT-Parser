package main

import (
	"fmt"
	"github.com/HDudz/SWIFT-Parser/internal/services"
)

func main() {

	db := services.LoadDB()
	router := services.LoadRoutes(db)

	err := services.StartServer(router, db)

	if err != nil {
		fmt.Println("failed to start server: ", err)
	}

}
