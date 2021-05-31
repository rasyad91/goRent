package mysql

import (
	"context"
	"database/sql"
	"fmt"
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
