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
			r.id, r.owner_id, r.renter_id, r.product_id, r.restriction_id, r.processed, r.start_date, r.end_date, r.duration, r.total_cost, r.created_at, r.updated_at,
			p.id, p.owner_id, p.brand, p.category, p.title, p.rating, p.description, p.price, p.created_at, p.updated_at
		from 
			rents r 
		left join 
			products p on (p.id = r.product_id)
		where 
			r.renter_id = ?`

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

func (m *DBrepo) EditUser(u model.User, editType string) error {
	// ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	// defer cancel()
	err := error(nil)
	if editType == "address" {
		// _, err = m.ExecContext(ctx, "UPDATE goRent.users SET block = ?, street_name = ?, unit_number = ?, postal_code = ? WHERE id = ?", u.Address.Block, u.Address.StreetName, u.Address.UnitNumber, u.Address.PostalCode, u.ID)
		fmt.Println("AddressChange test:", u)
	} else if editType == "profile" {
		// _, err = m.ExecContext(ctx, "UPDATE goRent.users SET username = ?, email = ? WHERE id = ?", u.Username, u.Email, u.ID)
		fmt.Println("ProfileChange test:", u)
	} else {
		// _, err = m.ExecContext(ctx, "UPDATE goRent.users SET password = ? WHERE id = ?", u.Password, u.ID)
		fmt.Println("PassWordChange test:", u)
	}
	if err != nil {
		return fmt.Errorf("db InsertUser: %v", err)
	}
	return nil
}
