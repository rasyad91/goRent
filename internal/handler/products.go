package handler

import (
	"fmt"
	"goRent/internal/config"
	"goRent/internal/form"
	"goRent/internal/model"
	"goRent/internal/render"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

func (m *Repository) ShowProductByID(w http.ResponseWriter, r *http.Request) {

	m.App.Info.Println("showProduct")

	user := m.App.Session.Get(r.Context(), "user").(model.User)
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

	// append dates that are already booked and processed in system and dates that user has rent but not yet processed for that user
	dates := append(listDatesFromRents(rents), listDatesFromRents(user.Rents)...)

	data := make(map[string]interface{})
	data["product"] = p
	data["blocked"] = dates
	if err := render.Template(w, r, "product.page.html", &render.TemplateData{
		Data: data,
	}); err != nil {
		m.App.Error.Println(err)
	}
}

func (m *Repository) PostReview(w http.ResponseWriter, r *http.Request) {

	m.App.Info.Println("postReview")
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

	if err := m.DB.CreateProductReview(pr); err != nil {
		m.App.Error.Println(err)
	}

	m.App.Session.Put(r.Context(), "flash", "You have posted a review!")
	http.Redirect(w, r, fmt.Sprintf("/v1/products/%d", pr.ProductID), http.StatusSeeOther)
}

func (m *Repository) PostRent(w http.ResponseWriter, r *http.Request) {

	m.App.Info.Println("postRent")
	u := m.App.Session.Get(r.Context(), "user").(model.User)

	if err := r.ParseForm(); err != nil {
		m.App.Error.Println(err)
		render.ServerError(w, r, err)
		return
	}

	start := r.PostFormValue("start_date")
	startDate, err := time.Parse(config.DateLayout, start)
	if err != nil {
		m.App.Error.Println(err)
		return
	}
	end := r.PostFormValue("end_date")
	endDate, err := time.Parse(config.DateLayout, end)
	if err != nil {
		m.App.Error.Println(err)
		return
	}
	productTitle := r.PostFormValue("product_title")
	blocked := r.PostFormValue("blocked")

	productID, err := strconv.Atoi(r.PostFormValue("product_id"))
	if err != nil {
		m.App.Error.Println(err)
		return
	}
	ownerID, err := strconv.Atoi(r.PostFormValue("owner_id"))
	if err != nil {
		m.App.Error.Println(err)
		return
	}
	price, err := strconv.ParseFloat(r.PostFormValue("price"), 32)
	if err != nil {
		m.App.Error.Println(err)
		return
	}

	blockedDates := strings.Split(strings.Trim(strings.Trim(blocked, "["), "]"), " ")
	rentDates, err := listDates(start, end)
	if err != nil {
		m.App.Error.Println(err)
		return
	}

	if includes(blockedDates, rentDates...) {
		m.App.Session.Put(r.Context(), "warning", "Dates selected are already booked! Please select another date")
		http.Redirect(w, r, fmt.Sprintf("/v1/products/%d", productID), http.StatusSeeOther)
		return
	}

	totalCost := float32(len(rentDates)) * float32(price)

	rent := model.Rent{
		OwnerID:   ownerID,
		RenterID:  u.ID,
		ProductID: productID,
		TotalCost: totalCost,
		Duration:  len(rentDates),
		StartDate: startDate,
		EndDate:   endDate,
	}

	if err := m.DB.CreateRent(rent); err != nil {
		render.ServerError(w, r, err)
		m.App.Error.Println(err)
		return
	}

	eu, _ := m.DB.GetUser(u.Username)
	m.App.Session.Put(r.Context(), "user", eu)

	m.App.Session.Put(r.Context(), "flash", fmt.Sprintf("You have added %s to cart!", productTitle))
	http.Redirect(w, r, fmt.Sprintf("/v1/products/%d", productID), http.StatusSeeOther)
}

func listDatesFromRents(rents []model.Rent) []string {
	dates := []string{}
	for _, r := range rents {
		x := r.StartDate
		start := r.StartDate.Format(config.DateLayout)
		end := r.EndDate.Format(config.DateLayout)
		dates = append(dates, start)
		for start != end {
			x = x.AddDate(0, 0, 1)
			start = x.Format(config.DateLayout)
			dates = append(dates, start)
		}
	}
	return dates
}

func listDates(start, end string) ([]string, error) {
	dates := []string{}
	x, err := time.Parse(config.DateLayout, start)
	if err != nil {
		return nil, err
	}
	dates = append(dates, start)
	for start != end {
		x = x.AddDate(0, 0, 1)
		start = x.Format(config.DateLayout)
		dates = append(dates, start)
	}
	return dates, nil
}

// includes compares a slice with string/s and returns true, if the instance of the string/s is in the slice
func includes(s1 []string, s2 ...string) bool {
	for _, v1 := range s1 {
		for _, v2 := range s2 {
			if v1 == v2 {
				return true
			}
		}
	}
	return false
}

// func filterRents(processed bool) func (rents []model.Rent) []model.Rent {
// 	return func (rents []model.Rent) []model.Rent {

// 		retun
// 	}
// }
