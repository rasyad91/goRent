package repository

import (
	"context"
	"goRent/internal/model"
	"time"
)

type DatabaseRepo interface {
	// Users
	GetUser(username string) (model.User, error)
	InsertUser(user model.User) error
	EditUser(user model.User, editType string) error

	// Products
	GetAllProducts() ([]model.Product, error)
	GetProductByID(ctx context.Context, id int) (model.Product, error)
	GetProductNextIndex() (int, error)

	// Rents
	GetRentsByProductID(ctx context.Context, id int) ([]model.Rent, error)
	CreateRent(r model.Rent) (int, error)
	DeleteRent(rentID int) error
	ProcessRent(rent model.Rent) error
	IsRentAvailable(ctx context.Context, productId int, startDate, endDate time.Time) (bool, error)

	// Reviews
	CreateProductReview(pr model.ProductReview) error
}
