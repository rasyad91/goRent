package handler

import (
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

	g, ctx := errgroup.WithContext(r.Context())
	var mutex sync.Mutex

	// lock this
	mutex.Lock()
	defer mutex.Unlock()
	for _, v := range u.Rents {
		// get product id from rents

		// use product id to get rents on that product

		// map the rents with key = product id, value = []string of dates

		// check if any of the rent dates includes the dates in the database, based on each Rent
		// if include, abort processing for that rent

		// continue processing for others

		// update process in db and spin off go routine... should spin off go routine when processing or, when hitting the db
		if !v.Processed {
			v.Processed = true
			g.Go(func() error {

			})

		}
	}
	// update db

	// send email

	m.App.Session.Put(r.Context(), "confirm", true)
	m.App.Session.Put(r.Context(), "flash", "Your rent is completed!")

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
