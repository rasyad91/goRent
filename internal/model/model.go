package model

import "time"

type User struct {
	ID          int
	Username    string
	Email       string
	Password    string
	AccessLevel int
	Rating      float32
	Address     Address
	DeletedAt   time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type Address struct {
	PostalCode string
	StreetName string
	Block      string
	UnitNumber string
}

type Product struct {
	ID          int
	UserID      int // owner reference
	Brand       string
	Title       string
	Rating      float32
	Description string
	DeletedAt   time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
