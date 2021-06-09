package handler

import (
	"fmt"
	"goRent/internal/form"
	"goRent/internal/helper"
	"goRent/internal/model"
	"goRent/internal/render"
	"net/http"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/gorilla/mux"

	// "github.com/aws/aws-sdk-go/aws"
	awsS3 "github.com/aws/aws-sdk-go/aws/session"
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

	ch := make(chan string)

	productIndex, err := m.DB.GetProductNextIndex()

	if err != nil {
		m.App.Error.Println("error occured when retriving product ID from DB query", err)
	}

	for i := 1; i < 5; i++ {
		// storeImages(w, r, i)
		go storeImagesS3(w, r, i, productIndex, m.App.AWSS3Session, ch)
	}

	form := form.New(r.PostForm)
	form.Required("productname", "price", "brand", "productdescription")
	form.CheckLength("productname", 1, 255)
	form.CheckLength("price", 1, 5)
	form.CheckLength("productdescription", 1, 400)

	productname := r.FormValue("productname")
	price := r.FormValue("price")
	brand := r.FormValue("brand")
	productdescription := r.FormValue("productdescription")
	category := r.FormValue("category")

	fmt.Println("product name", productname)
	fmt.Println("price", price)
	fmt.Println("brand", brand)
	fmt.Println("productdescription", productdescription)
	fmt.Println("category", category)

	fmt.Fprintf(w, "successfully uploaded file to server")
	var imagelinks []string
	for i := range ch {
		imagelinks = append(imagelinks, i)
		fmt.Println("s3 stored URL", i)
	}

	//explore Selectcase for goroutine

}

func storeImagesS3(w http.ResponseWriter, r *http.Request, i, productIndex int, sess *awsS3.Session, ch chan string) {

	const MAX_UPLOAD_SIZE = 1024 * 1024 // 1MB
	uploader := s3manager.NewUploader(sess)

	r.Body = http.MaxBytesReader(w, r.Body, MAX_UPLOAD_SIZE)

	if err := r.ParseMultipartForm(MAX_UPLOAD_SIZE); err != nil {
		http.Error(w, "The uploaded file is too big. Please choose an file that's less than 1MB in size", http.StatusBadRequest)
		return
	}

	fileName := "file" + strconv.Itoa(i) //file1

	file, fileHeader, err := r.FormFile(fileName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	_ = fileHeader

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

	var s3fileExtension string
	if filetype == "image/jpeg" {
		s3fileExtension = ".jpeg"
	} else {
		s3fileExtension = ".png"
	}

	s3FileName := strconv.Itoa(productIndex) + "_" + strconv.Itoa(i) + s3fileExtension

	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String("wooteam-productslist/product_list/images/"),
		ACL:    aws.String("public-read"),
		Key:    aws.String(s3FileName),
		Body:   file,
	})

	if err != nil {
		fmt.Println("error with uploading file", err)
	} else {
		fmt.Println("upload to S3 bucket was successful; please check")
	}

	//return amz link:
	s3link := "https://wooteam-productslist.s3-ap-southeast-1.amazonaws.com/product_list/images/" + s3FileName
	ch <- s3link
}

func storeProfileImage(w http.ResponseWriter, r *http.Request, owner_ID int, sess *awsS3.Session) {

	const MAX_UPLOAD_SIZE = 1024 * 1024 // 1MB
	uploader := s3manager.NewUploader(sess)

	r.Body = http.MaxBytesReader(w, r.Body, MAX_UPLOAD_SIZE)

	if err := r.ParseMultipartForm(MAX_UPLOAD_SIZE); err != nil {
		http.Error(w, "The uploaded file is too big. Please choose an file that's less than 1MB in size", http.StatusBadRequest)
		return
	}

	fileName := "file" + strconv.Itoa(owner_ID) //file1

	file, fileHeader, err := r.FormFile(fileName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	_ = fileHeader

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

	var s3fileExtension string
	if filetype == "image/jpeg" {
		s3fileExtension = ".jpeg"
	} else {
		s3fileExtension = ".png"
	}

	s3FileName := strconv.Itoa(owner_ID) + s3fileExtension

	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String("wooteam-productslist/profile_images/"),
		ACL:    aws.String("public-read"),
		Key:    aws.String(s3FileName),
		Body:   file,
	})

	if err != nil {
		fmt.Println("error with uploading file", err)
	} else {
		fmt.Println("upload to S3 bucket was successful; please check")
	}

	// s3link := "https://wooteam-productslist.s3-ap-southeast-1.amazonaws.com/product_list/images/" + s3FileName
}
