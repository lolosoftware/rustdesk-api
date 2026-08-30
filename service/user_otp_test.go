package service

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/lejianwen/rustdesk-api/v2/config"
	"github.com/lejianwen/rustdesk-api/v2/model"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestGenerateTOTPCode(t *testing.T) {
	secret := "JBSWY3DPEHPK3PXP"
	code := GenerateTOTPCode(secret, time.Unix(0, 0))
	if code == "" {
		t.Fatal("GenerateTOTPCode returned empty code")
	}
	if code != "282760" {
		t.Fatalf("GenerateTOTPCode secret=%s at unix 0 = %q, want %q", secret, code, "282760")
	}
}

func TestVerifyOTPCode(t *testing.T) {
	secret := "JBSWY3DPEHPK3PXP"
	code := GenerateTOTPCode(secret, time.Now())
	if !VerifyOTPCode(secret, code) {
		t.Fatalf("VerifyOTPCode(secret, %q) = false, want true", code)
	}
	if VerifyOTPCode(secret, fmt.Sprintf("%06d", (100000+1)%1000000)) {
		t.Fatal("VerifyOTPCode should reject invalid code")
	}
}

func TestUserOTPEnrollmentRequiresValidCode(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	if err = db.AutoMigrate(&model.User{}); err != nil {
		t.Fatal(err)
	}
	DB = db
	Config = &config.Config{Admin: config.Admin{Title: "RustDesk Test"}}

	isAdmin := false
	u := &model.User{Username: "otp-user", IsAdmin: &isAdmin}
	if err = db.Create(u).Error; err != nil {
		t.Fatal(err)
	}

	users := &UserService{}
	enrollment, err := users.BeginUserOTPEnrollment(u)
	if err != nil {
		t.Fatal(err)
	}
	if enrollment.Secret == "" || enrollment.ProvisioningURI == "" || enrollment.QRCode == "" {
		t.Fatalf("incomplete enrollment: %#v", enrollment)
	}

	pending := users.InfoById(u.Id)
	if pending.OtpEnabled {
		t.Fatal("OTP must remain disabled before code confirmation")
	}
	validCode := GenerateTOTPCode(enrollment.Secret, time.Now())
	invalidCode := "0" + validCode[1:]
	if validCode[0] == '0' {
		invalidCode = "1" + validCode[1:]
	}
	if err = users.ConfirmUserOTP(pending, invalidCode); !errors.Is(err, ErrOTPInvalidCode) {
		t.Fatalf("invalid confirmation code returned %v, want ErrOTPInvalidCode", err)
	}

	if err = users.ConfirmUserOTP(pending, validCode); err != nil {
		t.Fatal(err)
	}
	enabled := users.InfoById(u.Id)
	if !enabled.OtpEnabled {
		t.Fatal("OTP was not enabled after a valid confirmation code")
	}
	if _, err = users.BeginUserOTPEnrollment(enabled); !errors.Is(err, ErrOTPAlreadyEnabled) {
		t.Fatalf("duplicate enrollment returned %v, want ErrOTPAlreadyEnabled", err)
	}
	if err = users.DisableUserOTP(enabled, validCode); err != nil {
		t.Fatal(err)
	}
	disabled := users.InfoById(u.Id)
	if disabled.OtpEnabled || disabled.OtpSecret != "" {
		t.Fatalf("OTP credentials were not cleared: enabled=%v secret=%q", disabled.OtpEnabled, disabled.OtpSecret)
	}
}
