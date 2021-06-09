package repository

import "goRent/internal/model"

type DatabaseRepo interface {
	// Users
	GetUser(username string) (model.User, error)
	InsertUser(user model.User) error
	EditUser(user model.User, editType string) error

	// Products
	GetAllProducts() ([]model.Product, error)
	GetProductByID(id int) (model.Product, error)

	// Rents
	GetRentsByProductID(id int) ([]model.Rent, error)
	CreateRent(r model.Rent) error
	DeleteRent(rentID int) error

	// Reviews
	CreateProductReview(pr model.ProductReview) error
}
