package handler

import (
	"fmt"
	"goRent/internal/config"
	"goRent/internal/driver/mysqlDriver"
	"goRent/internal/model"
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
	m.App.Session.Put(r.Context(), "warning", "hello")

	data := make(map[string]interface{})

	u := m.App.Session.Get(r.Context(), "user").(model.User)
	fmt.Println("PRINTINT U", u)
	fmt.Println("checking authenticate", m.App.Session.Exists(r.Context(), "userID"))

	if err := render.Template(w, r, "home.page.html", &render.TemplateData{
		Data: data,
	}); err != nil {
		m.App.Error.Println(err)
	}
}

// Logout logs the user out
func (m *Repository) Logout(w http.ResponseWriter, r *http.Request) {

	// delete the remember me cookie, if any
	delCookie := http.Cookie{
		Name:     fmt.Sprintf("_%s_gowatcher_remember", m.App.PreferenceMap["identifier"]),
		Value:    "",
		Domain:   m.App.Domain,
		Path:     "/",
		MaxAge:   0,
		HttpOnly: true,
	}
	http.SetCookie(w, &delCookie)

	_ = m.App.Session.RenewToken(r.Context())
	_ = m.App.Session.Destroy(r.Context())
	_ = m.App.Session.RenewToken(r.Context())

	m.App.Session.Put(r.Context(), "flash", "You've been logged out successfully!")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
