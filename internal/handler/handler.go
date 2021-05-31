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

// NewMySQLHandler creates db repo
func NewMySQLHandler(db *mysqlDriver.DB, errorLog, infoLog *log.Logger) *Repository {
	return &Repository{
		DB:    mysql.NewRepo(db.SQL),
		Error: errorLog,
		Info:  infoLog,
	}
}

// New creates the handlers
func New(r *Repository) {
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

func (m *Repository) GetCourse(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to REST API\n")
}

func (m *Repository) PostCourse(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to REST API\n")
}

func (m *Repository) PutCourse(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to REST API\n")
}

func (m *Repository) DeleteCourse(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to REST API\n")
}
