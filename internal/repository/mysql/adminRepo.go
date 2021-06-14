package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"goRent/internal/model"
	"time"
)

func (m *DBrepo) GetAllUsers() ([]model.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	result := []model.User{}
	rows, _ := m.DB.QueryContext(ctx, "SELECT id,username,image_url,access_level FROM gorent.users")
	defer rows.Close()
	for rows.Next() {
		u := model.User{}
		if err := rows.Scan(
			&u.ID,
			&u.Username,
			&u.Image_URL,
			&u.AccessLevel,
		); err != nil {
			fmt.Println("ERROR", err)
			if err == sql.ErrNoRows {
				return []model.User{}, err
			}
		}
		result = append(result, u)
	}
	fmt.Println("INSIDE GETALLUSERS")
	fmt.Println(result)
	return result, nil
}
func (m *DBrepo) GrantAccess(u string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err := m.ExecContext(ctx, "UPDATE users SET access_level = 1 where id = ?", u)
	if err != nil {
		return fmt.Errorf("db InsertUser: %v", err)
	}
	return nil
}

func (m *DBrepo) RemoveAccess(u string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err := m.ExecContext(ctx, "UPDATE users SET access_level = 5 where id = ?", u)
	if err != nil {
		return fmt.Errorf("db InsertUser: %v", err)
	}
	return nil
}
