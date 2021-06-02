package repository

import "goRent/internal/model"

type DatabaseRepo interface {
	GetAllCourses()
	GetUser(username string) (model.User, bool)
	GetAllProducts() ([]model.Product, error)
}
