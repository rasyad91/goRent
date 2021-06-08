package handler

import (
	"fmt"
	"goRent/internal/config"
	"goRent/internal/helper"
	"goRent/internal/model"
	"goRent/internal/render"
	"net/http"
	"strconv"
	"strings"
	"time"
)

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
	rentDates, err := helper.ListDates(start, end)
	if err != nil {
		m.App.Error.Println(err)
		return
	}

	if helper.Includes(blockedDates, rentDates...) {
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

func (m *Repository) DeleteRent(w http.ResponseWriter, r *http.Request) {

	m.App.Info.Println("deleteRent")

	if err := r.ParseForm(); err != nil {
		m.App.Error.Println(err)
		render.ServerError(w, r, err)
		return
	}

	rentID, err := strconv.Atoi(r.PostFormValue("rent_id"))
	if err != nil {
		m.App.Error.Println(err)
		render.ServerError(w, r, err)
		return
	}
	fmt.Println(rentID)

	if err := m.DB.DeleteRent(rentID); err != nil {
		render.ServerError(w, r, err)
		m.App.Error.Println(err)
		return
	}

	u := m.App.Session.Get(r.Context(), "user").(model.User)
	eu, _ := m.DB.GetUser(u.Username)
	m.App.Session.Put(r.Context(), "user", eu)

	m.App.Session.Put(r.Context(), "flash", fmt.Sprintf("Rent #%d removed from cart!", rentID))
	http.Redirect(w, r, "/v1/user/cart", http.StatusSeeOther)
}
