package handler

import (
	"fmt"
	"goRent/internal/render"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func (m *Repository) ShowProductByID(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	fmt.Println("PARAMS", params)
	productId, err := strconv.Atoi(params["productId"])
	fmt.Println("productid", productId)
	if err != nil {
		m.App.Error.Println(err)
		return
	}

	p, err := m.DB.GetProductByID(productId)
	if err != nil {
		m.App.Error.Println(err)
		return
	}
	fmt.Println(p.Reviews)
	data := make(map[string]interface{})
	data["product"] = p

	if err := render.Template(w, r, "product.page.html", &render.TemplateData{
		Data: data,
	}); err != nil {
		m.App.Error.Println(err)
	}

}
