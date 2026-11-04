package services

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/url"
	"os"
)

// TokenResponse maps the payload returned by Google's token endpoint
type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token,omitempty"`
	Scope        string `json:"scope"`
	TokenType    string `json:"token_type"`
	IDToken      string `json:"id_token"`
}

// GenerateStateToken creates a cryptographically secure random state for CSRF protection
func GenerateStateToken() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", fmt.Errorf("failed to generate state token: %w", err)
	}
	return hex.EncodeToString(b), nil
}

// RequestGoogleAuthURL builds the Google OAuth2 consent screen redirect URL
func RequestGoogleAuthURL(state string) string {
	clientID := os.Getenv("GOOGLE_CLIENT_ID")
	urlRedirect := os.Getenv("GOOGLE_REDIRECT_URL")

	params := url.Values{}
	params.Set("client_id", clientID)
	params.Set("redirect_uri", urlRedirect)
	params.Set("response_type", "code")
	params.Set("scope", "openid email profile")
	params.Set("state", state)
	params.Set("access_type", "offline")
	params.Set("prompt", "consent")

	return "https://accounts.google.com/o/oauth2/v2/auth?" + params.Encode()
}
