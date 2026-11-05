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
)

// GitHubTokenResponse maps the payload returned by GitHub's token endpoint
type GitHubTokenResponse struct {
	AccessToken string `json:"access_token"`
	Scope string `json:"scope"`
	TokenType string `json:"token_type"`
}

// GitHubUser represents the public user profile returned by GitHub's API
type GitHubUser struct {
	ID int `json:"id"`
	Login string `json:"login"`
	Name string `json:"name"`
	Email string `json:"email"`
	AvatarURL string `json:"avatar_url"`
}

// GitHubEmail represents an email entry from GitHub's user emails endpoint
type GitHubEmail struct {
	Email string `json:"email"`
	Primary bool `json:"primary"`
	Verified bool `json:"verified"`
}

// GenerateGitHubStateToken creates a cryptographically secure random state for CSRF protection
func GenerateGitHubStateToken() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", fmt.Errorf("failed to generate github state token: %w", err)
	}
	return hex.EncodeToString(b), nil
}

// RequestGitHubAuthURL builds the GitHub OAuth consent screen redirect URL
func RequestGitHubAuthURL(state string) string {
	clientID := os.Getenv("CLIENT_ID_GITHUB")
	urlRedirect := os.Getenv("REDIRECT_URL_GITHUB")

	params := url.Values{}
	params.Set("client_id", clientID)
	params.Set("redirect_uri", urlRedirect)
	params.Set("scope", "read:user user:email") // Request user profile and verified emails
	params.Set("state", state)

	return "https://github.com/login/oauth/authorize?" + params.Encode()
}

// ExchangeGitHubCode executes the back-channel POST request to exchange
// the authorization code for an access token.
func ExchangeGitHubCode(code string) (*GitHubTokenResponse, error) {
	clientID := os.Getenv("CLIENT_ID_GITHUB")
	clientSecret := os.Getenv("CLIENT_SECRET_GITHUB")

	formData := url.Values{}
	formData.Set("client_id", clientID)
	formData.Set("client_secret", clientSecret)
	formData.Set("code", code)

	req, err := http.NewRequest("POST", "https://github.com/login/oauth/access_token", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create github token request: %w", err)
	}

	req.URL.RawQuery = formData.Encode()
	req.Header.Set("Accept", "application/json") // Crucial: forces GitHub to return JSON instead of URL-encoded query string

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send github token request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("github token endpoint returned status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var tokenResp GitHubTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return nil, fmt.Errorf("failed to decode github token response: %w", err)
	}

	return &tokenResp, nil
}

// FetchGitHubUser calls GitHub's REST API using the access token to retrieve user claims
func FetchGitHubUser(ctx context.Context, accessToken string) (*GitHubUser, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", "https://api.github.com/user", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create github user request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute github user request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("github user endpoint returned status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var ghUser GitHubUser
	if err := json.NewDecoder(resp.Body).Decode(&ghUser); err != nil {
		return nil, fmt.Errorf("failed to decode github user response: %w", err)
	}

	// If the user's primary email is private on GitHub, /user returns an empty email.
	// We must query the /user/emails endpoint explicitly using user:email scope.
	if ghUser.Email == "" {
		email, err := fetchPrimaryGitHubEmail(ctx, accessToken)
		if err == nil {
			ghUser.Email = email
		}
	}

	return &ghUser, nil
}

// fetchPrimaryGitHubEmail queries the separate /user/emails endpoint for private or primary addresses
func fetchPrimaryGitHubEmail(ctx context.Context, accessToken string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", "https://api.github.com/user/emails", nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("status code %d", resp.StatusCode)
	}

	var emails []GitHubEmail
	if err := json.NewDecoder(resp.Body).Decode(&emails); err != nil {
		return "", err
	}

	for _, e := range emails {
		if e.Primary && e.Verified {
			return e.Email, nil
		}
	}

	return "", fmt.Errorf("no verified primary email found")
}