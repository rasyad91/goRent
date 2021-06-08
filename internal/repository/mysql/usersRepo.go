package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"goRent/internal/model"
	"goRent/internal/repository"
	"time"
)

type DBrepo struct {
	*sql.DB
}

// const (
// 	layoutISO = "2006-01-02"
// )

// NewRepo creates the repository
func NewRepo(Conn *sql.DB) repository.DatabaseRepo {
	return &DBrepo{
		DB: Conn,
	}
}

func (m *DBrepo) GetUser(username string) (model.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	u := model.User{}
	tx, err := m.DB.BeginTx(ctx, nil)
	if err != nil {
		return model.User{}, fmt.Errorf("db GetUser: %v", err)
	}
	if err := tx.QueryRowContext(ctx, "SELECT * FROM goRent.Users where username=?", username).
		Scan(
			&u.ID,
			&u.Username,
			&u.Email,
			&u.Password,
			&u.AccessLevel,
			&u.Rating,
			&u.Address.PostalCode,
			&u.Address.StreetName,
			&u.Address.Block,
			&u.Address.UnitNumber,
			&u.CreatedAt,
			&u.UpdatedAt,
		); err != nil {
		return model.User{}, fmt.Errorf("db GetUser: %v", err)
	}
	// add in concurrency here - 1 get Rent 1 get Product
	// get rents
	rent_query := `select 
			r.id, r.owner_id, r.renter_id, r.product_id, r.restriction_id, r.processed, r.start_date, r.end_date, r.duration, r.total_cost, r.created_at, r.updated_at,
			p.id, p.owner_id, p.brand, p.category, p.title, p.rating, p.description, p.price, p.created_at, p.updated_at
		from 
			rents r 
		left join 
			products p on (p.id = r.product_id)
		where 
			r.renter_id = ?`
	product_query := `select * from products p where p.owner_id = ?`
	booking_query := `select 
		r.id, r.owner_id, r.renter_id, r.product_id, r.restriction_id, r.processed, r.start_date, r.end_date, r.duration, r.total_cost, r.created_at, r.updated_at,
		p.id, p.owner_id, p.brand, p.category, p.title, p.rating, p.description, p.price, p.created_at, p.updated_at
	from 
		rents r 
	left join 
		products p on (p.id = r.product_id)
	where 
		r.owner_id = ?`
	//initializing concurrency

	// rent_query
	rent_rows, err := tx.QueryContext(ctx, rent_query, u.ID)
	if err != nil {
		return model.User{}, fmt.Errorf("db GetUser: %v", err)
	}
	defer rent_rows.Close()
	for rent_rows.Next() {
		r := model.Rent{}
		if err := rent_rows.Scan(
			&r.ID,
			&r.OwnerID,
			&r.RenterID,
			&r.ProductID,
			&r.RestrictionID,
			&r.Processed,
			&r.StartDate,
			&r.EndDate,
			&r.Duration,
			&r.TotalCost,
			&r.CreatedAt,
			&r.UpdatedAt,
			&r.Product.ID,
			&r.Product.OwnerID,
			&r.Product.Brand,
			&r.Product.Category,
			&r.Product.Title,
			&r.Product.Rating,
			&r.Product.Description,
			&r.Product.Price,
			&r.Product.CreatedAt,
			&r.Product.UpdatedAt,
		); err != nil {
			return model.User{}, fmt.Errorf("db GetUser: %v", err)
		}
		u.Rents = append(u.Rents, r)
	}
	//booking_query
	booking_rows, err := tx.QueryContext(ctx, booking_query, u.ID)
	if err != nil {
		return model.User{}, fmt.Errorf("db GetUser: %v", err)
	}
	defer booking_rows.Close()
	for booking_rows.Next() {
		b := model.Rent{}
		if err := booking_rows.Scan(
			&b.ID,
			&b.OwnerID,
			&b.RenterID,
			&b.ProductID,
			&b.RestrictionID,
			&b.Processed,
			&b.StartDate,
			&b.EndDate,
			&b.Duration,
			&b.TotalCost,
			&b.CreatedAt,
			&b.UpdatedAt,
			&b.Product.ID,
			&b.Product.OwnerID,
			&b.Product.Brand,
			&b.Product.Category,
			&b.Product.Title,
			&b.Product.Rating,
			&b.Product.Description,
			&b.Product.Price,
			&b.Product.CreatedAt,
			&b.Product.UpdatedAt,
		); err != nil {
			return model.User{}, fmt.Errorf("db GetUser: %v", err)
		}
		u.Bookings = append(u.Bookings, b)
	}

	//product query
	product_rows, err := tx.QueryContext(ctx, product_query, u.ID)
	if err != nil {
		return model.User{}, fmt.Errorf("db GetUser: %v", err)
	}
	defer product_rows.Close()
	for product_rows.Next() {
		p := model.Product{}
		if err := product_rows.Scan(
			&p.ID,
			&p.OwnerID,
			&p.Brand,
			&p.Category,
			&p.Title,
			&p.Rating,
			&p.Description,
			&p.Price,
			// &p.Images,
			&p.CreatedAt,
			&p.UpdatedAt,
		); err != nil {
			return model.User{}, fmt.Errorf("db GetUser: %v", err)
		}
		u.Products = append(u.Products, p)
	}

	if err := booking_rows.Err(); err != nil {
		return model.User{}, fmt.Errorf("db GetUser: %v", err)
	}
	fmt.Println("PRODUCTS QUERY:")
	for _, item := range u.Products {
		fmt.Println(item.Title)
	}
	fmt.Println("Rents QUERY:")
	for _, item := range u.Rents {
		fmt.Println(item.Product.Title)
	}

	fmt.Println("Bookings QUERY:")
	for _, item := range u.Bookings {
		fmt.Println(item.Product.Title)
	}
	return u, nil
}

func (m *DBrepo) InsertUser(u model.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err := m.ExecContext(ctx, "INSERT INTO goRent.users (username,email,password,postal_code,street_name,block,unit_number,created_at,updated_at) VALUES (?,?,?,?,?,?,?,?,?);",
		u.Username, u.Email, u.Password,
		u.Address.PostalCode, u.Address.StreetName, u.Address.Block, u.Address.UnitNumber,
		time.Now(), time.Now())
	if err != nil {
		return fmt.Errorf("db InsertUser: %v", err)
	}
	return nil
}

func (m *DBrepo) EditUser(u model.User, editType string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := error(nil)
	if editType == "address" {
		_, err = m.ExecContext(ctx, "UPDATE goRent.users SET block = ?, street_name = ?, unit_number = ?, postal_code = ? WHERE id = ?", u.Address.Block, u.Address.StreetName, u.Address.UnitNumber, u.Address.PostalCode, u.ID)
		fmt.Println("AddressChange test:", u)
	} else if editType == "profile" {
		_, err = m.ExecContext(ctx, "UPDATE goRent.users SET username = ?, email = ? WHERE id = ?", u.Username, u.Email, u.ID)
		fmt.Println("ProfileChange test:", u)
	} else {
		_, err = m.ExecContext(ctx, "UPDATE goRent.users SET password = ? WHERE id = ?", u.Password, u.ID)
		fmt.Println("PassWordChange test:", u)
	}
	if err != nil {
		return fmt.Errorf("db EditUser: %v", err)
	}
	return nil
}
