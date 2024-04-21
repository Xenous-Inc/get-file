package main

import (
	"fmt"
	"hui/internal/server"

	_ "github.com/joho/godotenv/autoload"
)

func main() {

	server := server.New()

	server.RegisterFiberRoutes()
	// port, _ := strconv.Atoi(os.Getenv("PORT"))
	err := server.Listen(fmt.Sprintf(":%d", 8080))
	if err != nil {
		panic(fmt.Sprintf("cannot start server: %s", err))
	}
}
