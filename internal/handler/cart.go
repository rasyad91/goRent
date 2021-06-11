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

func (m *Repository) CheckoutConfirm(w http.ResponseWriter, r *http.Request) {

	fmt.Println("Get checkoutconfirm")
	data := make(map[string]interface{})

	data["failedRents"] = []model.Rent{}
	data["passedRents"] = []model.Rent{}

	if m.App.Session.Get(r.Context(), "failedRents") != nil {
		data["failedRents"] = m.App.Session.Get(r.Context(), "failedRents").([]model.Rent)
	}
	if m.App.Session.Get(r.Context(), "passedRents") != nil {
		data["passedRents"] = m.App.Session.Get(r.Context(), "passedRents").([]model.Rent)
	}

	if err := render.Template(w, r, "checkout-confirm.page.html", &render.TemplateData{
		Data: data,
	}); err != nil {
		m.App.Error.Println(err)
	}
}

// Sample email
// msg := model.MailData{
// 	To:       reservation.Email,
// 	From:     "me@here.com",
// 	Subject:  "Reservation Confirmation",
// 	Content:  "",
// 	Template: "basic.html",
// }
// msg.Content = fmt.Sprintf(`
// 	<strong>Reservation Confirmation</strong><br>
// 	Dear Mr/Ms %s, <br>
// 	This is to confirm your reservation from %s to %s.
// `,
// 	reservation.LastName,
// 	reservation.StartDate.Format(datelayout),
// 	reservation.EndDate.Format(datelayout),
// )
// m.App.MailChan <- msg
