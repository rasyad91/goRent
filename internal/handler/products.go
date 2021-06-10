package handler

import (
	"fmt"
	"goRent/internal/form"
	"goRent/internal/helper"
	"goRent/internal/model"
	"goRent/internal/render"
	"net/http"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/gorilla/mux"

	"golang.org/x/sync/errgroup"

	// "github.com/aws/aws-sdk-go/aws"
	awsS3 "github.com/aws/aws-sdk-go/aws/session"
)

func (m *Repository) ShowProductByID(w http.ResponseWriter, r *http.Request) {

	t := time.Now()
	fmt.Println("start timing...")
	m.App.Info.Println("showProduct")
	params := mux.Vars(r)
	productID, err := strconv.Atoi(params["productID"])
	if err != nil {
		render.ServerError(w, r, err)
		m.App.Error.Println(err)
		return
	}

	g, ctx := errgroup.WithContext(r.Context())
	p := model.Product{}
	dates := []string{}
	rents := []model.Rent{}

	g.Go(func() error {
		p, err = m.DB.GetProductByID(ctx, productID)
		if err != nil {
			return err
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			return nil
		}
	})

	g.Go(func() error {
		rents, err = m.DB.GetRentsByProductID(ctx, productID)
		if err != nil {
			return err
		}
		dates = helper.ListDatesFromRents(rents)
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			return nil
		}
	})

	if err := g.Wait(); err != nil {
		render.ServerError(w, r, err)
		m.App.Error.Println(err)
		return
	}

	if helper.IsAuthenticated(r) {
		user := m.App.Session.Get(r.Context(), "user").(model.User)
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
func (m *Repository) UserProducts(w http.ResponseWriter, r *http.Request) {

	data := make(map[string]interface{})
	user := m.App.Session.Get(r.Context(), "user").(model.User)
	data["products"] = user.Products
	data["user"] = user
	if err := render.Template(w, r, "userProduct.page.html", &render.TemplateData{
		Data: data,
	}); err != nil {
		m.App.Error.Println(err)
	}
}
func (m *Repository) AddProduct(w http.ResponseWriter, r *http.Request) {

	data := make(map[string]interface{})
	// user := m.App.Session.Get(r.Context(), "user").(model.User)
	// data["products"] = user.Products
	// data["user"] = user
	if err := render.Template(w, r, "addproduct.page.html", &render.TemplateData{
		Data: data,
		Form: &form.Form{},
	}); err != nil {
		m.App.Error.Println(err)
	}
}

func (m *Repository) CreateProduct(w http.ResponseWriter, r *http.Request) {

	// data := make(map[string]interface{})

	form := form.New(r.PostForm)

	productIndex, err := m.DB.GetProductNextIndex()

	if err != nil {
		m.App.Error.Println("error occured when retriving product ID from DB query", err)
	}

	g, ctx := errgroup.WithContext(r.Context())

	for i := 1; i < 5; i++ {
		id := i
		g.Go(func() error {

			fileName := "file" + strconv.Itoa(id) //file1/2/3/4/
			file, header, err := r.FormFile(fileName)
			if err != nil || header.Size == 0 {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return err
			} else {
				defer file.Close()
				storeImagesS3(w, r, id, productIndex, m.App.AWSS3Session)
			}
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
				return nil
			}
		})
	}

	if err := g.Wait(); err != nil {
		// render.ServerError(w, r, err)
		m.App.Error.Println(err)
		form.Errors.Add("fileupload", "A miniumum of 4 images are required")

	}

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
	fmt.Println("routine ended")

	if len(form.Errors) != 0 {
		if err := render.Template(w, r, "addproduct.page.html", &render.TemplateData{
			Form: form,
			// Data: data,
		}); err != nil {
			m.App.Error.Println(err)
		}
		return
	}

}

func storeImagesS3(w http.ResponseWriter, r *http.Request, i, productIndex int, sess *awsS3.Session) {

	const MAX_UPLOAD_SIZE = 1024 * 1024 // 1MB
	uploader := s3manager.NewUploader(sess)

	r.Body = http.MaxBytesReader(w, r.Body, MAX_UPLOAD_SIZE)

	if err := r.ParseMultipartForm(MAX_UPLOAD_SIZE); err != nil {
		http.Error(w, "The uploaded file is too big. Please choose an file that's less than 1MB in size", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	fileName := "file" + strconv.Itoa(i) //file1

	file, header, err := r.FormFile(fileName)
	if err != nil || header.Size == 0 {
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
		return
	}
	fmt.Println("upload to S3 bucket was successful; please check")

	//return amz link:
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

}

// func checkEmptyUpload(w http.ResponseWriter, r *http.Request, i int, ch chan<- error, wg *sync.WaitGroup) {

// }
