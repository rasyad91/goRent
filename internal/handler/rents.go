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
	"sync"
	"time"
)

func (m *Repository) PostRent(w http.ResponseWriter, r *http.Request) {

	m.App.Info.Println("postRent")
	if err := r.ParseForm(); err != nil {
		m.App.Error.Println(err)
		render.ServerError(w, r, err)
		return
	}
	productID, err := strconv.Atoi(r.PostFormValue("product_id"))
	if err != nil {
		m.App.Error.Println(err)
		return
	}

	if !helper.IsAuthenticated(r) {
		m.App.Session.Put(r.Context(), "url", fmt.Sprintf("/v1/products/%d", productID))
		m.App.Session.Put(r.Context(), "warning", "Sorry! You have to login first to make a booking.")
		m.App.Info.Println("user not logged in to make rent, redirecting to login")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	u := m.App.Session.Get(r.Context(), "user").(model.User)

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

	if ownerID == u.ID {
		m.App.Session.Put(r.Context(), "warning", "You fool! You cannot book your own product!")
		http.Redirect(w, r, fmt.Sprintf("/v1/products/%d", productID), http.StatusSeeOther)
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

	t := time.Now()
	fmt.Println("Create Rent: Start timing ...")

	c := make(chan int, 1)
	go func(rent model.Rent) {
		id, err := m.DB.CreateRent(rent)
		if err != nil {
			render.ServerError(w, r, err)
			m.App.Error.Println(err)
			return
		}
		c <- id
		close(c)
	}(rent)

	rent.ID = <-c
	u.Rents = append(u.Rents, rent)

	m.App.Session.Put(r.Context(), "user", u)
	fmt.Println("Time taken: ", time.Since(t))

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

	t := time.Now()
	fmt.Println("Delete Rent: Start timing ...")

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		if err := m.DB.DeleteRent(rentID); err != nil {
			render.ServerError(w, r, err)
			m.App.Error.Println(err)
			return
		}
		wg.Done()
	}()

	u := m.App.Session.Get(r.Context(), "user").(model.User)
	for i, v := range u.Rents {
		if v.ID == rentID {
			u.Rents = append(u.Rents[:i], u.Rents[i+1:]...)
		}
	}

	wg.Wait()
	m.App.Session.Put(r.Context(), "user", u)
	fmt.Println("Time taken: ", time.Since(t))

	m.App.Session.Put(r.Context(), "flash", fmt.Sprintf("Rent #%d removed from cart!", rentID))
	http.Redirect(w, r, "/v1/user/cart", http.StatusSeeOther)
}
