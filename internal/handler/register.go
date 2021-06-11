package handler

import (
	"database/sql"
	"fmt"
	"goRent/internal/form"
	"goRent/internal/model"
	"goRent/internal/render"
	"net/http"
	"sync"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func (m *Repository) Register(w http.ResponseWriter, r *http.Request) {
	m.App.Info.Println("Register: GET", m.App.Session.Get(r.Context(), "user"))
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
	fmt.Println("Register post: Start timing ...")

	data := make(map[string]interface{})

	if err := r.ParseForm(); err != nil {
		m.App.Error.Println(err)
	}
	t := time.Now()

	var wg sync.WaitGroup
	wg.Add(3)

	passwordChan := make(chan []byte)

	go func() {
		bpassword, err := bcrypt.GenerateFromPassword([]byte(r.FormValue("password")), bcrypt.DefaultCost)
		if err != nil {
			m.App.Error.Println(err)
		}
		passwordChan <- bpassword
		close(passwordChan)
	}()

	fmt.Println("bcrypt Time taken: ", time.Since(t))

	form := form.New(r.PostForm)

	go func() {
		form.Required("username", "email", "password", "block", "streetName", "unitNumber", "postalCode")
		form.CheckLength("username", 1, 255)
		form.CheckLength("password", 8, -1)
		form.CheckLength("email", 1, 255)
		form.CheckLength("block", 1, 10)
		form.CheckLength("streetName", 1, 255)
		form.CheckLength("unitNumber", 1, 10)
		form.CheckLength("postalCode", 1, 10)
		form.CheckEmail("email")
		wg.Done()
	}()
	fmt.Println("form checks Time taken: ", time.Since(t))

	go func(u string) {
		if _, err := m.DB.GetUser(u); err != nil {
			if err != sql.ErrNoRows {
				m.App.Error.Println(err)
			}
		} else {
			m.App.Info.Println("YES THIS USERNAME IS ALREADY IN USE")
			form.Errors.Add("username", "Username already in use")
		}
		wg.Done()
	}(r.PostFormValue("username"))

	go func(u string) {
		if err := m.DB.EmailExist(u); err != nil {
			m.App.Info.Println("YES THIS EMAIL IS ALREADY IN USE")
			form.Errors.Add("email", "Email already in use")
		}
		wg.Done()
	}(r.PostFormValue("email"))

	fmt.Println("query db for existing Time taken: ", time.Since(t))

	newUser := model.User{
		Username: r.FormValue("username"),
		Email:    r.FormValue("email"),
		Password: string(<-passwordChan),
		Address: model.Address{
			Block:      r.FormValue("block"),
			StreetName: r.FormValue("streetName"),
			UnitNumber: r.FormValue("unitNumber"),
			PostalCode: r.FormValue("postalCode"),
		},
		Image_URL: "https://cdn.pixabay.com/photo/2015/10/05/22/37/blank-profile-picture-973460_960_720.png",
	}

	fmt.Println("after get user Time taken: ", time.Since(t))

	wg.Wait()

	data["register"] = newUser

	fmt.Println(form.Errors)
	fmt.Println("last Time taken: ", time.Since(t))
	if len(form.Errors) != 0 {
		if err := render.Template(w, r, "register.page.html", &render.TemplateData{
			Form: form,
			Data: data,
		}); err != nil {
			m.App.Error.Println(err)
		}
		return
	}
	fmt.Println("CREATING NEW USER:", newUser)
	if err := m.DB.InsertUser(newUser); err != nil {
		fmt.Println(err)

	} else {
		m.App.Info.Println("Successfully Registered!")
	}
	m.App.Session.Put(r.Context(), "flash", "You've registered successfully!")

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
