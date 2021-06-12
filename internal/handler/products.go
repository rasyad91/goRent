package handler

import (
	"fmt"
	"goRent/internal/config"
	"goRent/internal/form"
	"goRent/internal/helper"
	"goRent/internal/model"
	"goRent/internal/render"
	"net/http"
	"net/url"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"
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
		return nil
	})

	g.Go(func() error {
		rents, err = m.DB.GetRentsByProductID(ctx, productID)
		if err != nil {
			return err
		}
		dates = helper.ListDatesFromRents(rents)
		return nil
	})

	if err := g.Wait(); err != nil {
		render.ServerError(w, r, err)
		m.App.Error.Println(err)
		return
	}
	fmt.Println(dates)
	if helper.IsAuthenticated(r) {
		user := m.App.Session.Get(r.Context(), "user").(model.User)
		// append dates that are already booked and processed in system and dates that user has rent but not yet processed for that user
		userRentofPID := []model.Rent{}
		for _, v := range user.Rents {
			if v.ProductID == productID {
				userRentofPID = append(userRentofPID, v)
			}
		}
		dates = append(helper.ListDatesFromRents(rents), helper.ListDatesFromRents(userRentofPID)...)
	}
	fmt.Println(dates)

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
	var sm sync.Mutex
	sm.Lock()
	{
		defer sm.Unlock()
		if err := m.DB.CreateProductReview(pr); err != nil {
			m.App.Error.Println(err)
		}
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

	type imageIndex struct {
		index     int
		imageType string
	}

	var (
		productCategory    string
		productPrice       float32
		s3ImageInformation []imageIndex
	)

	u := m.App.Session.Get(r.Context(), "user").(model.User)

	form := form.New(r.PostForm)

	productIndex, err := m.DB.GetProductNextIndex()

	if err != nil {
		m.App.Error.Println("error occured when retriving product ID from DB query", err)
	}

	g, ctx := errgroup.WithContext(r.Context())
	var imgCount int = 0

	for i := 1; i < 5; i++ {
		id := i
		g.Go(func() error {
			fileName := "file" + strconv.Itoa(id) //file1/2/3/4/
			file, header, err := r.FormFile(fileName)
			if err != nil || header.Size == 0 {
				// http.Error(w, err.Error(), http.StatusBadRequest)
				fmt.Println(err)
				return err
			} else {
				defer file.Close()
				// uploader := s3manager.NewUploader(m.App.AWSS3Session)
				s3imgType, s3err := storeImagesS3(w, r, id, productIndex, m.App.AWSS3Session)
				if s3err != nil {
					m.App.Error.Println("S3 error", err)
					fmt.Println("")
					form.Errors.Add("fileupload", "Please only use .jpeg/ .png files not exceeding 1MB in size")
				} else {
					imgCount++
					s3ImageInformation = append(s3ImageInformation, imageIndex{index: id, imageType: s3imgType})
				}
			}
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
				return nil
			}
		})
	}
	g.Go(func() error {

		form.Required("productname", "price", "brand", "productdescription")
		form.CheckLength("productname", 1, 255)
		form.CheckLength("price", 1, 5)
		form.CheckLength("productdescription", 1, 400)
		productCategory = form.RetrieveCategory("category")
		productPrice = form.ProcessPrice("price")

		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			return nil
		}
	})

	if err := g.Wait(); err != nil {
		m.App.Error.Println(err)
		if imgCount == 0 {
			form.Errors.Add("fileupload", "Please at least upload one image")
		}
		// form.Errors.Add("fileupload", "A miniumum of 4 images are required")
	}

	// sort.Slice(imageIndex)
	fmt.Println("this is the unsorted index", s3ImageInformation)

	sort.Slice(s3ImageInformation, func(i, j int) bool {
		return s3ImageInformation[i].index < s3ImageInformation[j].index
	})

	fmt.Println("this is the sorted index", s3ImageInformation)

	var productImageURL []string
	for _, v := range s3ImageInformation {

		s := config.AWSProductImageLink
		productImageURL = append(productImageURL, s+strconv.Itoa(productIndex)+"_"+strconv.Itoa(v.index)+v.imageType)

	}

	var newProduct = model.Product{

		ID:          productIndex,
		OwnerID:     u.ID,
		OwnerName:   u.Username,
		Brand:       r.FormValue("brand"),
		Category:    productCategory,
		Title:       r.FormValue("productname"),
		Rating:      0,
		Description: r.FormValue("productdescription"),
		Price:       productPrice,
		Reviews:     []model.ProductReview{},
		Images:      productImageURL,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if len(form.Errors) != 0 {
		if err := render.Template(w, r, "addproduct.page.html", &render.TemplateData{
			Form: form,
		}); err != nil {
			m.App.Error.Println(err)
		}
		return
	} else {

		g2, ctx := errgroup.WithContext(r.Context())

		g2.Go(func() error {
			err := m.DB.InsertProduct(newProduct)
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

		// index for elastic search

		g2.Go(func() error {
			put1, err := m.App.AWSClient.Index().
				Index("sample_product_list").
				Type("sampleproducttype").
				Id(strconv.Itoa(newProduct.ID)).
				BodyJson(newProduct).
				Do(r.Context())

			if err != nil {
				// Handle error
				panic(err)
			}
			fmt.Printf("Indexed tweet %s to index %s, type %s\n", put1.Id, put1.Index, put1.Type)

			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
				return nil
			}
		})

		if err := g2.Wait(); err != nil {
			m.App.Error.Println("error from g2", err)
			// form.Errors.Add("fileupload", "A miniumum of 4 images are required")
		}

		m.App.Session.Put(r.Context(), "flash", "You've successfully created your product!")
		m.App.Info.Println("Register: redirecting to user's account page")
		http.Redirect(w, r, "/v1/user/products", http.StatusSeeOther)

	}
}

func storeImagesS3(w http.ResponseWriter, r *http.Request, i, productIndex int, sess *awsS3.Session) (string, error) {

	const MAX_UPLOAD_SIZE = 1024 * 1024 // 1MB

	uploader := s3manager.NewUploader(sess)

	r.Body = http.MaxBytesReader(w, r.Body, MAX_UPLOAD_SIZE)

	if err := r.ParseMultipartForm(MAX_UPLOAD_SIZE); err != nil {
		// http.Error(w, "The uploaded file is too big. Please choose an file that's less than 1MB in size", http.StatusBadRequest)
		return "", fmt.Errorf("the uplaoded file is too big. Please choose af ile that's less than 1 MB in size :%s", err)
	}
	defer r.Body.Close()

	fileName := "file" + strconv.Itoa(i) //file1

	file, header, err := r.FormFile(fileName)
	if err != nil || header.Size == 0 {
		// http.Error(w, err.Error(), http.StatusBadRequest)
		return "", fmt.Errorf("there are no files being uploaded: %s", err)
	}
	defer file.Close()

	buff := make([]byte, 512)
	_, err = file.Read(buff)
	if err != nil {
		// http.Error(w, err.Error(), http.StatusInternalServerError)
		return "", fmt.Errorf("failed to read buff: %s", err)
	}

	filetype := http.DetectContentType(buff)
	if filetype != "image/jpeg" && filetype != "image/png" {
		// http.Error(w, "The provided file format is not allowed. Please upload a JPEG or PNG image", http.StatusBadRequest)
		// fmt.Println("error in file type !!!")
		return "", fmt.Errorf("only .jpeg and .png files are allowed: %s", err)
	}

	// Reset the file
	file.Seek(0, 0)

	s3FileName := strconv.Itoa(productIndex) + "_" + strconv.Itoa(i) + filepath.Ext(header.Filename)

	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket:               aws.String(config.AWSProductBucket),
		ACL:                  aws.String("public-read"),
		Key:                  aws.String(s3FileName),
		Body:                 file,
		ContentType:          aws.String(filetype),
		ServerSideEncryption: aws.String("AES256"),
		StorageClass:         aws.String("INTELLIGENT_TIERING"),
	})

	if err != nil {
		// fmt.Println("error with uploading file", err)
		return "", fmt.Errorf("error with uploading file: %s", err)
	}
	fmt.Println("upload to S3 bucket was successful; please check")

	return filepath.Ext(header.Filename), nil

	//return amz link:
}

func storeProfileImage(w http.ResponseWriter, r *http.Request, owner_ID int, sess *awsS3.Session) (string, error) {
	const MAX_UPLOAD_SIZE = 1024 * 1024 // 1MB
	uploader := s3manager.NewUploader(sess)
	r.Body = http.MaxBytesReader(w, r.Body, MAX_UPLOAD_SIZE)

	if err := r.ParseMultipartForm(MAX_UPLOAD_SIZE); err != nil {
		http.Error(w, "The uploaded file is too big. Please choose an file that's less than 1MB in size", http.StatusBadRequest)
		return "", err
	}

	// fileName := "file" + strconv.Itoa(owner_ID) //file1

	file, header, err := r.FormFile("profileImage")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return "", err
	}

	defer file.Close()

	buff := make([]byte, 512)
	_, err = file.Read(buff)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return "", err
	}

	filetype := http.DetectContentType(buff)
	if filetype != "image/jpeg" && filetype != "image/png" {
		http.Error(w, "The provided file format is not allowed. Please upload a JPEG or PNG image", http.StatusBadRequest)
		return "", err
	}

	// Reset the file
	file.Seek(0, 0)

	// s3FileName := "-1" + s3fileExtension

	s3FileName := "-1" + filepath.Ext(header.Filename)

	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket:               aws.String(config.AWSProfileBucketLink),
		ACL:                  aws.String("public-read"),
		Key:                  aws.String(s3FileName),
		Body:                 file,
		ContentType:          aws.String(filetype),
		ServerSideEncryption: aws.String("AES256"),
		StorageClass:         aws.String("INTELLIGENT_TIERING"),
	})

	// _, err = uploader.Upload(&s3manager.UploadInput{
	// 	Bucket: aws.String("wooteam-productslist/profile_images/"),
	// 	ACL:    aws.String("public-read"),
	// 	Key:    aws.String(s3FileName),
	// 	Body:   file,
	// })

	if err != nil {
		fmt.Println("error with uploading file", err)
	} else {
		fmt.Println("upload to S3 bucket was successful; please check")
	}
	return config.AWSProfileImageLink + s3FileName, nil

}

func (m *Repository) EditProduct(w http.ResponseWriter, r *http.Request) {

	data := make(map[string]interface{})

	x := r.URL.Query()
	fmt.Println(x)
	res := strings.ToLower(url.QueryEscape(x["edit"][0])) //hockey+sticks

	productIdEdit, err := strconv.Atoi(res)
	if err != nil {
		m.App.Error.Println("error occured while converting product index in query string to int", err)
	}

	product, err := m.DB.GetProductByID(r.Context(), productIdEdit)

	data["product"] = product

	fmt.Println("this shows product infomation", product)
	// user := m.App.Session.Get(r.Context(), "user").(model.User)
	// data["products"] = user.Products
	// data["user"] = user
	if err := render.Template(w, r, "editproduct.page.html", &render.TemplateData{
		Data: data,
		Form: &form.Form{},
	}); err != nil {
		m.App.Error.Println(err)
	}
}

func (m *Repository) EditProductPost(w http.ResponseWriter, r *http.Request) {

	data := make(map[string]interface{})

	type imageIndex struct {
		index     int
		imageType string
	}

	var (
		productPrice       float32
		s3ImageInformation []imageIndex
		productDescription string
	)

	form := form.New(r.PostForm)

	var productID = r.FormValue("productid")

	productIDInt, err := strconv.Atoi(productID)

	if err != nil {
		m.App.Error.Println("error converting string to int", err)
	}

	product, err := m.DB.GetProductByID(r.Context(), productIDInt)

	if err != nil {
		m.App.Error.Println("error retrieving information from db using productID", err)
	}

	data["product"] = product
	fmt.Println("this shows product infomation", product)
	g, ctx := errgroup.WithContext(r.Context())

	for i := 1; i < 5; i++ {
		id := i
		g.Go(func() error {
			fileName := "file" + strconv.Itoa(id) //file1/2/3/4/
			file, header, err := r.FormFile(fileName)
			if err != nil || header.Size == 0 {
				// http.Error(w, err.Error(), http.StatusBadRequest)
				fmt.Println(err)
				return err
			} else {
				defer file.Close()
				s3imgType, s3err := storeImagesS3(w, r, id, 1, m.App.AWSS3Session)
				if s3err != nil {
					m.App.Error.Println("S3 error", err)
					fmt.Println("")
					form.Errors.Add("fileupload", "Please only use .jpeg/ .png files not exceeding 1MB in size")
				} else {
					s3ImageInformation = append(s3ImageInformation, imageIndex{index: id, imageType: s3imgType})
				}
			}
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
				return nil
			}
		})
	}
	g.Go(func() error {

		form.Required("productname", "price", "brand")
		form.CheckLength("productname", 1, 255)
		productPrice = form.ProcessPrice("price")
		if len(r.FormValue("productdescription")) == 0 {
			//retain old value.
			fmt.Println("triggered1")
			productDescription = product.Description
		} else {
			fmt.Println("triggered2")
			form.CheckLength("productdescription", 1, 400)
			productDescription = r.FormValue("productdescription")

		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			return nil
		}
	})

	if err := g.Wait(); err != nil {
		m.App.Error.Println(err)
		// form.Errors.Add("fileupload", "A miniumum of 4 images are required")
	}

	// sort.Slice(imageIndex)
	fmt.Println("this is the unsorted index", s3ImageInformation)

	sort.Slice(s3ImageInformation, func(i, j int) bool {
		return s3ImageInformation[i].index < s3ImageInformation[j].index
	})

	fmt.Println("this is the sorted index", s3ImageInformation)

	s3ImageInformation = append(s3ImageInformation, imageIndex{index: id, imageType: s3imgType})

	var productImageURL []string
	for _, v := range s3ImageInformation {

		s := config.AWSProductImageLink
		productImageURL = append(productImageURL, s+strconv.Itoa(1)+"_"+strconv.Itoa(v.index)+v.imageType)

	}

	var editedProduct = model.Product{

		ID:          productIDInt,
		Brand:       r.FormValue("brand"),
		Title:       r.FormValue("productname"),
		Description: productDescription,
		Price:       productPrice,
		Images:      []string{"https://wooteam-productslist.s3.ap-southeast-1.amazonaws.com/product_list/images/50_1.jpeg"},
	}

	if len(form.Errors) != 0 {
		if err := render.Template(w, r, "editproduct.page.html", &render.TemplateData{
			Data: data,
			Form: form,
		}); err != nil {
			m.App.Error.Println(err)
		}
		return

	} else {

		g2, ctx := errgroup.WithContext(r.Context())

		g2.Go(func() error {
			err := m.DB.UpdateProducts(editedProduct)
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

		// index for elastic search

		g2.Go(func() error {
			put1, err := m.App.AWSClient.Index().
				Index("sample_product_list").
				Type("sampleproducttype").
				Id(strconv.Itoa(editedProduct.ID)).
				BodyJson(editedProduct).
				Do(r.Context())

			if err != nil {
				// Handle error
				panic(err)
			}
			fmt.Printf("Indexed tweet %s to index %s, type %s\n", put1.Id, put1.Index, put1.Type)

			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
				return nil
			}
		})

		if err := g2.Wait(); err != nil {
			m.App.Error.Println("error from g2", err)
			// form.Errors.Add("fileupload", "A miniumum of 4 images are required")
		}

		m.App.Session.Put(r.Context(), "flash", "You've successfully edited your product!")
		m.App.Info.Println("Register: redirecting to user's account page")
		productRedirectLink := fmt.Sprintf("/v1/products/%s", productID)
		http.Redirect(w, r, productRedirectLink, http.StatusSeeOther)

	}

	if err := render.Template(w, r, "editproduct.page.html", &render.TemplateData{
		Data: data,
		Form: form,
	}); err != nil {
		m.App.Error.Println(err)
	}
}
