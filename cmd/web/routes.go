package main

import (
	"goRent/internal/handler"
	"net/http"

	"github.com/gorilla/mux"
)

func routes() http.Handler {

	mux := mux.NewRouter()

	// default middleware
	mux.Use(SessionLoad)
	mux.Use(RecoverPanic)
	mux.Use(NoSurf)
	mux.Use(LastGetURL)

	// needs authentication
	u := mux.PathPrefix("/v1/user").Subrouter()
	u.Use(Authenticate)

	mux.HandleFunc("/", handler.Repo.Home).Methods("GET")
	mux.HandleFunc("/search", handler.Repo.Search).Methods("GET")
	mux.HandleFunc("/searchresult", handler.Repo.SearchResult).Methods("GET")

	mux.HandleFunc("/login", handler.Repo.Login).Methods("GET")
	mux.HandleFunc("/login", handler.Repo.LoginPost).Methods("POST")

	mux.HandleFunc("/register", handler.Repo.Register).Methods("GET")
	mux.HandleFunc("/register", handler.Repo.RegisterPost).Methods("POST")

	mux.HandleFunc("/v1/products/{productID}", handler.Repo.ShowProductByID).Methods("GET")

	u.HandleFunc("/logout", handler.Repo.Logout).Methods("GET")
	u.HandleFunc("/account", handler.Repo.UserAccount).Methods("GET")
	u.HandleFunc("/account/profile", handler.Repo.EditUserAccount).Methods("GET")
	u.HandleFunc("/account/profile", handler.Repo.EditUserAccountPost).Methods("POST")
	u.HandleFunc("/account/payment", handler.Repo.Payment).Methods("GET")

	u.HandleFunc("/cart", handler.Repo.GetCart).Methods("GET")
	u.HandleFunc("/products", handler.Repo.UserProducts).Methods("GET")
	u.HandleFunc("/rents", handler.Repo.UserRents).Methods("GET")
	u.HandleFunc("/bookings", handler.Repo.UserBookings).Methods("GET")

	u.HandleFunc("/createcategories", handler.Repo.CreateCategoryDataBase).Methods("GET")
	u.HandleFunc("/searchtrend", handler.Repo.SearchTrend).Methods("GET")

	u.HandleFunc("/cart", handler.Repo.GetCart).Methods("GET")
	u.HandleFunc("/cart/checkout", handler.Repo.GetCheckout).Methods("GET")
	u.HandleFunc("/cart/checkout/confirm", handler.Repo.ConfirmRents).Methods("POST")
	u.HandleFunc("/cart/checkout/confirm", handler.Repo.CheckoutConfirm).Methods("GET")

	u.HandleFunc("/addproduct", handler.Repo.AddProduct).Methods("GET")
	u.HandleFunc("/createproduct", handler.Repo.CreateProduct).Methods("POST")

	u.HandleFunc("/editproduct", handler.Repo.EditProduct).Methods("GET")
	u.HandleFunc("/editproduct", handler.Repo.EditProductPost).Methods("POST")

	a := mux.PathPrefix("/admin").Subrouter()
	// a.Use(Authenticate)
	a.Use(Authorized)

	a.HandleFunc("/overview", handler.Repo.AdminAccount).Methods("GET")
	a.HandleFunc("/overview", handler.Repo.AdminAccountPost).Methods("POST")
	a.HandleFunc("/products", handler.Repo.AdminProducts).Methods("GET")
	a.HandleFunc("/rents", handler.Repo.AdminRentals).Methods("GET")

	p := mux.PathPrefix("/v1/products").Subrouter()
	p.Use(Authenticate)
	p.HandleFunc("/addRent", handler.Repo.PostRent).Methods("POST")
	p.HandleFunc("/removeRent", handler.Repo.DeleteRent).Methods("POST")
	p.HandleFunc("/{productID}/review", handler.Repo.PostReview).Methods("POST")

	fileServer := http.FileServer(http.Dir("./static/"))
	mux.PathPrefix("/static/").Handler(http.StripPrefix("/static", fileServer))

	return mux
}
