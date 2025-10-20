package main

import (
	"net/http"
)

// The routes() method returns a servemux containing application routes
func (app *application) routes() http.Handler {

	// ServeMux is a router which in this case register home function as handler for URL "/"
	// create our own mux becoz DefaultServeMux is a global variable, any package can access it (Security Conerns)
	mux := http.NewServeMux()

	// create a file server which serves files out of the "./ui/static" directory
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	// URL "/" sub-tree pattern is a catch-all, all URL requests will be handled by this(its like "/**")
	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/snippet/view", app.snippetView)
	mux.HandleFunc("/snippet/create", app.snippetCreate)

	return app.recoverPanic(app.logRequest(secureHeaders(mux)))
}
