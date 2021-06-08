package handler

import (
	"fmt"
	"goRent/internal/form"
	"goRent/internal/model"
	"goRent/internal/render"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func (m *Repository) UserAccount(w http.ResponseWriter, r *http.Request) {
	data := make(map[string]interface{})
	u := m.App.Session.Get(r.Context(), "user").(model.User)
	fmt.Println(u)
	data["user"] = model.User{
		Username: u.Username,
		Email:    u.Email,
		Password: "",
		Address:  u.Address,
		Products: u.Products,
		Rents:    u.Rents,
		Bookings: u.Bookings,
	}
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
	u := m.App.Session.Get(r.Context(), "user").(model.User)
	data := make(map[string]interface{})
	form := form.New(r.PostForm)
	data["editUser"] = model.User{
		Username: u.Username,
		Email:    u.Email,
		Password: "",
		Address:  u.Address,
	}
	if err := r.ParseForm(); err != nil {
		m.App.Error.Println(err)
	}
	action := r.FormValue("action")
	if action == "address" {
		form.CheckLength("block", 1, 10)
		form.CheckLength("streetName", 1, 255)
		form.CheckLength("unitNumber", 1, 10)
		form.CheckLength("postalCode", 1, 10)
		u.Address = model.Address{
			Block:      r.FormValue("block"),
			StreetName: r.FormValue("streetName"),
			UnitNumber: r.FormValue("unitNumber"),
			PostalCode: r.FormValue("postalCode"),
		}
	} else if action == "profile" {
		form.CheckLength("username", 1, 255)
		form.CheckLength("email", 1, 255)
		form.CheckEmail("email")
		u.Username = r.FormValue("username")
		u.Email = r.FormValue("email")
	} else {
		oldPassword := r.FormValue("password_old")
		err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(oldPassword))
		if err != nil {
			fmt.Println("OLD PASSWORD DOESN'T MATCH")
			form.Errors.Add("password_old", "Incorrect password!")
		}
		newPassword_1, newPassword_2 := r.FormValue("password_1"), r.FormValue("password_2")
		if newPassword_1 != newPassword_2 {
			fmt.Println("NEW PASSWORDS DON'T MATCH EITHER")
			form.Errors.Add("password", "New passwords doesnt' match!")
		} else {
			bpassword, err := bcrypt.GenerateFromPassword([]byte(r.FormValue("password_1")), bcrypt.DefaultCost)
			if err != nil {
				m.App.Error.Println(err)
			}
			u.Password = string(bpassword)
		}
	}

	if len(form.Errors) != 0 {
		fmt.Println("THERE ARE ERRORS - RENDERING ERRORS HERE")
		if err := render.Template(w, r, "profile.page.html", &render.TemplateData{
			Form: form,
			Data: data,
		}); err != nil {
			m.App.Error.Println(err)
		}
		return
	}
	m.DB.EditUser(u, action)
	fmt.Println("SUCCESSFULLY TRIGGERED DB")
	if action == "address" {
		m.App.Session.Put(r.Context(), "flash", "Address Updated!")
	} else if action == "profile" {
		m.App.Session.Put(r.Context(), "flash", "Profile Name Updated!")
	} else {
		m.App.Session.Put(r.Context(), "flash", "Password Updated!")
	}
	http.Redirect(w, r, "/v1/user/account/profile", http.StatusSeeOther)

}
func (m *Repository) Payment(w http.ResponseWriter, r *http.Request) {
	data := make(map[string]interface{})

	if err := render.Template(w, r, "payment.page.html", &render.TemplateData{
		Data: data,
	}); err != nil {
		m.App.Error.Println(err)
	}
}

func (m *Repository) GetCart(w http.ResponseWriter, r *http.Request) {

	data := make(map[string]interface{})

	if err := render.Template(w, r, "cart.page.html", &render.TemplateData{
		Data: data,
	}); err != nil {
		m.App.Error.Println(err)
	}
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
