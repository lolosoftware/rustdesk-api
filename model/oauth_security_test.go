package model

import "testing"

func TestGoogleAndOIDCForcePKCES256(t *testing.T) {
	for _, oauthType := range []string{OauthTypeGoogle, OauthTypeOidc} {
		enabled := false
		oauth := &Oauth{
			OauthType:  oauthType,
			Op:         oauthType,
			PkceEnable: &enabled,
			PkceMethod: PKCEMethodPlain,
		}
		if err := oauth.FormatOauthInfo(); err != nil {
			t.Fatal(err)
		}
		if oauth.PkceEnable == nil || !*oauth.PkceEnable || oauth.PkceMethod != PKCEMethodS256 {
			t.Fatalf("%s did not force PKCE/S256: %#v", oauthType, oauth)
		}
		if oauth.AutoRegister == nil {
			t.Fatalf("%s did not initialize the auto-register policy", oauthType)
		}
	}
}
