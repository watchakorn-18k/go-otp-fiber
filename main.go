package main

import (
	"bytes"
	"encoding/base64"
	"go-opt-fiber/domain"
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

	// เสิร์ฟไฟล์ static จากโฟลเดอร์ public
	app.Static("/", "./public")

	app.Post("/generate_link", generateLinkHandler)
	app.Post("/verify_otp", verifyOTPHandler)

	log.Fatal(app.Listen(":3000"))
}

func generateLinkHandler(c *fiber.Ctx) error {
	username := c.FormValue("username")
	expireTime := uint(30)

	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "WK-18K Server",
		AccountName: username,
		Period:      expireTime,
	})

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Error generating key")
	}

	var buf bytes.Buffer
	img, err := key.Image(200, 200)
	png.Encode(&buf, img)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Error generating key")
	}

	err = domain.SaveOrUpdateUser(username, key.Secret(), key.URL(), buf.Bytes(), userCollection)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Error saving key")
	}

	return c.Status(fiber.StatusOK).SendString("<img src='data:image/png;base64," + encodeImageToBase64(buf.Bytes()) + "' alt='QR Code'>")
}

func verifyOTPHandler(c *fiber.Ctx) error {
	username := c.FormValue("username")
	otpCode := c.FormValue("otp")

	user, err := domain.FindUserByUsername(username, userCollection)
	if err != nil {
		return c.Status(fiber.StatusNotFound).SendString("<p style='color: red;'>User not found</p>")
	}

	if isValidOTP(otpCode, user.Secret) {
		return c.SendString("<p style='color: green;'>OTP is valid</p>")
	}

	return c.SendString("<p style='color: red;'>Invalid OTP</p>")
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

func encodeImageToBase64(image []byte) string {
	return base64.StdEncoding.EncodeToString(image)
}
