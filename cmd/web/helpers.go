package main

import (
	"fmt"
	"net/http"
	"runtime/debug"
)

// serverError helper writes error and stack trace to errorLog and sends a 500
func (app *application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.errorLog.Output(2, trace)

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

// clientError sends specific status codes and description to the user
func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (app *application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}

func (app *application) render(w http.ResponseWriter, r *http.Request, name string, td *templateData) {

	// get the template set from cache for the supplied page name.
	ts, ok := app.templateCache[name]
	if !ok {
		app.serverError(w, fmt.Errorf("The template %s does not exist", name))
		return
	}

	// execute the template set, passing in any dynamic data
	err := ts.Execute(w, td)
	if err != nil {
		app.serverError(w, err)
	}
}
