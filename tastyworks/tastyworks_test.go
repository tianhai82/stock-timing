package tastyworks

import (
	"testing"
)

func Test_login(t *testing.T) {
	token, err := login()
	if err != nil {
		t.Fatalf("fail to login: %v", err)
	}
	err = saveToken(token)
	if err != nil {
		t.Fatalf("fail to save token: %v", err)
	}
	cachedToken, err := getCachedSessionToken()
	if err != nil {
		t.Fatalf("fail to get cached token: %v", err)
	}
	if cachedToken != token {
		t.Fatalf("cached token is different from original token")
	}
	ok, err := validateSessionToken(cachedToken)
	if err != nil {
		t.Fatalf("fail to validate token: %v", err)
	}
	if !ok {
		t.Fatalf("token validation failed")
	}
}
