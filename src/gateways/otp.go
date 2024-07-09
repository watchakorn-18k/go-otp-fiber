package gateways

import (
	"go-opt-fiber/src/domain/entities"

	"github.com/gofiber/fiber/v2"
)

func (h *HTTPGateway) GenerateOTP(ctx *fiber.Ctx) error {
	username := ctx.Params("username")
	if username == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(entities.ResponseModel{
			Message: "query username is require",
		})
	}
	data, err := h.OTPServer.GenerateOTP(username)
	if err != nil {
		return ctx.Status(400).JSON(entities.ResponseMessage{Message: err.Error()})
	}
	return ctx.Status(201).JSON(entities.ResponseModel{
		Message: "Generate success",
		Data:    data,
	})
}

func (h *HTTPGateway) VerifyOTP(ctx *fiber.Ctx) error {
	data := new(entities.VerifyRequest)
	if err := ctx.BodyParser(data); err != nil {
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(entities.ResponseModel{
			Message: "Unprocessable Entity",
		})
	}
	if err := h.OTPServer.VerifyOTP(data); err != nil {
		return ctx.Status(400).JSON(entities.ResponseMessage{
			Message: err.Error(),
		})
	}
	return ctx.Status(200).JSON(entities.ResponseMessage{
		Message: "valid OTP",
	})
}
