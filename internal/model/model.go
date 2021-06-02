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
	Products    []Product // where ID = Product.OwnerID
	Rents       []Rent    // where ID = Rent.RenterID
	Bookings    []Rent    // where ID = Rent.OwnerID
}

type Address struct {
	PostalCode string
	StreetName string
	Block      string
	UnitNumber string
}

type Product struct {
	ID          int
	OwnerID     int
	Brand       string
	Title       string
	Rating      float32
	Description string
	Price       float32
	Reviews     []Review // where ID = Review.ProductID
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

type Review struct {
	ID        int
	OwnerID   int
	RenterID  int
	ProductID int
	Type      string
	Title     string
	Body      string
	Rating    float32
	CreatedAt time.Time
	UpdatedAt time.Time
}
