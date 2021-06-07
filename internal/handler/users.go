package handler

import (
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
