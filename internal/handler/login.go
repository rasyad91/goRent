package handler

import (
	"goRent/internal/form"
	"goRent/internal/render"
	"net/http"
)

func (m *Repository) Login(w http.ResponseWriter, r *http.Request) {
	data := make(map[string]interface{})
	if err := render.Template(w, r, "login.page.html", &render.TemplateData{
		Data: data,
		Form: &form.Form{},
	}); err != nil {
		m.App.Error.Println(err)
	}
}

func (m *Repository) LoginPost(w http.ResponseWriter, r *http.Request) {
	data := make(map[string]interface{})

	if err := r.ParseForm(); err != nil {
		m.App.Error.Println(err)
		return
	}
	//...

	if err := render.Template(w, r, "login.page.html", &render.TemplateData{
		Data: data,
	}); err != nil {
		m.App.Error.Println(err)
	}
}
