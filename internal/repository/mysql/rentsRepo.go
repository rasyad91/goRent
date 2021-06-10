package mysql

import (
	"context"
	"fmt"
	"goRent/internal/model"
	"time"
)

func (m *DBrepo) CreateRent(r model.Rent) (int, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `INSERT INTO rents (owner_id, renter_id, product_id, restriction_id, processed,duration, total_cost, start_date, end_date, created_at, updated_at)
				VALUES (?,?,?,?,?,?,?,?,?,?,?)`

	result, err := m.DB.ExecContext(ctx, query,
		r.OwnerID,
		r.RenterID,
		r.ProductID,
		1,
		false,
		r.Duration,
		r.TotalCost,
		r.StartDate,
		r.EndDate,
		time.Now(),
		time.Now(),
	)
	if err != nil {
		return 0, fmt.Errorf("db createrent: %v", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (m *DBrepo) GetRentsByProductID(id int) ([]model.Rent, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rents := []model.Rent{}

	query := `select 
	id, owner_id, renter_id, product_id, restriction_id, start_date, end_date, created_at, updated_at
		from rents
		where product_id = ? and processed = true`

	rows, err := m.DB.QueryContext(ctx, query, id)
	if err != nil {
		return nil, fmt.Errorf("db line76 getrentbyproductid: %v", err)
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
			&r.StartDate,
			&r.EndDate,
			&r.CreatedAt,
			&r.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("db getrentbyproductid: %v", err)
		}
		rents = append(rents, r)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("db getrentbyproductid: %v", err)
	}
	return rents, nil
}

func (m *DBrepo) DeleteRent(rentID int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `DELETE from rents where id = ?`

	_, err := m.DB.ExecContext(ctx, query, rentID)
	if err != nil {
		return fmt.Errorf("db deleterent: %v", err)
	}

	return nil
}
