package mysql

import (
	"database/sql"
	"goRent/internal/repository"
)

// dbRepo struct that satisfy the repository.Database interface
type dbRepo struct {
	*sql.DB
}

// NewRepo creates the repository in mysql
func NewRepo(conn *sql.DB) repository.Database {
	return &dbRepo{
		DB: conn,
	}
}
