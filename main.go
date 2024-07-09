package main

import (
	"fmt"
	"go-opt-fiber/src/domain/datasources"
	"go-opt-fiber/src/domain/repositories"
	"go-opt-fiber/src/gateways"
	"go-opt-fiber/src/middlewares"
	"go-opt-fiber/src/services"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}

	app := fiber.New()
	middlewares.Logger(app)
	mongoDB := datasources.NewMongoDB(10)
	userRepo := repositories.NewUsersRepository(mongoDB)
	sv1 := services.NewOTPService(userRepo)
	gateways.NewHTTPGateway(app, sv1)
	log.Fatal(app.Listen(":3000"))
}
