package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lejianwen/rustdesk-api/v2/config"
	"github.com/lejianwen/rustdesk-api/v2/global"
	apiResp "github.com/lejianwen/rustdesk-api/v2/http/response/api"
	"github.com/lejianwen/rustdesk-api/v2/lib/jwt"
	"github.com/lejianwen/rustdesk-api/v2/lib/lock"
	"github.com/lejianwen/rustdesk-api/v2/model"
	"github.com/lejianwen/rustdesk-api/v2/service"
	"github.com/lejianwen/rustdesk-api/v2/utils"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/sirupsen/logrus"
	"golang.org/x/text/language"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestRustDeskOTPLoginFlow(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db, err := gorm.Open(sqlite.Open("file:api-otp-login?mode=memory&cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	if err = db.AutoMigrate(&model.User{}, &model.UserToken{}, &model.LoginLog{}); err != nil {
		t.Fatal(err)
	}

	cfg := &config.Config{}
	cfg.App.TokenExpire = 24 * time.Hour
	logger := logrus.New()
	service.New(cfg, db, logger, jwt.NewJwt("", cfg.App.TokenExpire), lock.NewLocal())
	global.Config = *cfg
	global.Logger = logger
	global.LoginLimiter = utils.NewLoginLimiter(utils.SecurityPolicy{CaptchaThreshold: -1})
	global.Validator.ValidStruct = func(*gin.Context, interface{}) []string { return nil }
	bundle := i18n.NewBundle(language.English)
	global.Localizer = func(string) *i18n.Localizer { return i18n.NewLocalizer(bundle, "en") }

	password, err := utils.EncryptPassword("password")
	if err != nil {
		t.Fatal(err)
	}
	u := &model.User{
		Username:   "otp-user",
		Password:   password,
		Status:     model.COMMON_STATUS_ENABLE,
		IsAdmin:    new(bool),
		OtpEnabled: true,
		OtpSecret:  "JBSWY3DPEHPK3PXP",
	}
	if err = db.Create(u).Error; err != nil {
		t.Fatal(err)
	}

	first := performLoginRequest(t, map[string]interface{}{
		"type":     apiResp.AuthRequestTypeAccount,
		"username": u.Username,
		"password": "password",
		"id":       "device-id",
		"uuid":     "device-uuid",
		"deviceInfo": map[string]string{
			"type": model.LoginLogClientApp,
			"os":   "Windows",
		},
	})
	if first.Code != http.StatusOK {
		t.Fatalf("password step status = %d, body = %s", first.Code, first.Body.String())
	}
	var challenge apiResp.LoginRes
	if err = json.Unmarshal(first.Body.Bytes(), &challenge); err != nil {
		t.Fatal(err)
	}
	if challenge.Type != apiResp.AuthResponseTypeEmailCheck ||
		challenge.TfaType != apiResp.AuthResponseTypeTfaCheck || challenge.Secret == "" {
		t.Fatalf("unexpected OTP challenge response: %#v", challenge)
	}
	var tokenCount int64
	if err = db.Model(&model.UserToken{}).Count(&tokenCount).Error; err != nil {
		t.Fatal(err)
	}
	if tokenCount != 0 {
		t.Fatal("access token issued before OTP verification")
	}

	secondBody := map[string]interface{}{
		// This is the request type used by the official RustDesk client for
		// both email and TOTP verification codes.
		"type":     apiResp.AuthRequestTypeEmailCode,
		"username": u.Username,
		"tfaCode":  service.GenerateTOTPCode(u.OtpSecret, time.Now()),
		"secret":   challenge.Secret,
		"id":       "device-id",
		"uuid":     "device-uuid",
		"deviceInfo": map[string]string{
			"type": model.LoginLogClientApp,
			"os":   "Windows",
		},
	}
	second := performLoginRequest(t, secondBody)
	if second.Code != http.StatusOK {
		t.Fatalf("OTP step status = %d, body = %s", second.Code, second.Body.String())
	}
	var login apiResp.LoginRes
	if err = json.Unmarshal(second.Body.Bytes(), &login); err != nil {
		t.Fatal(err)
	}
	if login.Type != apiResp.AuthResponseTypeToken || login.AccessToken == "" {
		t.Fatalf("unexpected login response: %#v", login)
	}

	replay := performLoginRequest(t, secondBody)
	if replay.Code != http.StatusBadRequest {
		t.Fatalf("consumed challenge replay status = %d, want %d", replay.Code, http.StatusBadRequest)
	}
}

func performLoginRequest(t *testing.T, body map[string]interface{}) *httptest.ResponseRecorder {
	t.Helper()
	payload, err := json.Marshal(body)
	if err != nil {
		t.Fatal(err)
	}
	recorder := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(recorder)
	context.Request = httptest.NewRequest(http.MethodPost, "/api/login", bytes.NewReader(payload))
	context.Request.Header.Set("Content-Type", "application/json")
	(&Login{}).Login(context)
	return recorder
}
