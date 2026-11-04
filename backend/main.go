package main

import (
	"log"
	"net/http"
	"github.com/keagilekeamogetswe/federated-identities/handlers"
)


func main(){
	mux := http.NewServeMux()
	mux.HandleFunc("/api/", handlers.APIHandler)

	//The server
	server := http.Server{
		Addr: ":8080",
		Handler: mux,
	}
	log.Println("Server listening...")
	// Since I am just starting out with go I'm avoiding error handling to keep things simple.
	server.ListenAndServe()
}