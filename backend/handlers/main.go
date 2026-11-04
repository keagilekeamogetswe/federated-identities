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
	switch identity_provider {
		case "google":
			GoogleRequestIdentityHandler(
				response,
				request,
			)

		case "github":
			GithubRequestIdentityHandler(
				response,
				request,
			)

		default:
			http.Error(
				response,
				"Unknown provider",
				http.StatusBadRequest,
			)
		}
}
func ProviderConfirm(response http.ResponseWriter, request *http.Request){
	// Intercept the dynamic route
	identity_provider := strings.ToLower(chi.URLParam(request, "name"))
	// Route based on identity provider
	switch identity_provider {
		case "google":
			GoogleConfirmIdentityHandler(
				response,
				request,
			)

		case "github":
			GithubConfirmIdentityHandler(
				response,
				request,
			)

		default:
			http.Error(
				response,
				"Unknown provider",
				http.StatusBadRequest,
			)
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