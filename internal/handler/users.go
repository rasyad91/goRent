package handler

import (
<<<<<<< HEAD
	"fmt"
=======
>>>>>>> 8612e7db05bc5f70736e568200c36b80d8984e92
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

func (m *Repository) GetCart(w http.ResponseWriter, r *http.Request) {

	data := make(map[string]interface{})

<<<<<<< HEAD
func (m *Repository) UserRents(w http.ResponseWriter, r *http.Request) {
	fmt.Println("x")
}
func (m *Repository) UserAccount(w http.ResponseWriter, r *http.Request) {
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

	if err := render.Template(w, r, "user.page.html", &render.TemplateData{
=======
	if err := render.Template(w, r, "cart.page.html", &render.TemplateData{
>>>>>>> 8612e7db05bc5f70736e568200c36b80d8984e92
		Data: data,
	}); err != nil {
		m.App.Error.Println(err)
	}
}
