package repository

import "goRent/internal/model"

type DatabaseRepo interface {
	GetUser(username string) (model.User, error)
	InsertUser(user model.User) error
	GetAllProducts() ([]model.Product, error)
}
