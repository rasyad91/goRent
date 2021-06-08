package handler

import (
	"database/sql"
	"fmt"
	"goRent/internal/form"
	"goRent/internal/model"
	"goRent/internal/render"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func (m *Repository) Login(w http.ResponseWriter, r *http.Request) {
	fmt.Println("SUCCESSFULLY HIT LOGINGET")
	data := make(map[string]interface{})
	data["login"] = model.User{
		Username: "",
		Email:    "",
		Password: "",
		Address:  model.Address{},
	}
	if err := render.Template(w, r, "login.page.html", &render.TemplateData{
		Data: data,
		Form: &form.Form{},
	}); err != nil {
		m.App.Error.Println(err)
	}
}

func (m *Repository) LoginPost(w http.ResponseWriter, r *http.Request) {
	fmt.Println("HITTING LOGINPOST")
	data := make(map[string]interface{})
	form := form.New(r.PostForm)
	if err := r.ParseForm(); err != nil {
		m.App.Error.Println(err)
		return
	}
	_ = m.App.Session.RenewToken(r.Context())
	Username := r.FormValue("username")
	password := r.FormValue("password")
	eu, err := m.DB.GetUser(Username)
	if err != nil {
		if err != sql.ErrNoRows {
			m.App.Error.Println(err)
		}
		form.Errors.Add("login", "Username or password incorrect")
	}
	fmt.Println("SUCCESSFULLY PULLED USER INFO")
	// fmt.Println(password)
	// t, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	// fmt.Println(t)
	// fmt.Println(string(t))
	// fmt.Println([]byte(t))
	// fmt.Println(string(t))

	// fmt.Println(eu.Password)
	err = bcrypt.CompareHashAndPassword([]byte(eu.Password), []byte(password))
	if err != nil {
		form.Errors.Add("login", "Username or password incorrect")
	} else {
		http.Redirect(w, r, "/home", http.StatusSeeOther)
		fmt.Println("PASSWORD MATCHES")
	}
	if len(form.Errors) != 0 {
		if err := render.Template(w, r, "login.page.html", &render.TemplateData{
			Form: form,
			Data: data,
		}); err != nil {
			m.App.Error.Println(err)
		}
		return
	}

	for _, v := range eu.Rents {
		fmt.Printf("%#v\n", v)
	}
	for _, v := range eu.Bookings {
		fmt.Println(v)
	}

	m.App.Session.Put(r.Context(), "userID", eu.ID)
	m.App.Session.Put(r.Context(), "flash", fmt.Sprintf("Welcome, %s", eu.Username))
	m.App.Session.Put(r.Context(), "user", eu)

	http.Redirect(w, r, "/", http.StatusSeeOther)
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
	fmt.Println("Setting Cookie here")
	_ = m.App.Session.RenewToken(r.Context())
	_ = m.App.Session.Destroy(r.Context())
	_ = m.App.Session.RenewToken(r.Context())
	// fmt.Println("LOGOUT Session", m.App.Session.Get(r.Context(), "user"))
	m.App.Session.Put(r.Context(), "flash", "You've been logged out successfully!")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// WHEN USER REACH LOGIN SCREEN:
// // LoginScreen shows the home (login) screen
// func (repo *DBRepo) LoginScreen(w http.ResponseWriter, r *http.Request) {
//  // if already logged in, take to dashboard
//  if repo.App.Session.Exists(r.Context(), "userID") {
//   http.Redirect(w, r, "/admin/overview", http.StatusSeeOther)
//   return
//  }

//  err := helpers.RenderPage(w, r, "login", nil, nil)
//  if err != nil {
//   printTemplateError(w, err)
//  }
// }
// WHEN USER ATTEMPTS TO LOGIN:

// // Login attempts to log the user in
// func (repo *DBRepo) Login(w http.ResponseWriter, r *http.Request) {
//  _ = repo.App.Session.RenewToken(r.Context())
//  err := r.ParseForm()
//  if err != nil {
//   log.Println(err)
//   ClientError(w, r, http.StatusBadRequest)
//   return
//  }

//  id, hash, err := repo.DB.Authenticate(r.Form.Get("email"), r.Form.Get("password"))
//  if err == models.ErrInvalidCredentials {
//   app.Session.Put(r.Context(), "error", "Invalid login")
//   err := helpers.RenderPage(w, r, "login", nil, nil)
//   if err != nil {
//    printTemplateError(w, err)
//                         return
//   }

//  // we authenticated. Get the user.
//  u, err := repo.DB.GetUserById(id)
//  if err != nil {
//   log.Println(err)
//   ClientError(w, r, http.StatusBadRequest)
//   return
//  }

//  app.Session.Put(r.Context(), "userID", id)
//  app.Session.Put(r.Context(), "flash", "You've been logged in successfully!")
//  app.Session.Put(r.Context(), "user", u)

//  http.Redirect(w, r, "/admin/overview", http.StatusSeeOther)
// }

// m.App.Session.Put(r.Context(), "success", "hello")
// fmt.Printf("%#v", m.App.Session)
// //to do
// - add in notification for register and login
// - feed in user info for header
