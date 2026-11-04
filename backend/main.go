package main

import (
	"fmt"
	"log"
	"net/http"
)
// Homehandler will be removed the only routes that will be server are api/*
func APIHandler (res http.ResponseWriter, req *http.Request) {
	fmt.Fprintln(res, "Hello world")
}

func main(){
	mux := http.NewServeMux()
	mux.HandleFunc("/api/", APIHandler)

	//The server
	server := http.Server{
		Addr: ":8080",
		Handler: mux,
	}
	log.Println("Server listening...")
	// Since I am just starting out with go I'm avoiding error handling to keep things simple.
	server.ListenAndServe()
}