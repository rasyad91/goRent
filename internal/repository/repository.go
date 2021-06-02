package repository

import "goRent/internal/model"

type DatabaseRepo interface {
	GetAllCourses()
<<<<<<< HEAD
	GetUser(username string) (model.User, bool)
=======
	GetAllProducts() ([]model.Product, error)
>>>>>>> origin/alvinProducts
}
