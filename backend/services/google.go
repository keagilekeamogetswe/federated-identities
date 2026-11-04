package services

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"

	"google.golang.org/api/idtoken"
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

// DirectServerConfirmation executes the back-channel POST request to exchange
// the authorization code for access and ID tokens.
func DirectServerConfirmation(code string) (*TokenResponse, error) {
	clientID := os.Getenv("GOOGLE_CLIENT_ID")
	clientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")
	urlRedirect := os.Getenv("GOOGLE_REDIRECT_URL")

	formData := url.Values{}
	formData.Set("code", code)
	formData.Set("client_id", clientID)
	formData.Set("client_secret", clientSecret)
	formData.Set("redirect_uri", urlRedirect)
	formData.Set("grant_type", "authorization_code")

	resp, err := http.PostForm("https://oauth2.googleapis.com/token", formData)
	if err != nil {
		return nil, fmt.Errorf("failed to send token request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("google token endpoint returned status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var tokenResp TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return nil, fmt.Errorf("failed to decode token response: %w", err)
	}

	return &tokenResp, nil
}

// ConfirmGoogleID cryptographically verifies Google's JWT ID token and extracts user claims
func ConfirmGoogleID(ctx context.Context, idTokenString string) (*idtoken.Payload, error) {
	clientID := os.Getenv("GOOGLE_CLIENT_ID")

	// Validates signature, audience (client ID), expiration, and issuer via Google's public keys
	payload, err := idtoken.Validate(ctx, idTokenString, clientID)
	if err != nil {
		return nil, fmt.Errorf("invalid google id token: %w", err)
	}

	return payload, nil
}