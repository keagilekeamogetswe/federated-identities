package handlers

import (
	"fmt"
	"net/http"
)

func GoogleRequestIdentityHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "selected google Provider")
}
func GoogleConfirmIdentityHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "selected google Provider")
}