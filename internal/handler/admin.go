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
		m.App.Session.Put(r.Context(), "warning", "Sorry! You do not have access to this!")
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
			m.App.Error.Println(err)

		} else {
			m.App.Session.Put(r.Context(), "flash", "Access Granted!")
		}

	} else if action == "removeAccess" {
		userID := r.FormValue("userid")
		err := m.DB.RemoveAccess(userID)
		if err != nil {
			m.App.Session.Put(r.Context(), "warning", "Access not reverted!")
			m.App.Error.Println(err)

		} else {
			m.App.Session.Put(r.Context(), "flash", "Access removed successfully!")
		}
	} else if action == "massiveDelete" {
		userID := r.FormValue("userid")
		fmt.Println("GIGANTIC MASSIVE DELETE!!")
		err := DeleteProductsElasticUserID(r, m.App.AWSClient, userID)
		if err != nil {
			m.App.Error.Println("Error deleting product from elastic database", err)
		}
		err = m.DB.DeleteUser(userID)
		if err != nil {
			m.App.Session.Put(r.Context(), "warning", "User not removed!")
			m.App.Error.Println(err)
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
		m.App.Session.Put(r.Context(), "warning", "Sorry! You do not have access to this!")
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
	result, _ := m.DB.GetAllProducts()

	data["AllProducts"] = result
	urlQuery(result, r)
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
		m.App.Session.Put(r.Context(), "warning", "Sorry! You do not have access to this!")
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
	result, _ := m.DB.GetAllRents()
	urlQueryRent(result, r)
	data["AllRents"] = result
	fmt.Println(result)
	if err := render.Template(w, r, "adminRentals.page.html", &render.TemplateData{
		Data: data,
		Form: &form.Form{},
	}); err != nil {
		m.App.Error.Println(err)
	}

}

func urlQuery(result []model.Product, r *http.Request) {
	sortby, ok := r.URL.Query()["sortby"]
	sortType, sortok := r.URL.Query()["sort"]
	fmt.Println("SORTBY:", sortby)
	fmt.Println("SORT:", sortType)
	if ok && sortok {
		if sortby[0] == "prodID" {
			if sortType[0] == "asc" {
				sort.SliceStable(result, func(i, j int) bool {
					return result[i].ID < result[j].ID
				})
			} else {
				sort.SliceStable(result, func(i, j int) bool {
					return result[i].ID > result[j].ID
				})
			}

		} else if sortby[0] == "ownerID" {
			if sortType[0] == "asc" {
				sort.SliceStable(result, func(i, j int) bool {
					return result[i].OwnerID < result[j].OwnerID
				})
			} else {
				sort.SliceStable(result, func(i, j int) bool {
					return result[i].OwnerID > result[j].OwnerID
				})
			}
		} else if sortby[0] == "brand" {
			if sortType[0] == "asc" {
				sort.SliceStable(result, func(i, j int) bool {
					return strings.ToLower(result[i].Brand) < strings.ToLower(result[j].Brand)
				})
			} else {
				sort.SliceStable(result, func(i, j int) bool {
					return strings.ToLower(result[i].Brand) > strings.ToLower(result[j].Brand)
				})
			}
		} else if sortby[0] == "title" {
			if sortType[0] == "asc" {
				sort.SliceStable(result, func(i, j int) bool {
					return strings.ToLower(result[i].Title) < strings.ToLower(result[j].Title)
				})
			} else {
				sort.SliceStable(result, func(i, j int) bool {
					return strings.ToLower(result[i].Title) > strings.ToLower(result[j].Title)
				})
			}
		} else if sortby[0] == "price" {
			if sortType[0] == "asc" {
				sort.SliceStable(result, func(i, j int) bool {
					return result[i].Price < result[j].Price
				})
			} else {
				sort.SliceStable(result, func(i, j int) bool {
					return result[i].Price > result[j].Price
				})
			}
		}
	}
}
func urlQueryRent(result []model.Rent, r *http.Request) {
	sortby, ok := r.URL.Query()["sortby"]
	sortType, sortok := r.URL.Query()["sort"]
	fmt.Println("SORTBY:", sortby)
	fmt.Println("SORT:", sortType)
	if ok && sortok {
		if sortby[0] == "id" {
			if sortType[0] == "asc" {
				sort.SliceStable(result, func(i, j int) bool {
					return result[i].ID < result[j].ID
				})
			} else {
				sort.SliceStable(result, func(i, j int) bool {
					return result[i].ID > result[j].ID
				})
			}

		} else if sortby[0] == "ownerid" {
			if sortType[0] == "asc" {
				sort.SliceStable(result, func(i, j int) bool {
					return result[i].OwnerID < result[j].OwnerID
				})
			} else {
				sort.SliceStable(result, func(i, j int) bool {
					return result[i].OwnerID > result[j].OwnerID
				})
			}
		} else if sortby[0] == "renterid" {
			if sortType[0] == "asc" {
				sort.SliceStable(result, func(i, j int) bool {
					return result[i].RenterID < result[j].RenterID
				})
			} else {
				sort.SliceStable(result, func(i, j int) bool {
					return result[i].RenterID > result[j].RenterID
				})
			}
		} else if sortby[0] == "prodid" {
			if sortType[0] == "asc" {
				sort.SliceStable(result, func(i, j int) bool {
					return result[i].ProductID < result[j].ProductID
				})
			} else {
				sort.SliceStable(result, func(i, j int) bool {
					return result[i].ProductID > result[j].ProductID
				})
			}
		} else if sortby[0] == "duration" {
			if sortType[0] == "asc" {
				sort.SliceStable(result, func(i, j int) bool {
					return result[i].Duration < result[j].Duration
				})
			} else {
				sort.SliceStable(result, func(i, j int) bool {
					return result[i].Duration > result[j].Duration
				})
			}
		} else if sortby[0] == "cost" {
			if sortType[0] == "asc" {
				sort.SliceStable(result, func(i, j int) bool {
					return result[i].TotalCost < result[j].TotalCost
				})
			} else {
				sort.SliceStable(result, func(i, j int) bool {
					return result[i].TotalCost > result[j].TotalCost
				})
			}
		}
	}
}
