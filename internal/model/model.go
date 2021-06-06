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
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Products    []Product    // where ID = Product.OwnerID from products table
	Rents       []Rent       // where ID = Rent.RenterID from rents table
	Bookings    []Rent       // where ID = Rent.OwnerID from rents table
	UserReviews []UserReview // where ID = UserReviews.ReceiverID from reviewstable
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
	Category    string
	Title       string
	Rating      float32
	Description string
	Price       float32
	Reviews     []ProductReview // where ID = ProductReview.ProductID from reviews table
	Images      []string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type Restriction struct {
	ID          int // ID = 1, booked by user, ID = 2, blocked by owner
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
	Processed     bool // false = in cart, true = checkedout
	TotalCost     float32
	StartDate     time.Time
	EndDate       time.Time
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type UserReview struct {
	ID           int
	ReviewerID   int // the one thats making the review
	ReviewerName string
	ReceiverID   int // the one thats get reviewed
	Body         string
	Rating       float32
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type ProductReview struct {
	ID           int
	ReviewerID   int // the one thats making the review
	ReviewerName string
	ProductID    int
	Body         string
	Rating       float32
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
