package handler

import (
	"fmt"
	"goRent/internal/config"
	"goRent/internal/form"
	"goRent/internal/model"
	"goRent/internal/render"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func (m *Repository) ShowProductByID(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	productID, err := strconv.Atoi(params["productID"])
	if err != nil {
		m.App.Error.Println(err)
		return
	}
	p, err := m.DB.GetProductByID(productID)
	if err != nil {
		m.App.Error.Println(err)
		return
	}
	rents, err := m.DB.GetRentsByProductID(productID)
	if err != nil {
		m.App.Error.Println(err)
		return
	}
	blockedDates := listBlockedDates(rents)
	data := make(map[string]interface{})
	data["product"] = p
	data["blocked"] = blockedDates
	if err := render.Template(w, r, "product.page.html", &render.TemplateData{
		Data: data,
	}); err != nil {
		m.App.Error.Println(err)
	}
}

func listBlockedDates(rents []model.Rent) []string {

	blockedDates := []string{}
	for _, r := range rents {
		x := r.StartDate
		start := r.StartDate.Format(config.DateLayout)
		end := r.EndDate.Format(config.DateLayout)
		blockedDates = append(blockedDates, start)
		for start != end {
			x = x.AddDate(0, 0, 1)
			start = x.Format(config.DateLayout)
			blockedDates = append(blockedDates, start)
			fmt.Println(x)
			fmt.Println(start)
		}
	}
	return blockedDates
}

func (m *Repository) PostReview(w http.ResponseWriter, r *http.Request) {

	u := m.App.Session.Get(r.Context(), "user").(model.User)
	productID, err := strconv.Atoi(mux.Vars(r)["productID"])
	if err != nil {
		m.App.Error.Println(err)
		return
	}

	if err := r.ParseForm(); err != nil {
		m.App.Error.Println(err)
		render.ServerError(w, r, err)
		return
	}

	form := form.New(r.Form)
	form.Required("body")

	form.CheckLength("body", 0, 500)
	rating, err := strconv.Atoi(r.PostFormValue("rating"))
	if err != nil {
		m.App.Error.Println(err)
		return
	}

	pr := model.ProductReview{
		ReviewerID:   u.ID,
		ReviewerName: u.Username,
		ProductID:    productID,
		Body:         r.PostFormValue("body"),
		Rating:       float32(rating),
	}

	if err := m.DB.AddProductReview(pr); err != nil {
		m.App.Error.Println(err)
	}

	m.App.Session.Put(r.Context(), "flash", "You have posted a review!")

	http.Redirect(w, r, fmt.Sprintf("/v1/products/%d", pr.ProductID), http.StatusSeeOther)
}
