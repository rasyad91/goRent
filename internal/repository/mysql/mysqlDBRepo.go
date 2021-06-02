package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"goRent/internal/model"
	"goRent/internal/repository"
	"log"
	"time"
)

type DBrepo struct {
	*sql.DB
}

const (
	layoutISO = "2006-01-02"
)

// NewRepo creates the repository
func NewRepo(Conn *sql.DB) repository.DatabaseRepo {
	return &DBrepo{
		DB: Conn,
	}
}

func (m *DBrepo) GetAllCourses() {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `select id, title, created_at, updated_at from courses`

	rows, err := m.QueryContext(ctx, query)
	if err != nil {
		log.Println(err)
	}
	defer rows.Close()

	for rows.Next() {

		err := rows.Scan()
		if err != nil {
			log.Println(err)

		}
	}

	if err := rows.Err(); err != nil {
		fmt.Println(err)
	}

}

func (m *DBrepo) GetUser(username string) (model.User, bool) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	results, _ := m.QueryContext(ctx, "SELECT * FROM goRent.Users where username=?", username)
	if results.Next() {
		var person model.User
		var add model.Address
		_ = results.Scan(&person.ID, &person.Username, &person.Email, &person.Password,
			&person.AccessLevel, &person.Rating, &add.PostalCode, &add.StreetName, &add.Block, &add.UnitNumber,
			&person.DeletedAt, &person.CreatedAt, &person.UpdatedAt)
		person.Address = add
		fmt.Println(person)
		fmt.Println(add)
		return person, true
	} else {
		return model.User{}, false
	}

}

func (m *DBrepo) InsertUser(u model.User) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err := m.ExecContext(ctx, "INSERT INTO goRent.users (username,email,password,postal_code,street_name,block,unit_number,created_at,updated_at) VALUES (?,?,?,?,?,?,?,?,?);",
		u.Username, u.Email, u.Password,
		u.Address.PostalCode, u.Address.StreetName, u.Address.Block, u.Address.UnitNumber,
		(u.CreatedAt).Format(layoutISO), (u.UpdatedAt).Format(layoutISO))
	fmt.Println("NORMAL DATE", u.DeletedAt)
	fmt.Println("ISO DATE", (u.DeletedAt).Format(layoutISO))

	fmt.Println("INSERTION ERR:", err)
	if err != nil {
		return false
	} else {
		return true
	}

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
