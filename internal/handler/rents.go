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

	"golang.org/x/sync/errgroup"
)

func (m *Repository) PostRent(w http.ResponseWriter, r *http.Request) {

	m.App.Info.Println("postRent")
	if err := r.ParseForm(); err != nil {
		m.App.Error.Println(err)
		render.ServerError(w, r, err)
		return
	}

	// if !helper.IsAuthenticated(r) {

	// 	m.App.Session.Put(r.Context(), "url", r.URL.String())
	// 	m.App.Session.Put(r.Context(), "warning", "Please login first before making a booking")
	// 	http.Redirect(w, r, "/login", http.StatusFound)
	// 	return
	// }

	productID, err := strconv.Atoi(r.PostFormValue("product_id"))
	if err != nil {
		m.App.Error.Println(err)
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
		Product:   model.Product{Title: productTitle, Price: float32(price)},
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
	for _, v := range u.Rents {
		fmt.Println(v)
	}

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

func (m *Repository) ConfirmRents(w http.ResponseWriter, r *http.Request) {

	u := m.App.Session.Get(r.Context(), "user").(model.User)

	g, _ := errgroup.WithContext(r.Context())
	var sm sync.Mutex
	failedRents := []model.Rent{}
	passedRents := []model.Rent{}

	rents := u.Rents

	// lock this
	sm.Lock()
	{
		defer sm.Unlock()

		for i, v := range rents {
			i, v := i, v
			// get product id from rents
			if !v.Processed {
				g.Go(func() error {
					fmt.Println("In processing------")
					fmt.Printf("id: %d, productname: %s, startDate: %s, endDate: %s\n", v.ID, v.Product.Title, v.StartDate, v.EndDate)
					if err := m.DB.ProcessRent(v); err != nil {
						if err.Error() == "rent not available" {
							fmt.Println("errRentNotAvailable: ", v.ID, v.Product.Title, v.StartDate, v.EndDate)
							failedRents = append(failedRents, v)
							return nil
						}
						fmt.Println("in else")
						return err

					}
					rents[i].Processed = true
					passedRents = append(passedRents, rents[i])
					fmt.Println("Pass processing------")
					fmt.Printf("id: %d, processed: %t productname: %s, startDate: %s, endDate: %s\n", v.ID, rents[i].Processed, v.Product.Title, v.StartDate, v.EndDate)
					return nil
				})
			}
		}

		if err := g.Wait(); err != nil {
			fmt.Println("in g.wait")
			fmt.Println(err)
			return
		}

		var wg sync.WaitGroup
		for _, rent := range failedRents {
			// need to remove rent from db and remove rent from user
			rent := rent

			fmt.Printf("id: %d, processed: %t productname: %s, startDate: %s, endDate: %s\n", rent.ID, rent.Processed, rent.Product.Title, rent.StartDate, rent.EndDate)
			wg.Add(1)
			go func() {
				if err := m.DB.DeleteRent(rent.ID); err != nil {
					render.ServerError(w, r, err)
					m.App.Error.Println(err)
					return
				}
				wg.Done()
			}()
		}
		wg.Wait()

		fmt.Println("length of failed rents: ", len(failedRents))
		for _, rent := range failedRents {
			fmt.Println()
			fmt.Println("FAILED RENTS")
			fmt.Printf("id: %d, processed: %t productname: %s, startDate: %s, endDate: %s\n", rent.ID, rent.Processed, rent.Product.Title, rent.StartDate, rent.EndDate)

		}

		fmt.Println("length of passed rents: ", len(passedRents))
		for _, rent := range passedRents {
			fmt.Println()
			fmt.Println("PASSED RENTS")
			fmt.Printf("id: %d, processed: %t productname: %s, startDate: %s, endDate: %s\n", rent.ID, rent.Processed, rent.Product.Title, rent.StartDate, rent.EndDate)

		}
	}

	u, err := m.DB.GetUser(u.Username)
	if err != nil {
		render.ServerError(w, r, err)
		m.App.Error.Println(err)
		return
	}

	// send email
	m.App.Session.Put(r.Context(), "user", u)
	m.App.Session.Put(r.Context(), "passedRents", passedRents)
	m.App.Session.Put(r.Context(), "failedRents", failedRents)

	http.Redirect(w, r, "/v1/user/cart/checkout/confirm", http.StatusSeeOther)
}
