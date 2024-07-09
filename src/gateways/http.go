package gateways

import (
	service "go-opt-fiber/src/services"

	"github.com/gofiber/fiber/v2"
)

type HTTPGateway struct {
	OTPServer service.IOTPService
}

func NewHTTPGateway(app *fiber.App, otp service.IOTPService) {
	gateway := &HTTPGateway{
		OTPServer: otp,
	}

	RouteOTP(*gateway, app)
}
