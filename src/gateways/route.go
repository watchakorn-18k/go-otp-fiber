package gateways

import "github.com/gofiber/fiber/v2"

func RouteOTP(gateway HTTPGateway, app *fiber.App) {
	api := app.Group("/api/otp")

	api.Get("/generate_link/:username", gateway.GenerateOTP)
	api.Post("/verify_otp", gateway.VerifyOTP)
}
