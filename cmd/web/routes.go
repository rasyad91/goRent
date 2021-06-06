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

	mux.HandleFunc("/", handler.Repo.Home).Methods("GET")
	mux.HandleFunc("/search", handler.Repo.Search).Methods("GET")
	mux.HandleFunc("/searchresult", handler.Repo.SearchResult).Methods("GET")

	mux.HandleFunc("/login", handler.Repo.Login).Methods("GET")
	mux.HandleFunc("/login", handler.Repo.LoginPost).Methods("POST")

	mux.HandleFunc("/register", handler.Repo.Register).Methods("GET")
	mux.HandleFunc("/register", handler.Repo.RegisterPost).Methods("POST")

	mux.PathPrefix("/user").Subrouter().Use(Auth)

	mux.HandleFunc("/user/logout", handler.Repo.Logout).Methods("GET")
	mux.HandleFunc("/v1/user/account", handler.Repo.UserAccount).Methods("GET")

	mux.PathPrefix("/auth").Subrouter().Use(Auth)

	mux.HandleFunc("/v1/products/{productID}", handler.Repo.ShowProductByID).Methods("GET")
	mux.HandleFunc("/v1/products/{productID}/review", handler.Repo.PostReview).Methods("POST")

	fileServer := http.FileServer(http.Dir("./static/"))
	mux.PathPrefix("/static/").Handler(http.StripPrefix("/static", fileServer))

	// static files

	// fileServer := http.FileServer(http.Dir("/static/"))
	// mux.Handle("/static/", http.StripPrefix("/static", fileServer))
	// mux.PathPrefix("/").Handler(http.FileServer(http.Dir("./static/")))

	return mux
}
