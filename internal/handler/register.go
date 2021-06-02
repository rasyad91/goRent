package handler

import (
	"fmt"
	"goRent/internal/form"
	"goRent/internal/model"
	"goRent/internal/render"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func (m *Repository) Register(w http.ResponseWriter, r *http.Request) {
	m.App.Info.Println("Register: GET")
	data := make(map[string]interface{})
	data["register"] = model.User{
		Username: "",
		Email:    "",
		Password: "",
		Address:  model.Address{},
	}

	if err := render.Template(w, r, "register.page.html", &render.TemplateData{
		Form: &form.Form{},
		Data: data,
	}); err != nil {
		m.App.Error.Println(err)
	}
}

func (m *Repository) RegisterPost(w http.ResponseWriter, r *http.Request) {
	m.App.Info.Println("Register: POST")

	data := make(map[string]interface{})

	if err := r.ParseForm(); err != nil {
		m.App.Error.Println(err)
	}
	bpassword, err := bcrypt.GenerateFromPassword([]byte(r.FormValue("password")), bcrypt.DefaultCost)
	if err != nil {
		m.App.Error.Println(err)
	}

	newUser := model.User{
		Username: r.FormValue("username"),
		Email:    r.FormValue("email"),
		Password: string(bpassword),
		Address: model.Address{
			Block:      r.FormValue("block"),
			StreetName: r.FormValue("streetName"),
			UnitNumber: r.FormValue("unitNumber"),
			PostalCode: r.FormValue("postalCode"),
		},
	}

	form := form.New(r.PostForm)
	fmt.Println("FORM:", form)
	form.Required("username", "email", "password", "block", "streetName", "unitNumber", "postalCode")

	if len(r.FormValue("password")) > 20 || len(r.FormValue("password")) < 8 {
		fmt.Println("YES THIS USERNAME IS ALREADY IN USE")
		form.Errors.Add("password", "Password must be more than 8 and less than 20")
	}

	// errors["password"] = {
	// 	"password is required", "password must be more ..."
	// }

	_, isExist := m.DB.GetUser(newUser.Username)
	if isExist {
		fmt.Println("YES THIS USERNAME IS ALREADY IN USE")
		form.Errors.Add("username", "Username already in use")
	}

	data["register"] = newUser
	fmt.Println(form.Errors)

	if len(form.Errors) != 0 {
		fmt.Println(form.Errors.Get("inputUsername"))
		fmt.Println("in form. errors")
		if err := render.Template(w, r, "register.page.html", &render.TemplateData{
			Form: form,
			Data: data,
		}); err != nil {
			m.App.Error.Println(err)
		}
		return
	}

	addedSuccess := m.DB.InsertUser(newUser)
	if addedSuccess {
		fmt.Println("SUCCESSFULLY REGISTERED")
	}

	m.App.Info.Println("Register: redirecting to login page")
	http.Redirect(w, r, "/login", http.StatusSeeOther)

}

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
