package handler

import (
	"fmt"
	"goRent/internal/model"
	"goRent/internal/render"
	"net/http"
)

func (m *Repository) GetCart(w http.ResponseWriter, r *http.Request) {

	data := make(map[string]interface{})
	u := m.App.Session.Get(r.Context(), "user").(model.User)
	fmt.Println(u.Rents[0].Product.Images)

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

func (m *Repository) CheckoutConfirm(w http.ResponseWriter, r *http.Request) {

	fmt.Println("Get checkoutconfirm")
	u := m.App.Session.Get(r.Context(), "user").(model.User)
	data := make(map[string]interface{})

	data["failedRents"] = []model.Rent{}
	data["passedRents"] = []model.Rent{}

	msg := model.MailData{
		To:       "rasyadsubandrio@gmail.com",
		From:     "gorent.help@gmail.com",
		Subject:  "Rent confirm",
		Content:  "",
		Template: "basic.html",
	}

	msg.Content = fmt.Sprintf(`
					Hi %s,

					Your rents have been confirmed.
					Thank you for helping with our goal to reduce the consumption footprint.
					
					GoRent
						`, u.Username,
	)

	m.App.MailChan <- msg

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
