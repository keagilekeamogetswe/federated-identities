package handlers

import (
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
	http.Redirect(w, r, authURL, http.StatusTemporaryRedirect)
}
func GoogleConfirmIdentityHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "selected google Provider")
}