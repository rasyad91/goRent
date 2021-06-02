package repository

import "goRent/internal/model"

type DatabaseRepo interface {
	GetAllCourses()
	GetAllProducts() ([]model.Product, error)
}
