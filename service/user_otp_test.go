package service

import (
	"fmt"
	"testing"
	"time"
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
