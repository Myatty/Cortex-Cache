package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
)

// Home handler function
// http.ResponseWriter provides method for HTTP Response and sending it to user
// *http.Request is pointer to struct which holds info about current request(HTTP method and URL being requested)
func (app *application) home(w http.ResponseWriter, r *http.Request) {

	// checks if URL path is not "/", it returns error Page
	// if we dont return it will also writes (w.Write)
	if r.URL.Path != "/" {
		app.notFound(w)
		return
	}

	// NOTE that the file containing base template must be the *first* file in the slice
	files := []string{
		"./ui/html/base.tmpl.html",
		"./ui/html/partials/nav.tmpl.html",
		"./ui/html/pages/home.tmpl.html",
	}

	// template.ParseFiles() function reads the template file into a template set
	// pass files as variadic parameter
	ts, err := template.ParseFiles(files...)
	if err != nil {
		//app.errorLog.Print(err.Error())
		// http.Error(w, "Internal Server Error", 500)
		app.serverError(w, err)
		return
	}

	// the ExecuteTemplate() method to write the content of the "base" template as response body
	// The last parameter to Execute() represents any dynamic data that we want to pass in
	err = ts.ExecuteTemplate(w, "base", nil)
	if err != nil {
		app.errorLog.Print(err.Error())
		// http.Error(w, "Internal Server Error", 500)
		app.serverError(w, err)
	}

}

func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {

	// return 404 not found error if requested id is not correct
	id, err := strconv.Atoi(r.URL.Query().Get("id"))

	if err != nil || id < 1 {
		// http.NotFound(w, r)
		app.notFound(w)
		return
	}

	// use fmt.Fprintf to interpolate id value with response and write it to http.ResponseWriter
	fmt.Fprintf(w, "Displaying a specific snippet with ID: %d", id)
}

func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {

	// if we wanna send a non 200 status code, we must call w.WriteHeader()(which limit to only one for each response)
	// we must set all Headers before WriteHeader
	if r.Method != "POST" {

		w.Header().Set("Allow", http.MethodPost)
		// w.WriteHeader(405)
		// w.Write([]byte("Method Not Allowed!"))
		// http.Error(w, "Method not Allowed!", http.StatusMethodNotAllowed)
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	w.Write([]byte("Creating Snippet"))
}

// below is serving a single file(NOTE: it doesnt sanitize the path so BE CAREFUL, use filePath.Clean())
// func downloadHandler(w http.ResponseWriter, r *http.Request) {
// 		http.ServeFile(w, r, "./ui/static/file.zip")
// }
