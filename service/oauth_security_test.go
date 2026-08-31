package service

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/lejianwen/rustdesk-api/v2/config"
	"github.com/lejianwen/rustdesk-api/v2/model"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
)

func TestOauthCallbackCanOnlyBeConsumedOnce(t *testing.T) {
	const state = "single-use-oauth-state"
	users := &OauthService{}
	users.DeleteOauthCache(state)
	t.Cleanup(func() { users.DeleteOauthCache(state) })

	item := &OauthCacheItem{Op: model.OauthTypeGoogle}
	users.SetOauthCache(state, item, 0)

	if got := users.ConsumeOauthCallback(state); got != item {
		t.Fatal("first callback did not consume the pending transaction")
	}
	if got := users.ConsumeOauthCallback(state); got != nil {
		t.Fatal("OAuth callback transaction was accepted more than once")
	}
}

func TestOIDCCallbackRequiresIDToken(t *testing.T) {
	Config = &config.Config{}
	Logger = logrus.New()
	tokenServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"access_token":"access-token","token_type":"Bearer"}`))
	}))
	defer tokenServer.Close()

	oauthConfig := &oauth2.Config{
		ClientID:     "client-id",
		ClientSecret: "client-secret",
		Endpoint: oauth2.Endpoint{
			TokenURL: tokenServer.URL,
		},
	}
	users := &OauthService{}
	err, _ := users.callbackBase(oauthConfig, nil, "authorization-code", "", "nonce", true, &model.OidcUser{})
	if err == nil || err.Error() != "IDTokenMissing" {
		t.Fatalf("callbackBase returned %v, want IDTokenMissing", err)
	}
}

func TestOauthRegistrationRejectsUnverifiedEmail(t *testing.T) {
	users := &UserService{}
	err, user := users.RegisterByOauth(&model.OauthUser{
		OpenId:        "google-subject",
		Email:         "user@example.com",
		VerifiedEmail: false,
	}, model.OauthTypeGoogle)
	if err == nil || err.Error() != "OauthEmailNotVerified" || user != nil {
		t.Fatalf("RegisterByOauth returned err=%v user=%#v", err, user)
	}
}
