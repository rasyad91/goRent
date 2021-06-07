package handler

import (
	"goRent/internal/config"
	"goRent/internal/driver/mysqlDriver"
	"goRent/internal/render"
	"goRent/internal/repository"
	"goRent/internal/repository/mysql"
	"net/http"
)

type Repository struct {
	DB  repository.DatabaseRepo
	App *config.AppConfig
}

var Repo *Repository

// NewMySQLHandler creates db repo
func NewMySQLHandler(db *mysqlDriver.DB, app *config.AppConfig) *Repository {
	return &Repository{
		DB:  mysql.NewRepo(db.SQL),
		App: app,
	}
}

// New creates the handlers
func New(r *Repository) {
	Repo = r
}

func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {
	if f := m.App.Session.Get(r.Context(), "flash"); f != nil {
		m.App.Session.Put(r.Context(), "warning", nil)
	} else {
		m.App.Session.Put(r.Context(), "warning", "hello")
	}
	// m.App.Session.Put(r.Context(), "flash", "let's see")
	data := make(map[string]interface{})

	// u := m.App.Session.Get(r.Context(), "user").(model.User)
	// fmt.Println("PRINTINT U", u)
	// fmt.Println("checking authenticate", m.App.Session.Exists(r.Context(), "userID"))
	// fmt.Println(u.ID)
	if err := render.Template(w, r, "home.page.html", &render.TemplateData{
		Data: data,
	}); err != nil {
		m.App.Error.Println(err)
	}
}
