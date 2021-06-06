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

	// get rents
	query := `select 
	id, owner_id, renter_id, product_id, restriction_id, processed, start_date, end_date, created_at, updated_at
		from rents
		where renter_id = ?`
	rows, err := tx.QueryContext(ctx, query, u.ID)
	if err != nil {
		return model.User{}, fmt.Errorf("db GetUser: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		r := model.Rent{}
		if err := rows.Scan(
			&r.ID,
			&r.OwnerID,
			&r.RenterID,
			&r.ProductID,
			&r.RestrictionID,
			&r.Processed,
			&r.StartDate,
			&r.EndDate,
			&r.CreatedAt,
			&r.UpdatedAt,
		); err != nil {
			return model.User{}, fmt.Errorf("db GetUser: %v", err)
		}
		u.Rents = append(u.Rents, r)
	}
	if err := rows.Err(); err != nil {
		return model.User{}, fmt.Errorf("db GetUser: %v", err)
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

func (m *DBrepo) GetAllProducts() ([]model.Product, error) {

	var products []model.Product

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `select id, owner_id, brand, title, rating, description, price, created_at, updated_at from products`

	rows, err := m.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		p := model.Product{}
		err := rows.Scan(
			&p.ID,
			&p.OwnerID,
			&p.Brand,
			&p.Title,
			&p.Rating,
			&p.Description,
			&p.Price,
			&p.CreatedAt,
			&p.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		// RentalProductsList = append(RentalProductsList, strings.ToLower(title)+" - "+strings.ToLower(brand))
		products = append(products, p)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return products, nil
}
