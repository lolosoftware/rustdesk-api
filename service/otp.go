package service

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha1"
	"encoding/base32"
	"fmt"
	"strings"
	"time"

	"github.com/lejianwen/rustdesk-api/v2/model"
)

const otpWindow = 1

func GenerateOTPSecret() string {
	b := make([]byte, 20)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	return strings.TrimSpace(base32.StdEncoding.EncodeToString(b))
}

func GenerateTOTPCode(secret string, t time.Time) string {
	secret = strings.TrimSpace(secret)
	if secret == "" {
		return ""
	}
	key, err := base32.StdEncoding.DecodeString(strings.ToUpper(secret))
	if err != nil || len(key) == 0 {
		return ""
	}

	counter := uint64(t.Unix() / 30)
	msg := make([]byte, 8)
	for i := 7; i >= 0; i-- {
		msg[i] = byte(counter & 0xff)
		counter >>= 8
	}

	hm := hmac.New(sha1.New, key)
	_, _ = hm.Write(msg)
	sum := hm.Sum(nil)
	offset := int(sum[len(sum)-1] & 0x0f)
	value := (int(sum[offset]&0x7f) << 24) |
		(int(sum[offset+1]&0xff) << 16) |
		(int(sum[offset+2]&0xff) << 8) |
		(int(sum[offset+3] & 0xff))

	code := value % 1000000
	return fmt.Sprintf("%06d", code)
}

func VerifyOTPCode(secret, code string) bool {
	secret = strings.TrimSpace(secret)
	code = strings.TrimSpace(code)
	if secret == "" || code == "" {
		return false
	}
	if len(code) != 6 {
		return false
	}
	for i := -otpWindow; i <= otpWindow; i++ {
		if GenerateTOTPCode(secret, time.Now().Add(time.Duration(i)*30*time.Second)) == code {
			return true
		}
	}
	return false
}

func (us *UserService) GenerateUserOTPSecret() string {
	return GenerateOTPSecret()
}

func (us *UserService) VerifyUserOTP(u *model.User, code string) bool {
	if u == nil || !u.OtpEnabled || strings.TrimSpace(u.OtpSecret) == "" {
		return false
	}
	return VerifyOTPCode(u.OtpSecret, code)
}
