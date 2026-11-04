package handlers

import (
	"fmt"
	"net/http"
)

func APIHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello world")
}