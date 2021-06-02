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
	results, _ := m.Query("SELECT * FROM goRent.Users where username=?", username)
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
