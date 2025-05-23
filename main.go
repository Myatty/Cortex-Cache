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

func snippetView(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Displaying Snippets"))
}

func snippetCreate(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Creating Snippet"))
}

func main() {

	// ServeMux is a router which in this case register home function as handler for URL "/"
	// URL "/" sub-tree pattern is a catch-all, all URL requests will be handled by this(its like "/**")
	mux := http.NewServeMux()
	mux.HandleFunc("/", home)
	mux.HandleFunc("/snippet/view", snippetView)
	mux.HandleFunc("/snippet/create", snippetCreate)

	// http.ListenAndServe() starts new web server and now it listens on tcp port 4000
	// Note: any error returned by http.ListenAndServe is always non-nil
	log.Print("Starting server on port : 4000")
	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)
}
