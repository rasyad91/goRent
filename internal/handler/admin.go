package handler

import (
	"fmt"
	"goRent/internal/form"
	"goRent/internal/model"
	"goRent/internal/render"
	"net/http"
	"sort"
	"strings"
)

func (m *Repository) AdminAccount(w http.ResponseWriter, r *http.Request) {
	data := make(map[string]interface{})
	u := m.App.Session.Get(r.Context(), "user").(model.User)
	data["user"] = u

	result, _ := m.DB.GetAllUsers()
	sortby, ok := r.URL.Query()["sortby"]
	sortType, sortok := r.URL.Query()["sort"]
	fmt.Println("SORTBY:", sortby)
	fmt.Println("SORT:", sortType)
	if ok && sortok {
		if sortby[0] == "username" {
			if sortType[0] == "asc" {
				sort.SliceStable(result, func(i, j int) bool {
					return strings.ToLower(result[i].Username) < strings.ToLower(result[j].Username)
				})
			} else {
				sort.SliceStable(result, func(i, j int) bool {
					return strings.ToLower(result[i].Username) > strings.ToLower(result[j].Username)
				})
			}

		} else if sortby[0] == "access" {
			if sortType[0] == "asc" {
				sort.SliceStable(result, func(i, j int) bool {
					return result[i].AccessLevel < result[j].AccessLevel
				})
			} else {
				sort.SliceStable(result, func(i, j int) bool {
					return result[i].AccessLevel > result[j].AccessLevel
				})
			}
		} else if sortby[0] == "id" {
			if sortType[0] == "asc" {
				sort.SliceStable(result, func(i, j int) bool {
					return result[i].ID < result[j].ID
				})
			} else {
				sort.SliceStable(result, func(i, j int) bool {
					return result[i].ID > result[j].ID
				})
			}
		}
	}

	data["AllUsers"] = result
	if err := render.Template(w, r, "adminUsers.page.html", &render.TemplateData{
		Data: data,
		Form: &form.Form{},
	}); err != nil {
		m.App.Error.Println(err)
	}

}

func (m *Repository) AdminAccountPost(w http.ResponseWriter, r *http.Request) {
	fmt.Println("HITTING ADMIN POST")
	u := m.App.Session.Get(r.Context(), "user").(model.User)
	if u.AccessLevel != 1 {
		m.App.Session.Put(r.Context(), "warning", fmt.Sprintf("Sorry! You do not have access to this!"))
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
	if err := r.ParseForm(); err != nil {
		m.App.Error.Println(err)
	}
	data := make(map[string]interface{})
	// form := form.New(r.PostForm)

	if err := r.ParseForm(); err != nil {
		m.App.Error.Println(err)
	}
	action := r.FormValue("action")
	fmt.Println("Action on form is", action)
	if action == "accessGrant" {
		userID := r.FormValue("userid")
		err := m.DB.GrantAccess(userID)
		if err != nil {
			m.App.Session.Put(r.Context(), "warning", "Access not granted")
		} else {
			m.App.Session.Put(r.Context(), "flash", "Access Granted!")
		}

	} else if action == "removeAccess" {
		userID := r.FormValue("userid")
		err := m.DB.RemoveAccess(userID)
		if err != nil {
			m.App.Session.Put(r.Context(), "warning", "Access not reverted!")
		} else {
			m.App.Session.Put(r.Context(), "flash", "Access removed successfully!")
		}
	} else if action == "massiveDelete" {
		userID := r.FormValue("userid")
		fmt.Println("GIGANTIC MASSIVE DELETE!!")
		err := m.DB.DeleteUser(userID)
		if err != nil {
			m.App.Session.Put(r.Context(), "warning", "User not removed!")
		} else {
			m.App.Session.Put(r.Context(), "flash", "User removed successfully!")
		}
	}
	result, _ := m.DB.GetAllUsers()
	data["AllUsers"] = result
	http.Redirect(w, r, "/admin/overview", http.StatusSeeOther)

}
func (m *Repository) AdminProducts(w http.ResponseWriter, r *http.Request) {
	data := make(map[string]interface{})
	u := m.App.Session.Get(r.Context(), "user").(model.User)
	data["user"] = u
	if u.AccessLevel != 1 {
		m.App.Session.Put(r.Context(), "warning", fmt.Sprintf("Sorry! You do not have access to this!"))
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
	result, _ := m.DB.GetAllProducts()
	data["AllProducts"] = result
	if err := render.Template(w, r, "adminProducts.page.html", &render.TemplateData{
		Data: data,
		Form: &form.Form{},
	}); err != nil {
		m.App.Error.Println(err)
	}

}

func (m *Repository) AdminRentals(w http.ResponseWriter, r *http.Request) {
	data := make(map[string]interface{})
	u := m.App.Session.Get(r.Context(), "user").(model.User)
	data["user"] = u
	if u.AccessLevel != 1 {
		m.App.Session.Put(r.Context(), "warning", fmt.Sprintf("Sorry! You do not have access to this!"))
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
	result, _ := m.DB.GetAllRents()
	data["AllRents"] = result
	fmt.Println(result)
	if err := render.Template(w, r, "adminRentals.page.html", &render.TemplateData{
		Data: data,
		Form: &form.Form{},
	}); err != nil {
		m.App.Error.Println(err)
	}

}