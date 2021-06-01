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
	Products    []Product
	Rents       []Rent // where Rent.RenterID = User ID
	Bookings    []Rent // where Rent.OwnerID = User ID
}

type Address struct {
	PostalCode string
	StreetName string
	Block      string
	UnitNumber string
}

type Product struct {
	ID          int
	OwnerID     int // owner reference
	Brand       string
	Title       string
	Rating      float32
	Description string
	DeletedAt   time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type Restriction struct {
	ID          int
	Description string
	DeletedAt   time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type Rent struct {
	ID            int
	OwnerID       int
	RenterID      int
	ProductID     int
	RestrictionID int
	StartDate     time.Time
	EndDate       time.Time
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
