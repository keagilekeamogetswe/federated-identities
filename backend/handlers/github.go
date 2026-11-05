package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/keagilekeamogetswe/federated-identities/services"
)

// GithubRequestIdentityHandler generates a CSRF state token, sets a secure cookie,
// and requests the GitHub OAuth authorization URL.
func GithubRequestIdentityHandler(w http.ResponseWriter, r *http.Request) {
	stateToken, err := services.GenerateGitHubStateToken()
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Store state in a cookie for CSRF verification on callback
	http.SetCookie(w, &http.Cookie{
		Name:     "github_oauth_state",
		Value:    stateToken,
		Expires:  time.Now().Add(10 * time.Minute),
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
	})

	authURL := services.RequestGitHubAuthURL(stateToken)
	// Return the authorization URL as a JSON payload
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"url": authURL,
	})
}

// GithubConfirmIdentityHandler handles the callback from GitHub, validates the state,
// exchanges the authorization code for an access token, fetches the user profile,
// and extracts verified user identity claims.
func GithubConfirmIdentityHandler(w http.ResponseWriter, r *http.Request) {
	// 1. Verify CSRF state token from cookie vs query parameter
	cookie, err := r.Cookie("github_oauth_state")
	if err != nil {
		http.Error(w, "State cookie missing", http.StatusBadRequest)
		return
	}

	queryState := r.URL.Query().Get("state")
	if queryState == "" || queryState != cookie.Value {
		http.Error(w, "Invalid or mismatched state token", http.StatusBadRequest)
		return
	}

	// Clear the state cookie immediately (single-use defense)
	http.SetCookie(w, &http.Cookie{
		Name:     "github_oauth_state",
		Value:    "",
		Expires:  time.Now().Add(-1 * time.Hour),
		HttpOnly: true,
		MaxAge:   -1,
		Path:     "/",
	})

	// 2. Extract authorization code returned by GitHub
	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "Authorization code missing from callback", http.StatusBadRequest)
		return
	}

	// 3. Execute back-channel token exchange
	tokenResp, err := services.ExchangeGitHubCode(code)
	if err != nil {
		http.Error(w, fmt.Sprintf("GitHub token exchange failed: %v", err), http.StatusInternalServerError)
		return
	}

	// 4. Fetch the authenticated user profile (and handle private emails automatically via service)
	ghUser, err := services.FetchGitHubUser(r.Context(), tokenResp.AccessToken)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to fetch GitHub user profile: %v", err), http.StatusUnauthorized)
		return
	}

	// 5. Extract verified user profile claims
	githubID := strconv.Itoa(ghUser.ID) // Unique permanent numeric GitHub user ID
	email := ghUser.Email
	name := ghUser.Name
	if name == "" {
		name = ghUser.Login // Fallback to GitHub username/login if real name is not set
	}
	picture := ghUser.AvatarURL

	if email == "" {
		http.Error(w, "GitHub email could not be retrieved or is private", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":     "authenticated",
		"github_id":  githubID,
		"email":      email,
		"name":       name,
		"picture":    picture,
		"username":   ghUser.Login,
	})
}