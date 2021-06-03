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

<<<<<<< HEAD
	mux.HandleFunc("/", handler.Repo.Home).Methods("GET")
	mux.HandleFunc("/search", handler.Repo.Search).Methods("GET")
	mux.HandleFunc("/searchresult", handler.Repo.SearchResult).Methods("GET")
=======
<<<<<<< HEAD
=======
	router.HandleFunc("/user/logout", handler.Repo.Logout).Methods("GET")

>>>>>>> Login
	router.HandleFunc("/login", handler.Repo.Login).Methods("GET")
	router.HandleFunc("/login", handler.Repo.LoginPost).Methods("POST")
>>>>>>> d0bb68ff913320fbff4313a78419e5157181b05e

	mux.HandleFunc("/register", handler.Repo.Register).Methods("GET")
	mux.HandleFunc("/register", handler.Repo.RegisterPost).Methods("POST")

	mux.HandleFunc("/user/{userID}/bookings", handler.Repo.UserBookings).Methods("GET")
	mux.HandleFunc("/user/{userID}/rents", handler.Repo.UserRents).Methods("GET")
	mux.HandleFunc("/user/{userID}/products", handler.Repo.UserProducts).Methods("GET")

	mux.HandleFunc("/user/logout", handler.Repo.Logout).Methods("GET")

	// sub := mux.NewRoute().Subrouter()
	// sub.Use(handler.ValidationAPIMiddleware)

	mux.HandleFunc("/v1/products/{productId}", handler.Repo.ShowProductByID).Methods("GET")

	// fileServer := http.FileServer(http.Dir("./static/"))
	// mux.Handle("/static/", http.StripPrefix("/static", fileServer))
		// static files
		fileServer := http.FileServer(http.Dir("./static/"))
		mux.Handle("/static/", http.StripPrefix("/static/", fileServer))
		// mux.PathPrefix("/").Handler(http.FileServer(http.Dir("./static/")))

	return mux
}
