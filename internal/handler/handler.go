package handler

import (
	"fmt"
	"goRent/internal/config"
	"goRent/internal/driver/mysqlDriver"
	"goRent/internal/form"
	"goRent/internal/helper"
	"goRent/internal/model"
	"goRent/internal/render"
	"goRent/internal/repository"
	"goRent/internal/repository/mysql"
	"html"
	"net/http"
	"strings"
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
	if true {
		helper.ServerError(w, r, fmt.Errorf("Test"))
		return
	}
	data := make(map[string]interface{})

	if err := render.Template(w, r, "home.page.html", &render.TemplateData{
		Data: data,
	}); err != nil {
		m.App.Error.Println(err)
	}
}

func (m *Repository) Search(w http.ResponseWriter, r *http.Request) {

	m.CreateProductList()

	if r.Method == http.MethodPost {
		searchKW := r.FormValue("searchtext")

		sr := html.EscapeString(searchKW)
		srtL := strings.ToLower((sr))
		_ = srtL
		// updateUserLastSearch(srtL, u)
		// insertUserSearchLogs(srtL, u)

		http.Redirect(w, r, "/searchresult", http.StatusSeeOther)
	}

	if err := render.Template(w, r, "home.page.html", &render.TemplateData{
		Data: nil,
	}); err != nil {
		m.App.Error.Println(err)
	}

}

func (m *Repository) Register(w http.ResponseWriter, r *http.Request) {
	m.App.Info.Println("Register: no session in progress")
	data := make(map[string]interface{})

	if err := render.Template(w, r, "register.page.html", &render.TemplateData{
		Data: data,
	}); err != nil {
		m.App.Error.Println(err)
	}
}

func (m *Repository) RegisterPost(w http.ResponseWriter, r *http.Request) {
	m.App.Info.Println("Register: no session in progress")
	data := make(map[string]interface{})
	if r.Method == http.MethodPost {
		if err := r.ParseForm(); err != nil {
			m.App.Error.Println(err)
		}

		newUser := model.User{
			Username: r.FormValue("inputUsername"),
			Email:    r.FormValue("inputEmail"),
			Password: r.FormValue("inputPassword"),
			Address: model.Address{
				Block:      r.FormValue("addressblock"),
				StreetName: r.FormValue("inputAddress"),
				UnitNumber: r.FormValue("addressunit"),
				PostalCode: r.FormValue("postalcode"),
			},
		}

		form := form.New(r.PostForm)
		form.Required("inputUsername", "inputEmail", "inputEmail", "inputPassword", "addressblock", "inputAddress", "addressunit", "postalcode")
		if form.ExistingUser() {
			form.Errors.Add("username", "Username already in use")
		}
		fmt.Println(newUser)
		fmt.Println(m.DB.GetUser(newUser.Username))
		//...
		// Sample email
		// msg := model.MailData{
		// 	To:       reservation.Email,
		// 	From:     "me@here.com",
		// 	Subject:  "Reservation Confirmation",
		// 	Content:  "",
		// 	Template: "basic.html",
		// }
		// msg.Content = fmt.Sprintf(`
		// 	<strong>Reservation Confirmation</strong><br>
		// 	Dear Mr/Ms %s, <br>
		// 	This is to confirm your reservation from %s to %s.
		// `,
		// 	reservation.LastName,
		// 	reservation.StartDate.Format(datelayout),
		// 	reservation.EndDate.Format(datelayout),
		// )

		// m.App.MailChan <- msg

		m.App.Info.Println("Register: redirecting to login page")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	if err := render.Template(w, r, "register.page.html", &render.TemplateData{
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
