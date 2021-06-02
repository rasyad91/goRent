package handler

import (
	"fmt"
	"goRent/internal/render"
	"html"
	"net/http"
	"strings"
)

func (m *Repository) Search(w http.ResponseWriter, r *http.Request) {

	m.CreateProductList()

	if err := render.Template(w, r, "home.page.html", &render.TemplateData{
		Data: nil,
	}); err != nil {
		m.App.Error.Println(err)
	}

}

func (m *Repository) SearchResult(w http.ResponseWriter, r *http.Request) {
	x := r.URL.Query()
	fmt.Println(x)
	searchKW := r.FormValue("searchtext")

	sr := html.EscapeString(searchKW)
	srtL := strings.ToLower((sr))
	_ = srtL
	// m.CreateProductList()

	// if r.Method == http.MethodPost {
	// 	searchKW := r.FormValue("searchtext")

	// 	sr := html.EscapeString(searchKW)
	// 	srtL := strings.ToLower((sr))
	// 	_ = srtL

	// 	//can we store this value into middleware so the value can be passed to the next page?

	// 	http.Redirect(w, r, "/searchresult", http.StatusSeeOther)
	// }

	if err := render.Template(w, r, "searchresult.page.html", &render.TemplateData{
		Data: nil,
	}); err != nil {
		m.App.Error.Println(err)
	}

}
