package repository

import (
	"context"
	"goRent/internal/model"
)

type DatabaseRepo interface {
	//admin
	GetAllUsers() ([]model.User, error)
	GrantAccess(userid string) error
	RemoveAccess(userid string) error
	// Users
	GetUser(username string) (model.User, error)
	InsertUser(user model.User) error
	EditUser(user model.User, editType string) error
	EmailExist(email string) error

	// Products
	GetAllProducts() ([]model.Product, error)
	GetProductByID(ctx context.Context, id int) (model.Product, error)
	GetProductNextIndex() (int, error)
	InsertProduct(model.Product) error
	InsertProductImages(i int, s string) error
	UpdateProducts(p model.Product, s1 []model.ImgUrl, s2 []string) error

	// Rents
	GetRentsByProductID(ctx context.Context, id int) ([]model.Rent, error)
	CreateRent(r model.Rent) (int, error)
	DeleteRent(rentID int) error
	ProcessRent(rent model.Rent) error

	// Reviews
	CreateProductReview(pr model.ProductReview) error
}
