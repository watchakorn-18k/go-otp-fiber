package main

import (
	"bytes"
	"go-opt-fiber/domain"
	"go-opt-fiber/entities"
	"image/png"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
	"go.mongodb.org/mongo-driver/mongo"
)

var userCollection *mongo.Collection

func main() {
	app := fiber.New()
	app.Use(logger.New(logger.Config{
		Format: "[${ip}]:${port} ${time} ${pid} ${status} - ${method} ${path}: ${latency}\n",
	}))
	collDB, err := domain.ConnectToMongoDB()
	if err != nil {
		log.Fatal(err)
	}
	userCollection = collDB
	app.Get("/generate_link/:username", generateLinkHandler)
	app.Post("/verify_otp", verifyOTPHandler)

	log.Fatal(app.Listen(":3000"))
}

func generateLinkHandler(c *fiber.Ctx) error {
	username := c.Params("username")
	expireTime := uint(30)

	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "WK-18K Server",
		AccountName: username,
		Period:      expireTime,
	})

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error generating key",
		})
	}

	var buf bytes.Buffer
	img, err := key.Image(200, 200)
	png.Encode(&buf, img)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error generating key",
		})
	}

	err = domain.SaveOrUpdateUser(username, key.Secret(), key.URL(), buf.Bytes(), userCollection)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error saving key",
		})
	}

	return c.JSON(fiber.Map{
		"secret": key.Secret(),
		"url":    key.URL(),
	})
}

func verifyOTPHandler(c *fiber.Ctx) error {
	var req entities.VerifyRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
		})
	}

	user, err := domain.FindUserByUsername(req.Username, userCollection)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "User not found",
		})
	}

	if isValidOTP(req.OTP, user.Secret) {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "OTP is valid",
		})
	}

	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
		"message": "Invalid OTP",
	})
}

func isValidOTP(otpCode, secret string) bool {
	expireTime := uint(30)
	valid, _ := totp.ValidateCustom(
		otpCode,
		secret,
		time.Now().UTC(),
		totp.ValidateOpts{
			Period:    expireTime,
			Skew:      1,
			Digits:    otp.DigitsSix,
			Algorithm: otp.AlgorithmSHA1,
		},
	)
	return valid
}
