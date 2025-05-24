package main

import (
	"log"
	"net/http"
)

func main() {

	// ServeMux is a router which in this case register home function as handler for URL "/"
	// URL "/" sub-tree pattern is a catch-all, all URL requests will be handled by this(its like "/**")

	// create our own mux becoz DefaultServeMux is a global variable, any package can access it (Security Conerns)
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
