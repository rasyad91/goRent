package handler

import (
	"fmt"
	"goRent/internal/form"
	"goRent/internal/model"
	"goRent/internal/render"
	"net/http"
)

func (m *Repository) AdminAccount(w http.ResponseWriter, r *http.Request) {
	data := make(map[string]interface{})
	u := m.App.Session.Get(r.Context(), "user").(model.User)
	data["user"] = u
	if u.AccessLevel != 1 {
		m.App.Session.Put(r.Context(), "warning", fmt.Sprintf("Sorry! You do not have access to this!"))
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
	result, _ := m.DB.GetAllUsers()
	data["AllUsers"] = result
	if err := render.Template(w, r, "adminUsers.page.html", &render.TemplateData{
		Data: data,
		Form: &form.Form{},
	}); err != nil {
		m.App.Error.Println(err)
	}

}

func (m *Repository) AdminAccountPost(w http.ResponseWriter, r *http.Request) {
	fmt.Println("HITTING ADMIN POST")
	u := m.App.Session.Get(r.Context(), "user").(model.User)
	if u.AccessLevel != 1 {
		m.App.Session.Put(r.Context(), "warning", fmt.Sprintf("Sorry! You do not have access to this!"))
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
	if err := r.ParseForm(); err != nil {
		m.App.Error.Println(err)
	}
	data := make(map[string]interface{})
	// form := form.New(r.PostForm)

	if err := r.ParseForm(); err != nil {
		m.App.Error.Println(err)
	}
	action := r.FormValue("action")
	fmt.Println("Action on form is", action)
	if action == "accessGrant" {
		userID := r.FormValue("userid")
		err := m.DB.GrantAccess(userID)
		fmt.Println(err)

	} else {

	}
	if action == "accessGrant" {
		m.App.Session.Put(r.Context(), "flash", "Address Updated!")
	}
	result, _ := m.DB.GetAllUsers()
	data["AllUsers"] = result
	http.Redirect(w, r, "/admin/overview", http.StatusSeeOther)

}
