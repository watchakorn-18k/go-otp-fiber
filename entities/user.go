package entities

type VerifyRequest struct {
	OTP      string `json:"otp"`
	Username string `json:"username"`
}

type User struct {
	Username string `bson:"username"`
	Secret   string `bson:"secret"`
	URL      string `bson:"url"`
	QRCode   string `bson:"qrcode"`
}
