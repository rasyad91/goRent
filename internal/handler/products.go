package handler

import (
	"fmt"
	"goRent/internal/form"
	"goRent/internal/helper"
	"goRent/internal/model"
	"goRent/internal/render"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

func (m *Repository) ShowProductByID(w http.ResponseWriter, r *http.Request) {

	t := time.Now()
	fmt.Println("start timing...")
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
	fmt.Println("time taken", time.Since(t))
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
	u := m.App.Session.Get(r.Context(), "user").(model.User)
	data["products"] = u.Products
	if err := render.Template(w, r, "userProduct.page.html", &render.TemplateData{
		Data: data,
	}); err != nil {
		m.App.Error.Println(err)
	}
}

func (m *Repository) CreateProduct(w http.ResponseWriter, r *http.Request) {

	// data := make(map[string]interface{})
	for i := 1; i < 5; i++ {
		storeImages(w, r, i)
	}
	fmt.Fprintf(w, "successfully uploaded file to server")

}

func storeImages(w http.ResponseWriter, r *http.Request, i int) {

	const MAX_UPLOAD_SIZE = 1024 * 1024 // 1MB

	r.Body = http.MaxBytesReader(w, r.Body, MAX_UPLOAD_SIZE)
	if err := r.ParseMultipartForm(MAX_UPLOAD_SIZE); err != nil {
		http.Error(w, "The uploaded file is too big. Please choose an file that's less than 1MB in size", http.StatusBadRequest)
		return
	}

	fileName := "file" + strconv.Itoa(i)
	file, fileHeader, err := r.FormFile(fileName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	defer file.Close()

	buff := make([]byte, 512)
	_, err = file.Read(buff)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	filetype := http.DetectContentType(buff)
	if filetype != "image/jpeg" && filetype != "image/png" {
		http.Error(w, "The provided file format is not allowed. Please upload a JPEG or PNG image", http.StatusBadRequest)
		return
	}

	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = os.MkdirAll("./uploads", os.ModePerm)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Create a new file in the uploads directory & replace product ID with real productID
	var productID string = "1" //assume productID is 1.
	imageFileName := productID + "_" + strconv.Itoa(i)
	// dst, err := os.Create(fmt.Sprintf("./uploads/%d%s", time.Now().UnixNano(), filepath.Ext(fileHeader.Filename)))
	dst, err := os.Create(fmt.Sprintf("./uploads/%s%s", imageFileName, filepath.Ext(fileHeader.Filename)))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	defer dst.Close()

	// Copy the uploaded file to the filesystem
	// at the specified destination
	_, err = io.Copy(dst, file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}
