package helper

import (
	"fmt"
	"goRent/internal/config"
	"goRent/internal/model"
	"goRent/internal/render"
	"log"
	"net/http"
	"runtime/debug"

	"github.com/justinas/nosurf"
)

const (
	letterBytes   = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

var app *config.AppConfig

// var src = rand.NewSource(time.Now().UnixNano())

// NewHelpers creates new helpers
func NewHelpers(a *config.AppConfig) {
	app = a
}

// IsAuthenticated returns true if a user is authenticated
func IsAuthenticated(r *http.Request) bool {
	exists := app.Session.Exists(r.Context(), "userID")
	return exists
}

// ServerError will display error page for internal server error
func ServerError(w http.ResponseWriter, r *http.Request, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	_ = log.Output(2, trace)

	w.WriteHeader(http.StatusInternalServerError)
	w.Header().Set("Connection", "close")
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, post-check=0, pre-check=0")
	if err := render.Template(w, r, "serverError.page.html", &render.TemplateData{}); err != nil {
		log.Println(err)
	}
}

// DefaultData adds default data which is accessible to all templates
func DefaultData(td render.TemplateData, r *http.Request, w http.ResponseWriter) render.TemplateData {
	td.CSRFToken = nosurf.Token(r)
	td.IsAuthenticated = IsAuthenticated(r)
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
