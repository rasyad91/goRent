//Handler for Login
package handler

import (
	"database/sql"
	"fmt"
	"goRent/internal/form"
	"goRent/internal/model"
	"goRent/internal/render"
	"net/http"
	"time"

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

//Checks if username and password are correct, if so , GETUSER and store in session
func (m *Repository) LoginPost(w http.ResponseWriter, r *http.Request) {
	// fmt.Println("HITTING LOGINPOST")
	t := time.Now()
	fmt.Println("Start timing...")
	form := form.New(r.PostForm)
	if err := r.ParseForm(); err != nil {
		m.App.Error.Println(err)
		render.ServerError(w, r, err)
		return
	}
	_ = m.App.Session.RenewToken(r.Context())
	Username := r.FormValue("username")
	password := r.FormValue("password")

	eu, err := m.DB.GetUser(Username)
	// fmt.Println("LOGIN ACCESS CHECK", eu.AccessLevel)
	if err != nil {
		if err != sql.ErrNoRows {
			m.App.Error.Println(err)
			render.ServerError(w, r, err)
			return
		}
		form.Errors.Add("login", "Username or password incorrect")
	}
	m.App.Info.Println("successfully pulled user info")

	if err := bcrypt.CompareHashAndPassword([]byte(eu.Password), []byte(password)); err != nil {
		form.Errors.Add("login", "Username or password incorrect")
	}

	fmt.Println("time taken: ", time.Since(t))
	data := make(map[string]interface{})
	if len(form.Errors) != 0 {
		if err := render.Template(w, r, "login.page.html", &render.TemplateData{
			Form: form,
			Data: data,
		}); err != nil {
			m.App.Error.Println(err)
		}
		return
	}

	// fmt.Printf("user: %#v\n", eu.Email)

	// for _, v := range eu.Rents {
	// 	fmt.Printf("id :#%d processed:%t product:%s start:%s end:%s\n", v.ID, v.Processed, v.Product.Title, v.StartDate, v.EndDate)
	// }

	m.App.Session.Put(r.Context(), "userID", eu.ID)
	m.App.Session.Put(r.Context(), "flash", fmt.Sprintf("Welcome, %s", eu.Username))
	m.App.Session.Put(r.Context(), "user", eu)
	m.App.Session.Put(r.Context(), "accesslevel", eu.AccessLevel)

	url := m.App.Session.Get(r.Context(), "url")
	// fmt.Println("WHY THIS NOT WORKING", eu.AccessLevel == 1)
	if eu.AccessLevel == 1 {
		fmt.Println("Hitting login admin redirect")
		http.Redirect(w, r, "/admin/overview", http.StatusSeeOther)
		return
	}
	fmt.Println(url)
	if url != "" {
		fmt.Println("I'm hitting here")
		http.Redirect(w, r, m.App.Session.PopString(r.Context(), "url"), http.StatusSeeOther)
		return
	}
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
	_ = m.App.Session.RenewToken(r.Context())
	_ = m.App.Session.Destroy(r.Context())
	_ = m.App.Session.RenewToken(r.Context())
	// fmt.Println("LOGOUT Session", m.App.Session.Get(r.Context(), "user"))
	m.App.Session.Put(r.Context(), "flash", "You've been logged out successfully!")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
