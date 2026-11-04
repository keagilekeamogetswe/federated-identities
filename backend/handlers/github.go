package handlers

import (
	"fmt"
	"net/http"
)

func GithubRequestIdentityHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "selected Github Provider")
}

func GithubConfirmIdentityHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "selected Github Provider")
}