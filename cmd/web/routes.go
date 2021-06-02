package main

import (
	"goRent/internal/handler"
	"net/http"

	"github.com/gorilla/mux"
)

func routes() http.Handler {

	router := mux.NewRouter()
	router.HandleFunc("/api/v1/", handler.Repo.Home).Methods("GET")
	router.HandleFunc("/user/logout", handler.Repo.Logout).Methods("GET")
	router.HandleFunc("/login", handler.Repo.Login)
	router.HandleFunc("/register", handler.Repo.Register)

	sub := router.NewRoute().Subrouter()
	sub.Use(handler.ValidationAPIMiddleware)

	sub.HandleFunc("/api/v1/courses/{courseId}", handler.Repo.GetCourse).Methods("GET")
	sub.HandleFunc("/api/v1/courses/{courseId}", handler.Repo.PostCourse).Methods("POST")
	sub.HandleFunc("/api/v1/courses/{courseId}", handler.Repo.PutCourse).Methods("PUT")
	sub.HandleFunc("/api/v1/courses/{courseId}", handler.Repo.DeleteCourse).Methods("DELETE")

	fileServer := http.FileServer(http.Dir("./static/"))
	router.Handle("/static/*", http.StripPrefix("/static", fileServer))

	return router
}
