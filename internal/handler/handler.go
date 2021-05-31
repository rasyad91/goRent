package handler

import (
	"fmt"
	"goRent/internal/driver/mysqlDriver"
	"goRent/internal/repository"
	"goRent/internal/repository/mysql"
	"log"
	"net/http"
)

type Repository struct {
	DB    repository.DatabaseRepo
	Error *log.Logger
	Info  *log.Logger
}

var Repo *Repository

// NewMySQLHandlers creates db repo for postgres
func NewMySQLHandlers(db *mysqlDriver.DB, errorLog, infoLog *log.Logger) *Repository {
	return &Repository{
		DB:    mysql.NewRepo(db.SQL),
		Error: errorLog,
		Info:  infoLog,
	}
}

// NewHandlers creates the handlers
func NewHandlers(r *Repository) {
	Repo = r
}

func ValidationAPIMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Do stuff here
		v := r.URL.Query()
		key, ok := v["key"]
		if !ok || key[0] != "2c78afaf-97da-4816-bbee-9ad239abb296" {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("401 - Invalid key"))
			return
		}
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}

func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to REST API\n")
}
