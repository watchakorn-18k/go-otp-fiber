# Go OTP Fiber Application

<div align="center">

[api-version](https://github.com/watchakorn-18k/go-otp-fiber/tree/master)/[htmx-version](https://github.com/watchakorn-18k/go-otp-fiber/tree/show-html-version)

</div>

This project is a simple OTP (One Time Password) generation and verification service built using the Go programming language and the Fiber web framework. It uses TOTP (Time-based One-Time Password) for generating and verifying OTPs.

## Features

- Generate a TOTP key for a user and return the QR code image.
- Verify the TOTP code provided by the user.

## Prerequisites

- Go 1.22+
- MongoDB
- Fiber framework

## Getting Started

### Installation

1. Clone the repository:

    ```sh
    git clone https://github.com/watchakorn-18k/go-otp-fiber
    cd go-otp-fiber
    ```

2. Install dependencies:

    ```sh
    go mod tidy
    ```

3. Set up your MongoDB database and update the connection string in the `domain` package.

### Running the Application

1. Start the application:

    ```sh
    go run main.go
    ```

2. The server will start on port `3000`.

### API Endpoints

#### Generate TOTP Key

- **Endpoint:** `GET /api/otp/generate_link/:username`
- **Description:** Generates a TOTP key for the specified username and returns the secret, URL, and QR code image.
- **Parameters:**
    - `username`: The username for which to generate the TOTP key.

- **Response:**

    ```json
    {
        "secret": "generated-secret",
        "url": "otpauth-url",
        "qrcode": "base64-encoded-qrcode"
    }
    ```

#### Verify TOTP Code

- **Endpoint:** `POST /api/otp/verify_otp`
- **Description:** Verifies the TOTP code provided by the user.
- **Request Body:**

    ```json
    {
        "username": "user123",
        "otp": "123456"
    }
    ```

- **Response:**

    - Success:

        ```json
        {
            "message": "valid OTP"
        }
        ```

    - Failure:

        ```json
        {
            "message": "invalid OTP"
        }
        ```

### Dependencies

- [Fiber](https://github.com/gofiber/fiber): An Express-inspired web framework for Go.
- [MongoDB Go Driver](https://github.com/mongodb/mongo-go-driver): The official MongoDB driver for Go.
- [pquerna/otp](https://github.com/pquerna/otp): A Go library for generating and verifying TOTP/HOTP codes.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.


