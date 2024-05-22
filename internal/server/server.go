package server

import (
	"github.com/gofiber/fiber/v2"
)

type FiberServer struct {
	*fiber.App
}

func New() *FiberServer {
	server := &FiberServer{
		App: fiber.New(fiber.Config{
			ServerHeader: "getfile",
			AppName:      "getfile",
		}),
	}
	server.Server().LogAllErrors = true

	return server
}
