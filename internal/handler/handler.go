package handler

import (
	"goRent/internal/config"
	"goRent/internal/render"
	"goRent/internal/repository"
	"net/http"
)

type Repository struct {
	DB  repository.DatabaseRepo
	App *config.AppConfig
}

// Repo is used by the handler
var Repo *Repository

// NewMySQLHandler creates db repo
func NewRepo(db repository.DatabaseRepo, app *config.AppConfig) *Repository {
	return &Repository{
		DB:  db,
		App: app,
	}
}

// New creates the handlers
func New(r *Repository) {
	Repo = r
}

func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {
	// m.App.Session.Put(r.Context(), "flash", "let's see")
	data := make(map[string]interface{})

	if err := render.Template(w, r, "home.page.html", &render.TemplateData{
		Data: data,
	}); err != nil {
		m.App.Error.Println(err)
	}
}
