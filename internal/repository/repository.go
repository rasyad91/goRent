package repository

import (
	"context"
	"goRent/internal/model"
)

type Database interface {
	AdminRepo
	UserRepo
	ProductRepo
	RentRepo
	ReviewRepo
}

type AdminRepo interface {
	GetAllUsers() ([]model.User, error)
	GrantAccess(userid string) error
	RemoveAccess(userid string) error
	GetAllRents() ([]model.Rent, error)
	DeleteUser(userid string) error
}

type UserRepo interface {
	GetUser(username string) (model.User, error)
	InsertUser(user model.User) error
	EditUser(user model.User, editType string) error
	EmailExist(email string) error
}

type ProductRepo interface {
	GetAllProducts() ([]model.Product, error)
	GetProductByID(ctx context.Context, id int) (model.Product, error)
	GetProductNextIndex() (int, error)
	InsertProduct(model.Product) error
	InsertProductImages(i int, s string) error
	UpdateProducts(p model.Product, s1 []model.ImgUrl, s2 []string) error
}

type RentRepo interface {
	GetRentsByProductID(ctx context.Context, id int) ([]model.Rent, error)
	CreateRent(r model.Rent) (int, error)
	DeleteRent(rentID int) error
	ProcessRent(rent model.Rent) error
}

type ReviewRepo interface {
	CreateProductReview(pr model.ProductReview) (float32, error)
}
