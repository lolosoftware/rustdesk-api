package service

import (
	"bytes"
	"encoding/base64"
	"errors"
	"image/png"
	"strings"
	"time"

	"github.com/lejianwen/rustdesk-api/v2/model"
	"github.com/pquerna/otp/totp"
)

var (
	ErrOTPAlreadyEnabled = errors.New("OTP is already enabled")
	ErrOTPNotConfigured  = errors.New("OTP is not configured")
	ErrOTPInvalidCode    = errors.New("invalid OTP code")
)

type OTPEnrollment struct {
	Secret          string `json:"secret"`
	ProvisioningURI string `json:"provisioning_uri"`
	QRCode          string `json:"qr_code"`
}

func GenerateOTPSecret() string {
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "RustDesk API",
		AccountName: "test",
		SecretSize:  20,
	})
	if err != nil {
		return ""
	}
	return key.Secret()
}

func GenerateTOTPCode(secret string, t time.Time) string {
	code, err := totp.GenerateCode(strings.TrimSpace(secret), t)
	if err != nil {
		return ""
	}
	return code
}

func VerifyOTPCode(secret, code string) bool {
	secret = strings.TrimSpace(secret)
	code = strings.TrimSpace(code)
	if secret == "" || len(code) != 6 {
		return false
	}
	for _, digit := range code {
		if digit < '0' || digit > '9' {
			return false
		}
	}
	return totp.Validate(code, secret)
}

func (us *UserService) BeginUserOTPEnrollment(u *model.User) (*OTPEnrollment, error) {
	if u == nil || u.Id == 0 {
		return nil, ErrOTPNotConfigured
	}
	if u.OtpEnabled {
		return nil, ErrOTPAlreadyEnabled
	}

	issuer := "RustDesk API"
	if Config != nil && strings.TrimSpace(Config.Admin.Title) != "" {
		issuer = strings.TrimSpace(Config.Admin.Title)
	}
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      issuer,
		AccountName: u.Username,
		SecretSize:  20,
	})
	if err != nil {
		return nil, err
	}
	image, err := key.Image(256, 256)
	if err != nil {
		return nil, err
	}
	var pngData bytes.Buffer
	if err = png.Encode(&pngData, image); err != nil {
		return nil, err
	}

	if err = DB.Model(&model.User{}).Where("id = ?", u.Id).Updates(map[string]interface{}{
		"otp_enabled": false,
		"otp_secret":  key.Secret(),
	}).Error; err != nil {
		return nil, err
	}

	return &OTPEnrollment{
		Secret:          key.Secret(),
		ProvisioningURI: key.URL(),
		QRCode:          "data:image/png;base64," + base64.StdEncoding.EncodeToString(pngData.Bytes()),
	}, nil
}

func (us *UserService) ConfirmUserOTP(u *model.User, code string) error {
	if u == nil || u.Id == 0 || u.OtpEnabled || strings.TrimSpace(u.OtpSecret) == "" {
		return ErrOTPNotConfigured
	}
	if !VerifyOTPCode(u.OtpSecret, code) {
		return ErrOTPInvalidCode
	}
	return DB.Model(&model.User{}).Where("id = ?", u.Id).Update("otp_enabled", true).Error
}

func (us *UserService) DisableUserOTP(u *model.User, code string) error {
	if u == nil || u.Id == 0 || !u.OtpEnabled || strings.TrimSpace(u.OtpSecret) == "" {
		return ErrOTPNotConfigured
	}
	if !VerifyOTPCode(u.OtpSecret, code) {
		return ErrOTPInvalidCode
	}
	return us.ResetUserOTP(u)
}

func (us *UserService) ResetUserOTP(u *model.User) error {
	if u == nil || u.Id == 0 {
		return ErrOTPNotConfigured
	}
	return DB.Model(&model.User{}).Where("id = ?", u.Id).Updates(map[string]interface{}{
		"otp_enabled": false,
		"otp_secret":  "",
	}).Error
}

func (us *UserService) VerifyUserOTP(u *model.User, code string) bool {
	if u == nil || !u.OtpEnabled || strings.TrimSpace(u.OtpSecret) == "" {
		return false
	}
	return VerifyOTPCode(u.OtpSecret, code)
}
