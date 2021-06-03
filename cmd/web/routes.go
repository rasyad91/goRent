package main

import (
	"goRent/internal/handler"
	"net/http"

	"github.com/gorilla/mux"
)

func routes() http.Handler {

	router := mux.NewRouter()

	// default middleware
	router.Use(SessionLoad)
	router.Use(RecoverPanic)
	router.Use(NoSurf)

	router.HandleFunc("/", handler.Repo.Home).Methods("GET")
	router.HandleFunc("/user/logout", handler.Repo.Logout).Methods("GET")

	router.HandleFunc("/search", handler.Repo.Search).Methods("GET")
	router.HandleFunc("/searchresult", handler.Repo.SearchResult).Methods("GET")

	router.HandleFunc("/login", handler.Repo.Login).Methods("GET")
	router.HandleFunc("/login", handler.Repo.LoginPost).Methods("POST")

	router.HandleFunc("/register", handler.Repo.Register).Methods("GET")
	router.HandleFunc("/register", handler.Repo.RegisterPost).Methods("POST")

	router.HandleFunc("/user/{userID}/bookings", handler.Repo.UserBookings).Methods("GET")
	router.HandleFunc("/user/{userID}/rents", handler.Repo.UserRents).Methods("GET")
	router.HandleFunc("/user/{userID}/products", handler.Repo.UserProducts).Methods("GET")

	router.HandleFunc("/user/logout", handler.Repo.Logout).Methods("GET")

	// sub := router.NewRoute().Subrouter()
	// sub.Use(handler.ValidationAPIMiddleware)

	router.HandleFunc("/v1/products/{productId}", handler.Repo.ShowProductByID).Methods("GET")

	fileServer := http.FileServer(http.Dir("../../static/"))
	router.Handle("/static/*", http.StripPrefix("/static", fileServer))

	return router
}
