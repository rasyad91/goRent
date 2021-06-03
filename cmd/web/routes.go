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

	router.HandleFunc("/", handler.Repo.Home).Methods("GET")

	router.HandleFunc("/login", handler.Repo.Login).Methods("GET")
	router.HandleFunc("/login", handler.Repo.LoginPost).Methods("POST")

	router.HandleFunc("/register", handler.Repo.Register).Methods("GET")
	router.HandleFunc("/register", handler.Repo.RegisterPost).Methods("POST")

	// sub := router.NewRoute().Subrouter()
	// sub.Use(handler.ValidationAPIMiddleware)

	fileServer := http.FileServer(http.Dir("../../static/"))
	router.Handle("/static/*", http.StripPrefix("/static", fileServer))

	return router
}
