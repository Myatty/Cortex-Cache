package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)

func main() {

	// default is :4000, value is stored in addr variable
	addr := flag.String("addr", ":4000", "HTTP network address")

	// flag.Parse() to parse cl flag
	flag.Parse()

	// third one is the flags to indicate what additional information to include (local date and time).
	// Note that the flags are joined using the bitwise OR operator |.
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)

	// use the log.Lshortfile flag to include the relevant file name and line number.
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	// ServeMux is a router which in this case register home function as handler for URL "/"
	// create our own mux becoz DefaultServeMux is a global variable, any package can access it (Security Conerns)
	mux := http.NewServeMux()

	// create a file server which serves files out of the "./ui/static" directory
	fileServer := http.FileServer(http.Dir("./ui/static/"))

	// register the file server as the handler for all URL paths that start with "/static/"
	// strips the "/static" prefix before the request reaches the file server.
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	// URL "/" sub-tree pattern is a catch-all, all URL requests will be handled by this(its like "/**")
	mux.HandleFunc("/", home)
	mux.HandleFunc("/snippet/view", snippetView)
	mux.HandleFunc("/snippet/create", snippetCreate)

	// The value returned from the flag.String() function is a pointer to the flag value, not the value itself.
	// Note: any error returned by http.ListenAndServe is always non-nil
	infoLog.Printf("Starting server on port %s", *addr)
	err := http.ListenAndServe(*addr, mux)
	errorLog.Fatal(err)
}
