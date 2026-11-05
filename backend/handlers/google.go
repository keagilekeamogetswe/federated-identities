package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/keagilekeamogetswe/federated-identities/services"
)

// GoogleRequestIdentityHandler generates a CSRF state token, sets a secure cookie,
// and redirects the user to Google's OAuth consent screen.
func GoogleRequestIdentityHandler(w http.ResponseWriter, r *http.Request) {
	stateToken, err := services.GenerateStateToken()
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Store state in a cookie for CSRF verification on callback
	http.SetCookie(w, &http.Cookie{
		Name:     "oauth_state",
		Value:    stateToken,
		Expires:  time.Now().Add(10 * time.Minute),
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
	})

	authURL := services.RequestGoogleAuthURL(stateToken)
<<<<<<< HEAD
	// Return the authorization URL as a JSON payload instead of redirecting
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"url": authURL,
	})
=======
	http.Redirect(w, r, authURL, http.StatusTemporaryRedirect)
>>>>>>> 946fa1c742e86fd2f4e129262d6fa98d06a87a12
}
// GoogleConfirmIdentityHandler handles the callback from Google, validates the state,
// exchanges the authorization code for tokens, and verifies the ID token.
func GoogleConfirmIdentityHandler(w http.ResponseWriter, r *http.Request) {
	// 1. Verify CSRF state token from cookie vs query parameter
	cookie, err := r.Cookie("oauth_state")
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
		Name:     "oauth_state",
		Value:    "",
		Expires:  time.Now().Add(-1 * time.Hour),
		HttpOnly: true,
		MaxAge:   -1,
		Path:     "/",
	})

	// 2. Extract authorization code returned by Google
	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "Authorization code missing from callback", http.StatusBadRequest)
		return
	}

	// 3. Execute back-channel token exchange
	tokenResp, err := services.DirectServerConfirmation(code)
	if err != nil {
		http.Error(w, fmt.Sprintf("Token exchange failed: %v", err), http.StatusInternalServerError)
		return
	}

	// 4. Cryptographically verify Google's JWT ID Token
	payload, err := services.ConfirmGoogleID(r.Context(), tokenResp.IDToken)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to verify ID token: %v", err), http.StatusUnauthorized)
		return
	}

	// 5. Extract verified user profile claims
	googleID := payload.Subject // Unique permanent Google user ID (`sub`)
	email, _ := payload.Claims["email"].(string)
	emailVerified, _ := payload.Claims["email_verified"].(bool)
	name, _ := payload.Claims["name"].(string)
	picture, _ := payload.Claims["picture"].(string)

	if !emailVerified {
		http.Error(w, "Google email is not verified", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":    "authenticated",
		"google_id": googleID,
		"email":     email,
		"name":      name,
		"picture":   picture,
	})
}