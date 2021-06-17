package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"goRent/internal/model"
	"time"
)

func (m *dbRepo) GetAllUsers() ([]model.User, error) {
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
	return result, nil
}
func (m *dbRepo) GrantAccess(u string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err := m.ExecContext(ctx, "UPDATE users SET access_level = 1 where id = ?", u)
	if err != nil {
		return fmt.Errorf("db GrantAccess: %v", err)
	}
	return nil
}

func (m *dbRepo) RemoveAccess(u string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err := m.ExecContext(ctx, "UPDATE users SET access_level = 5 where id = ?", u)
	if err != nil {
		return fmt.Errorf("db RemoveAccess: %v", err)
	}
	return nil
}

func (m *dbRepo) GetAllRents() ([]model.Rent, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	result := []model.Rent{}
	rows, _ := m.DB.QueryContext(ctx, "SELECT id, owner_id, renter_id, product_id, processed, duration,total_cost, start_date, end_date FROM gorent.rents order by product_id ")
	defer rows.Close()
	for rows.Next() {
		r := model.Rent{}
		if err := rows.Scan(
			&r.ID,
			&r.OwnerID,
			&r.RenterID,
			&r.ProductID,
			&r.Processed,
			&r.Duration,
			&r.TotalCost,
			&r.StartDate,
			&r.EndDate,
		); err != nil {
			fmt.Println("ERROR", err)
			if err == sql.ErrNoRows {
				return []model.Rent{}, err
			}
		}
		result = append(result, r)
	}
	return result, nil
}

func (m *dbRepo) DeleteUser(u string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err := m.ExecContext(ctx, "DELETE FROM users where id = ?", u)
	if err != nil {
		return fmt.Errorf("db DeleteUsers: %v", err)
	}
	return nil
}
