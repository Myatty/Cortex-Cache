package main

import (
	"fmt"
	"net/http"
	"strconv"
)

// Home handler function
// http.ResponseWriter provides method for HTTP Response and sending it to user
// *http.Request is pointer to struct which holds info about current request(HTTP method and URL being requested)
func home(w http.ResponseWriter, r *http.Request) {

	// checks if URL path is not "/", it returns error Page
	// if we dont return it will also writes (w.Write)
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	w.Write([]byte("Hello from Cortex Cache"))
}

func snippetView(w http.ResponseWriter, r *http.Request) {

	// return 404 not found error if requested id is not correct
	id, err := strconv.Atoi(r.URL.Query().Get("id"))

	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	// use fmt.Fprintf to interpolate id value with response and write it to http.ResponseWriter
	fmt.Fprintf(w, "Displaying a specific snippet with ID: %d", id)
}

func snippetCreate(w http.ResponseWriter, r *http.Request) {

	// if we wanna send a non 200 status code, we must call w.WriteHeader()(which limit to only one for each response)
	// we must set all Headers before WriteHeader
	if r.Method != "POST" {

		w.Header().Set("Allow", http.MethodPost)
		// w.WriteHeader(405)
		// w.Write([]byte("Method Not Allowed!"))
		http.Error(w, "Method not Allowed!", http.StatusMethodNotAllowed)

		return
	}

	w.Write([]byte("Creating Snippet"))
}
