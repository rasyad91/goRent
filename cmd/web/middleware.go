package main

import (
	"fmt"
	"goRent/internal/helper"
	"goRent/internal/render"
	"net/http"
	"strings"

	"github.com/justinas/nosurf"
)

// SessionLoad loads the session on requests
func SessionLoad(next http.Handler) http.Handler {
	return session.LoadAndSave(next)
}

// Auth checks for authentication
func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !helper.IsAuthenticated(r) {
			app.Session.Put(r.Context(), "warning", "Please login first")
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}
		w.Header().Add("Cache-Control", "no-store")

		next.ServeHTTP(w, r)
	})
}

// LastGetURL stores the last GET URL that is not /login or /register
func LastGetURL(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			if strings.HasPrefix(r.URL.String(), "/v1") {
				app.Session.Put(r.Context(), "url", r.URL.String())
			}
		}
		next.ServeHTTP(w, r)
	})
}

// RecoverPanic recovers from a panic
func RecoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			// Check if there has been a panic
			if err := recover(); err != nil {
				// return a 500 Internal Server response
				render.ServerError(w, r, fmt.Errorf("%s", err))
				app.Error.Println(err)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

// NoSurf implements CSRF protection
func NoSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)

	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		SameSite: http.SameSiteStrictMode,
		Domain:   app.Domain,
	})

	return csrfHandler
}
