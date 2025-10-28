package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"cortexcache.myatty.net/internal/models"
	"cortexcache.myatty.net/internal/validator"
	"github.com/julienschmidt/httprouter"
)

type snippetCreateForm struct {
	Title               string `form:"title"`
	Content             string `form:"content"`
	Expires             int    `form:"expires"`
	validator.Validator `form:"-"`
}

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
	data := app.newTemplateData(r)

	data.Form = snippetCreateForm{
		Expires: 365,
	}

	app.render(w, http.StatusOK, "create.tmpl.html", data)
}

func (app *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {

	var form snippetCreateForm

	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form.CheckField(validator.NotBlank(form.Title), "title", "This field cannot be blank")
	form.CheckField(validator.MaxChars(form.Title, 100), "title", "This field cannot be more than 100 characters long")
	form.CheckField(validator.NotBlank(form.Content), "content", "This field cannot be blank")
	form.CheckField(validator.PermittedInt(form.Expires, 1, 7, 365), "expires", "This field must equal 1, 7 or 365")

	// Dump errors(if exists) into HTTP response and return from Handler
	// pass form data as dynamic data and return 422 error
	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "create.tmpl.html", data)
		return
	}

	// Insert into database
	id, err := app.snippets.Insert(form.Title, form.Content, form.Expires)
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.sessionManager.Put(r.Context(), "flash", "Snippet Successfully Created.")

	// refactor "/snippet/view?id=%d" becoz httprouter can provide Clean URL format
	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}

// below is serving a single file(NOTE: it doesnt sanitize the path so BE CAREFUL, use filePath.Clean())
// func downloadHandler(w http.ResponseWriter, r *http.Request) {
// 		http.ServeFile(w, r, "./ui/static/file.zip")
// }
