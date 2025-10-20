package main

import (
	"errors"
	"fmt"

	"net/http"
	"strconv"

	"cortexcache.myatty.net/internal/models"
	"github.com/julienschmidt/httprouter"
)

// Home handler function
// http.ResponseWriter provides method for HTTP Response and sending it to user
// *http.Request is pointer to struct which holds info about current request(HTTP method and URL being requested)
func (app *application) home(w http.ResponseWriter, r *http.Request) {

	// httprouter can check this, so ill just remove this
	// checks if URL path is not "/", it returns error Page
	// if r.URL.Path != "/" {
	// 	app.notFound(w)
	// 	return
	// }

	snippets, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}

	data := app.newTemplateData(r)
	data.Snippets = snippets

	app.render(w, http.StatusOK, "home.tmpl.html", data)

}

func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {

	params := httprouter.ParamsFromContext(r.Context())

	// return 404 not found error if requested id is not valid
	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	snippet, err := app.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	data := app.newTemplateData(r)
	data.Snippet = snippet

	app.render(w, http.StatusOK, "view.tmpl.html", data)
}

func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Display the form fo creating new Snippet ... "))
}

func (app *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {

	// checking if its POST method or not is done by httprouter, so ill remove this also
	// if we wanna send a non 200 status code, we must call w.WriteHeader()(which limit to only one for each response)
	// we must set all Headers before WriteHeader
	// if r.Method != http.MethodPost {
	// 	w.Header().Set("Allow", http.MethodPost)
	// 	app.clientError(w, http.StatusMethodNotAllowed)
	// 	return
	// }

	title := "0 snail"
	content := "O snail\nClimb Mount Fuji,\nBut slowly, slowly!\n\n– Myint Myat"
	expires := 7

	id, err := app.snippets.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// refactor "/snippet/view?id=%d" becoz httprouter can provide Clean URL format
	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}

// below is serving a single file(NOTE: it doesnt sanitize the path so BE CAREFUL, use filePath.Clean())
// func downloadHandler(w http.ResponseWriter, r *http.Request) {
// 		http.ServeFile(w, r, "./ui/static/file.zip")
// }
