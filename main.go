package main

import (
	"log"
	"net/http"
)

// Home handler function
// http.ResponseWriter provides method for HTTP Response and sending it to user
// *http.Request is pointer to struct which holds info about current request(HTTP method and URL being requested)
func home(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello from Cortex Cache"))
}

func main() {

	// ServeMux is a router which in this case register home function as handler for URL "/"
	mux := http.NewServeMux()
	mux.HandleFunc("/", home)

	// http.ListenAndServe() starts new web server and now it listens on tcp port 4000
	// Note: any error returned by http.ListenAndServe is always non-nil
	log.Print("Starting server on port : 4000")
	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)
}
