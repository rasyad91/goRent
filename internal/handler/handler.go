package handler

import (
	"fmt"
	"goRent/internal/config"
	"goRent/internal/driver/mysqlDriver"
	"goRent/internal/form"
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

func (m *Repository) Login(w http.ResponseWriter, r *http.Request) {

	m.App.Info.Println("Login: no session in progress")

	data := make(map[string]interface{})

	if r.Method == http.MethodPost {
		if err := r.ParseForm(); err != nil {
			m.App.Error.Println(err)
			return
		}

		//...
	}

	if err := render.Template(w, r, "login.page.html", &render.TemplateData{
		Data: data,
	}); err != nil {
		m.App.Error.Println(err)
	}
}

func (m *Repository) Register(w http.ResponseWriter, r *http.Request) {

	m.App.Info.Println("Register: no session in progress")

	data := make(map[string]interface{})

	if r.Method == http.MethodPost {
		if err := r.ParseForm(); err != nil {
			m.App.Error.Println(err)
		}
		newUser := model.User{
			Username:   r.FormValue("inputUsername"),
			Email:      r.FormValue("inputEmail"),
			Password:   []byte(r.FormValue("inputPassword")),
			Block:      r.FormValue("addressblock"),
			StreetName: r.FormValue("inputAddress"),
			Unit:       r.FormValue("addressunit"),
			PostalCode: r.FormValue("postalcode"),
		}
		form := form.New(r.PostForm)
		form.Required("inputUsername", "inputEmail", "inputEmail", "inputPassword", "addressblock", "inputAddress", "addressunit", "postalcode")
		if form.ExistingUser() {
			form.Errors.Add("username", "Username already in use")
		}
		fmt.Println(newUser)
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
