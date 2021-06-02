package handler

import (
	"goRent/internal/render"
	"html"
	"net/http"
	"strings"
)

func (m *Repository) Search(w http.ResponseWriter, r *http.Request) {

	m.CreateProductList()

	if r.Method == http.MethodPost {
		searchKW := r.FormValue("searchtext")

		sr := html.EscapeString(searchKW)
		srtL := strings.ToLower((sr))
		_ = srtL
		// updateUserLastSearch(srtL, u)
		// insertUserSearchLogs(srtL, u)

		http.Redirect(w, r, "/searchresult", http.StatusSeeOther)
	}

	if err := render.Template(w, r, "home.page.html", &render.TemplateData{
		Data: nil,
	}); err != nil {
		m.App.Error.Println(err)
	}

}
