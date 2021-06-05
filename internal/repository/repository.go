package repository

import "goRent/internal/model"

type DatabaseRepo interface {
	// Users
	GetUser(username string) (model.User, error)
	InsertUser(user model.User) error

	// Products
	GetAllProducts() ([]model.Product, error)
	GetProductByID(id int) (model.Product, error)

	// Rents
	GetRentsByProductID(id int) ([]model.Rent, error)
}
