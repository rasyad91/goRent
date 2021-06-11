package repository

import (
	"context"
	"goRent/internal/model"
)

type DatabaseRepo interface {
	// Users
	GetUser(username string) (model.User, error)
	InsertUser(user model.User) error
	EditUser(user model.User, editType string) error
	EmailExist(email string) error

	// Products
	GetAllProducts() ([]model.Product, error)
	GetProductByID(ctx context.Context, id int) (model.Product, error)
	GetProductNextIndex() (int, error)

	// Rents
	GetRentsByProductID(ctx context.Context, id int) ([]model.Rent, error)
	CreateRent(r model.Rent) (int, error)
	DeleteRent(rentID int) error

	// Reviews
	CreateProductReview(pr model.ProductReview) error
}
