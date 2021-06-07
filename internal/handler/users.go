package handler

import (
	"fmt"
	"goRent/internal/form"
	"goRent/internal/model"
	"goRent/internal/render"
	"net/http"
)

func (m *Repository) UserAccount(w http.ResponseWriter, r *http.Request) {
	data := make(map[string]interface{})

	if err := render.Template(w, r, "account.page.html", &render.TemplateData{
		Data: data,
	}); err != nil {
		m.App.Error.Println(err)
	}
}
func (m *Repository) EditUserAccount(w http.ResponseWriter, r *http.Request) {
	data := make(map[string]interface{})
	u := m.App.Session.Get(r.Context(), "user").(model.User)
	data["editUser"] = model.User{
		Username: u.Username,
		Email:    u.Email,
		Password: "",
		Address:  u.Address,
	}
	if err := render.Template(w, r, "profile.page.html", &render.TemplateData{
		Data: data,
		Form: &form.Form{},
	}); err != nil {
		m.App.Error.Println(err)
	}
}

func (m *Repository) EditUserAccountPost(w http.ResponseWriter, r *http.Request) {
	m.App.Info.Println("Register: POST")
	// u := m.App.Session.Get(r.Context(), "user").(model.User)
	// data := make(map[string]interface{})

	if err := r.ParseForm(); err != nil {
		m.App.Error.Println(err)
	}
	newUser := model.User{
		Username: r.FormValue("username"),
		Email:    r.FormValue("email"),
		Password: r.FormValue("password"),
		Address: model.Address{
			Block:      r.FormValue("block"),
			StreetName: r.FormValue("streetName"),
			UnitNumber: r.FormValue("unitNumber"),
			PostalCode: r.FormValue("postalCode"),
		},
	}

	form := form.New(r.PostForm)
	fmt.Println("FORM:", form)
	form.Required("block", "streetName", "unitNumber", "postalCode")
	form.CheckLength("block", 1, 10)
	form.CheckLength("streetName", 1, 255)
	form.CheckLength("unitNumber", 1, 10)
	form.CheckLength("postalCode", 1, 10)
	fmt.Println(newUser)
	http.Redirect(w, r, "/login", http.StatusSeeOther)

}

func (m *Repository) GetCart(w http.ResponseWriter, r *http.Request) {

	// data := make(map[string]interface{})
}

// func (m *Repository) UserAccount(w http.ResponseWriter, r *http.Request) {
// 	if f := m.App.Session.Get(r.Context(), "flash"); f != nil {
// 		m.App.Session.Put(r.Context(), "warning", nil)
// 	} else {
// 		m.App.Session.Put(r.Context(), "warning", "hello")
// 	}
// 	// m.App.Session.Put(r.Context(), "flash", "let's see")
// 	data := make(map[string]interface{})

// 	// u := m.App.Session.Get(r.Context(), "user").(model.User)
// 	// fmt.Println("PRINTINT U", u)
// 	// fmt.Println("checking authenticate", m.App.Session.Exists(r.Context(), "userID"))

// 	if err := render.Template(w, r, "user.page.html", &render.TemplateData{
// 		Data: data,
// 	}); err != nil {
// 		m.App.Error.Println(err)
// 	}
// }
