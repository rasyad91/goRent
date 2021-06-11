package handler

import (
	"fmt"
	"goRent/internal/model"
	"goRent/internal/render"
	"net/http"
	"sync"

	"golang.org/x/sync/errgroup"
)

func (m *Repository) GetCart(w http.ResponseWriter, r *http.Request) {

	data := make(map[string]interface{})
	u := m.App.Session.Get(r.Context(), "user").(model.User)
	for _, v := range u.Rents {
		fmt.Println(v)
	}

	if err := render.Template(w, r, "cart.page.html", &render.TemplateData{
		Data: data,
	}); err != nil {
		m.App.Error.Println(err)
	}
}

func (m *Repository) GetCheckout(w http.ResponseWriter, r *http.Request) {

	data := make(map[string]interface{})

	if err := render.Template(w, r, "checkout.page.html", &render.TemplateData{
		Data: data,
	}); err != nil {
		m.App.Error.Println(err)
	}
}

func (m *Repository) PostCheckout(w http.ResponseWriter, r *http.Request) {

	// update rent
	u := m.App.Session.Get(r.Context(), "user").(model.User)

	g, _ := errgroup.WithContext(r.Context())
	var sm sync.Mutex
	failedRents := []model.Rent{}
	passedRents := []model.Rent{}

	rents := u.Rents

	// lock this
	sm.Lock()
	defer sm.Unlock()
	{
		defer sm.Unlock()

		for i, v := range rents {
			i, v := i, v
			// get product id from rents
			if !v.Processed {
				fmt.Println("RENT ID: ", v.ID)
				fmt.Println("PROCESSED: ", v.Processed)
				g.Go(func() error {
					fmt.Println("in processing------")
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

	for _, v := range rents {
		fmt.Println(v)
	}

	u.Rents = rents

	// send email
	m.App.Session.Put(r.Context(), "user", u)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
