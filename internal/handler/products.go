package handler

import (
	"fmt"
	"goRent/internal/form"
	"goRent/internal/helper"
	"goRent/internal/model"
	"goRent/internal/render"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func (m *Repository) ShowProductByID(w http.ResponseWriter, r *http.Request) {

	m.App.Info.Println("showProduct")
	params := mux.Vars(r)
	productID, err := strconv.Atoi(params["productID"])
	if err != nil {
		m.App.Error.Println(err)
		return
	}
	p, err := m.DB.GetProductByID(productID)
	if err != nil {
		m.App.Error.Println(err)
		return
	}
	rents, err := m.DB.GetRentsByProductID(productID)
	if err != nil {
		m.App.Error.Println(err)
		return
	}

	var user model.User
	dates := helper.ListDatesFromRents(rents)

	if helper.IsAuthenticated(r) {
		user = m.App.Session.Get(r.Context(), "user").(model.User)
		// append dates that are already booked and processed in system and dates that user has rent but not yet processed for that user
		dates = append(helper.ListDatesFromRents(rents), helper.ListDatesFromRents(user.Rents)...)
	}

	data := make(map[string]interface{})
	data["product"] = p
	data["blocked"] = dates
	if err := render.Template(w, r, "product.page.html", &render.TemplateData{
		Data: data,
	}); err != nil {
		m.App.Error.Println(err)
	}
}

func (m *Repository) PostReview(w http.ResponseWriter, r *http.Request) {

	m.App.Info.Println("postReview")
	u := m.App.Session.Get(r.Context(), "user").(model.User)
	productID, err := strconv.Atoi(mux.Vars(r)["productID"])
	if err != nil {
		m.App.Error.Println(err)
		return
	}

	if err := r.ParseForm(); err != nil {
		m.App.Error.Println(err)
		render.ServerError(w, r, err)
		return
	}

	form := form.New(r.Form)
	form.Required("body")

	form.CheckLength("body", 0, 500)
	rating, err := strconv.Atoi(r.PostFormValue("rating"))
	if err != nil {
		m.App.Error.Println(err)
		return
	}

	pr := model.ProductReview{
		ReviewerID:   u.ID,
		ReviewerName: u.Username,
		ProductID:    productID,
		Body:         r.PostFormValue("body"),
		Rating:       float32(rating),
	}

	// Add to mutex
	if err := m.DB.CreateProductReview(pr); err != nil {
		m.App.Error.Println(err)
	}

	m.App.Session.Put(r.Context(), "flash", "You have posted a review!")
	http.Redirect(w, r, fmt.Sprintf("/v1/products/%d", pr.ProductID), http.StatusSeeOther)
}

// func filterRents(processed bool) func (rents []model.Rent) []model.Rent {
// 	return func (rents []model.Rent) []model.Rent {

// 		retun
// 	}
// }

func (m *Repository) AddProduct(w http.ResponseWriter, r *http.Request) {

	data := make(map[string]interface{})

	if err := render.Template(w, r, "addproduct.page.html", &render.TemplateData{
		Data: data,
	}); err != nil {
		m.App.Error.Println(err)
	}
}

func (m *Repository) UploadImages(w http.ResponseWriter, r *http.Request) {

	// data := make(map[string]interface{})
	fmt.Fprintf(w, "Uploading File \n")

	r.ParseMultipartForm(7 << 10)

	file, handler, err := r.FormFile("productImage")
	if err != nil {
		m.App.Error.Println("Error retriveving file from form-data/ frontend:", err)
		return
	}

	defer file.Close()

	fmt.Printf("Uploaded File: %+v\n", handler.Filename)
	fmt.Printf("File Size: %+v\n", handler.Size)
	fmt.Printf("MINE Header: %+v\n", handler.Header)

	tempFile, err := ioutil.TempFile("temp-images", "upload-*.png")

	if err != nil {
		m.App.Error.Println("Error writing the image upload to the temp images directory", err)
	}

	defer tempFile.Close()

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		m.App.Error.Println("Error writing the image upload to the temp images directory", err)
	}

	tempFile.Write(fileBytes)

	fmt.Fprintf(w, "successfully uploaded file to server")

	// if err := render.Template(w, r, "addproduct.page.html", &render.TemplateData{
	// 	Data: data,
	// }); err != nil {
	// 	m.App.Error.Println(err)
	// }
}
