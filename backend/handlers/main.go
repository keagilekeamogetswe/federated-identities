package handlers

import (
	"strings"
	"net/http"
	"github.com/go-chi/chi/v5"
)


func ProviderRequest( response http.ResponseWriter, request *http.Request){
	// Intercept the dynamic route
	identity_provider := strings.ToLower(chi.URLParam(request, "name"))
	// Route based on identity provider
	if identity_provider == "google" {
		GoogleRequestIdentityHandler(response, request)
		return
	} else if (identity_provider == "github") {
		GithubRequestIdentityHandler(response, request)
		return
	}
}
func ProviderConfirm(response http.ResponseWriter, request *http.Request){
	// Intercept the dynamic route
	identity_provider := strings.ToLower(chi.URLParam(request, "name"))
	// Route based on identity provider
	if identity_provider == "google"{
		GoogleConfirmIdentityHandler(response, request)
		return
	}else if identity_provider == "gihub" {
		GithubConfirmIdentityHandler(response, request)
		return
	}
}


func MainHandler() http.Handler {
	api_router := chi.NewRouter()
	//GET is for requesting identity
	api_router.Get("/provider/{name}", ProviderRequest)
	//POST Method is dedicated for confirming identity
	api_router.Post("/provider/{name}", ProviderConfirm)
	return api_router
}