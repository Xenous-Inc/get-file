package main

import (
	"fmt"
	"hui/internal/server"

	_ "github.com/joho/godotenv/autoload"
)

func main() {

	server := server.New()

	server.RegisterFiberRoutes()
	err := server.Listen(fmt.Sprintf(":%d", 8080))
	if err != nil {
		panic(fmt.Sprintf("cannot start server: %s", err))
	}
}
