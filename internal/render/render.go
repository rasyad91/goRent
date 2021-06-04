package render

import (
	"fmt"
	"goRent/internal/config"
	"goRent/internal/form"
	"goRent/internal/helper"
	"goRent/internal/model"
	"html/template"
	"log"
	"net/http"
	"runtime/debug"

	"github.com/justinas/nosurf"
)

var app *config.AppConfig

// var src = rand.NewSource(time.Now().UnixNano())

// NewHelpers creates new helpers
func New(a *config.AppConfig) {
	app = a
}

// TemplateData stores data to be used in Templates
type TemplateData struct {
	Data            map[string]interface{}
	Form            *form.Form
	Products        []model.Product
	CSRFToken       string
	IsAuthenticated bool
	User            model.User
	Flash           string
	Warning         string
	Error           string
}

// Template parses and exectues template by its template name
func Template(w http.ResponseWriter, r *http.Request, tmpl string, td *TemplateData) error {

	t := DefaultData(w, r, *td)

	ts, err := template.New(tmpl).Funcs(function).ParseFiles(fmt.Sprintf("./templates/%s", tmpl), "./templates/base.layout.html", "./templates/header.layout.html", "./templates/footer.layout.html")
	if err != nil {
		return fmt.Errorf("ParseTemplate: Unable to find template pages: %w", err)
	}

	if err := ts.Execute(w, t); err != nil {
		return fmt.Errorf("ParseTemplate: Unable to execute template: %w", err)
	}

	return nil
}

// ServerError will display error page for internal server error
func ServerError(w http.ResponseWriter, r *http.Request, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	_ = log.Output(2, trace)

	w.WriteHeader(http.StatusInternalServerError)
	w.Header().Set("Connection", "closec")
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, post-check=0, pre-check=0")
	http.ServeFile(w, r, "./static/500.html")

}

// DefaultData adds default data which is accessible to all templates
func DefaultData(w http.ResponseWriter, r *http.Request, td TemplateData) TemplateData {
	td.CSRFToken = nosurf.Token(r)
	td.IsAuthenticated = helper.IsAuthenticated(r)
	// if logged in, store user id in template data
	if td.IsAuthenticated {
		u := app.Session.Get(r.Context(), "user").(model.User)
		td.User = u
	}
	td.Flash = app.Session.PopString(r.Context(), "flash")
	td.Warning = app.Session.PopString(r.Context(), "warning")
	td.Error = app.Session.PopString(r.Context(), "error")

	return td
}
