package mysql

import (
	"context"
	"fmt"
	"goRent/internal/model"
	"time"
)

func (m *DBrepo) GetProductByID(id int) (model.Product, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	p := model.Product{}

	query := `select 
				p.id, p.owner_id, p.brand, p.category, p.title, p.rating, p.description, p.price, p.created_at, p.updated_at
			from
				products p where id = ?
			`
	if err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&p.ID,
		&p.OwnerID,
		&p.Brand,
		&p.Category,
		&p.Title,
		&p.Rating,
		&p.Description,
		&p.Price,
		&p.CreatedAt,
		&p.UpdatedAt,
	); err != nil {
		return p, fmt.Errorf("db getproductbyid: %v", err)
	}

	query = `select id, reviewer_id, reviewer_name, product_id, title, body, rating, created_at, updated_at
				from product_reviews where product_id = ?`
	rows, err := m.DB.QueryContext(ctx, query, id)
	if err != nil {
		return p, fmt.Errorf("db getproductbyid: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		r := model.ProductReview{}
		rows.Scan(
			&r.ID,
			&r.ReviewerID,
			&r.ReviewerName,
			&r.ProductID,
			&r.Title,
			&r.Body,
			&r.Rating,
			&r.CreatedAt,
			&r.UpdatedAt,
		)
		p.Reviews = append(p.Reviews, r)
	}
	if err := rows.Err(); err != nil {
		return p, fmt.Errorf("db getproductbyid: %v", err)
	}
	return p, nil
}

func (m *DBrepo) GetRentsByProductID(id int) ([]model.Rent, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rents := []model.Rent{}

	query := `select id, owner_id, renter_id, product_id, restriction_id, start_date, end_date, created_at, updated_at
		where product_id = ?`

	rows, err := m.DB.QueryContext(ctx, query, id)
	if err != nil {
		return nil, fmt.Errorf("db getrentbyproductid: %v", err)
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