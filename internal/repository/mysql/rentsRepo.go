package mysql

import (
	"context"
	"errors"
	"fmt"
	"goRent/internal/model"
	"sync"
	"time"
)

var rentLock sync.Mutex

const (
	ErrRentNotAvailable = "rent not available"
)

func (m *dbRepo) CreateRent(r model.Rent) (int, error) {
	rentLock.Lock()
	defer rentLock.Unlock()

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

func (m *dbRepo) GetRentsByProductID(ctx context.Context, id int) ([]model.Rent, error) {

	rentLock.Lock()
	defer rentLock.Unlock()

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
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

func (m *dbRepo) DeleteRent(rentID int) error {

	rentLock.Lock()
	defer rentLock.Unlock()

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `DELETE from rents where id = ?`

	_, err := m.DB.ExecContext(ctx, query, rentID)
	if err != nil {
		return fmt.Errorf("db deleterent: %v", err)
	}

	return nil
}

func (m *dbRepo) ProcessRent(rent model.Rent) error {

	rentLock.Lock()
	defer rentLock.Unlock()

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	tx, err := m.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer func() {
		if err := recover(); err != nil {
			tx.Rollback()
		}
	}()

	available, err := m.IsRentAvailable(ctx, rent.ProductID, rent.StartDate, rent.EndDate)
	if err != nil {
		return err
	}
	fmt.Println("available result: ", available)
	if !available {
		return errors.New(ErrRentNotAvailable)
	}

	query := `UPDATE rents set processed = true, updated_at = ? where id = ?`
	if _, err := tx.ExecContext(ctx, query, time.Now(), rent.ID); err != nil {
		tx.Rollback()
		return fmt.Errorf("db processrent: %v", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("db processrent: %v", err)
	}

	return nil
}

// IsRentAvailable returns true if there is no clashing date, and product is available for rent
func (m *dbRepo) IsRentAvailable(ctx context.Context, productId int, startDate, endDate time.Time) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	tx, err := m.BeginTx(ctx, nil)
	if err != nil {
		fmt.Println(err)
		return false, err
	}

	defer func() {
		if err := recover(); err != nil {
			tx.Rollback()
		}
	}()

	if _, err := tx.ExecContext(ctx, `SET @sd := ?;`, startDate); err != nil {
		fmt.Println(err)
		return false, err
	}
	if _, err := tx.ExecContext(ctx, `SET @ed := ?;`, endDate); err != nil {
		fmt.Println(err)
		return false, err
	}

	query := `
		SELECT count(*)
		FROM rents 
		WHERE 	processed = true AND
				product_id = ? AND
				(
					(start_date BETWEEN @sd AND @ed) OR
					(end_date BETWEEN @sd AND @ed) OR
					(start_date <= @sd AND  end_date >= @ed)
				);
			`

	var count int
	if err := tx.QueryRowContext(ctx, query, productId).Scan(&count); err != nil {
		fmt.Println(err)
		return false, err
	}

	if count != 0 {
		return false, nil
	}

	return true, nil
}
