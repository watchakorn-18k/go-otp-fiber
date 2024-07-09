package services

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"go-opt-fiber/src/domain/entities"
	"go-opt-fiber/src/domain/repositories"
	"image/png"
	"time"

	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
	"go.mongodb.org/mongo-driver/mongo"
)

type OTPService struct {
	UsersRepository repositories.IUsersRepository
	expireTime      uint
}

type IOTPService interface {
	isValidOTP(otpCode, secret string) bool
	GenerateOTP(username string) (*entities.OTP, error)
	VerifyOTP(data *entities.VerifyRequest) error
}

func NewOTPService(repo0 repositories.IUsersRepository) IOTPService {
	return &OTPService{
		UsersRepository: repo0,
		expireTime:      30,
	}
}

func (s *OTPService) isValidOTP(otpCode, secret string) bool {
	valid, _ := totp.ValidateCustom(
		otpCode,
		secret,
		time.Now().UTC(),
		totp.ValidateOpts{
			Period:    s.expireTime,
			Skew:      1,
			Digits:    otp.DigitsSix,
			Algorithm: otp.AlgorithmSHA1,
		},
	)
	return valid
}

func (s *OTPService) GenerateOTP(username string) (*entities.OTP, error) {
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "WK-18K Server",
		AccountName: username,
		Period:      s.expireTime,
		SecretSize:  10,
	})
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	img, err := key.Image(200, 200)
	if err != nil {
		return nil, err
	}
	png.Encode(&buf, img)

	if err := s.UsersRepository.SaveOrUpdateUser(username, key.Secret(), key.URL()); err != nil {
		return nil, err
	}

	return &entities.OTP{
		Secret: key.Secret(),
		URL:    key.URL(),
		QRCode: fmt.Sprintf("data:image/jpeg;base64,%s", base64.StdEncoding.EncodeToString(buf.Bytes())),
	}, nil

}

func (s *OTPService) VerifyOTP(data *entities.VerifyRequest) error {
	user, err := s.UsersRepository.FindUserByUsername(data.Username)
	if err != nil && err != mongo.ErrNoDocuments {
		return err
	}
	if err == mongo.ErrNoDocuments {
		return errors.New("username not found")
	}
	if s.isValidOTP(data.OTP, user.Secret) {
		return nil
	} else {
		return errors.New("invalid OTP")
	}
}
