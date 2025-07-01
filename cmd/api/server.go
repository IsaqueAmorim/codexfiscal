package main

import (
	"github.com/IsaqueAmorim/codexfiscal/internal/config"
	"github.com/IsaqueAmorim/codexfiscal/internal/handler"
	"github.com/IsaqueAmorim/codexfiscal/internal/repository"
	"github.com/IsaqueAmorim/codexfiscal/internal/service"
	"github.com/gofiber/fiber/v2"
)

func main() {

	db, err := config.ConnectDatabase()
	if err != nil {
		panic("Failed to connect to the database")
	}
	defer db.Close()

	err = config.NewMigration(db).Run()
	if err != nil {
		panic("Failed to run migrations: " + err.Error())
	}

	// Start the server
	api := fiber.New(fiber.Config{
		Prefork:               false,
		DisableStartupMessage: true,
		BodyLimit:             1024 * 1024 * 100000, // 10 MB
	})

	handler := handler.NewNCMHandler(
		service.NewNCMService(
			repository.NewNCMRepository(db),
		))
	handler.RegisterNCMRoutes(api)
	if err := api.Listen(":8080"); err != nil {
		panic("Failed to start server: " + err.Error())
	}

	println("Server is running on port 8080")
}
