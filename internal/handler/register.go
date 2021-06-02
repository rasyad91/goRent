package handler

import (
	"fmt"
	"goRent/internal/form"
	"goRent/internal/model"
	"goRent/internal/render"
	"net/http"
)

func (m *Repository) Register(w http.ResponseWriter, r *http.Request) {
	m.App.Info.Println("Register: GET")
	data := make(map[string]interface{})

	if err := render.Template(w, r, "register.page.html", &render.TemplateData{
		Data: data,
	}); err != nil {
		m.App.Error.Println(err)
	}
}

func (m *Repository) RegisterPost(w http.ResponseWriter, r *http.Request) {
	m.App.Info.Println("Register: POST")

	data := make(map[string]interface{})

	if r.Method == http.MethodPost {
		if err := r.ParseForm(); err != nil {
			m.App.Error.Println(err)
		}

		newUser := model.User{
			Username: r.FormValue("inputUsername"),
			Email:    r.FormValue("inputEmail"),
			Password: r.FormValue("inputPassword"),
			Address: model.Address{
				Block:      r.FormValue("addressblock"),
				StreetName: r.FormValue("inputAddress"),
				UnitNumber: r.FormValue("addressunit"),
				PostalCode: r.FormValue("postalcode"),
			},
		}
		fmt.Println("PARSEFORM:", r.Form["username"])
		fmt.Println("PARSEFORM:", r.Form)
		form := form.New(r.PostForm)
		fmt.Println("FORM:", form)
		form.Required("inputUsername", "inputEmail", "inputEmail", "inputPassword", "addressblock", "inputAddress", "addressunit", "postalcode")

		_, isExist := m.DB.GetUser(newUser.Username)
		if isExist {
			form.Errors.Add("username", "Username already in use")
		}
		fmt.Println("REGISTER:", newUser)

		//...
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
		m.App.Info.Println("Register: redirecting to login page")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	if err := render.Template(w, r, "register.page.html", &render.TemplateData{
		Data: data,
	}); err != nil {
		m.App.Error.Println(err)
	}
}
