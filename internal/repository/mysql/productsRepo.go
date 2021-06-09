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

	query = `select username from users where id = ?`
	if err := m.DB.QueryRowContext(ctx, query, p.OwnerID).Scan(
		&p.OwnerName,
	); err != nil {
		return p, fmt.Errorf("db getproductbyid: %v", err)
	}

	query = `select id, reviewer_id, reviewer_name, product_id, body, rating, created_at, updated_at
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

func (m *DBrepo) CreateProductReview(pr model.ProductReview) error {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// create new transaction
	tx, err := m.DB.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("db addproductreview: %v", err)
	}

	// get count of reviews of particular product
	var reviewCount int
	var rating float32

	query := `select count(pr.id), p.rating 
				from product_reviews pr
				left join 
				products p on (p.id = pr.product_id)
				where pr.product_id = ?
				group by p.rating`
	if err := tx.QueryRowContext(ctx, query, pr.ProductID).Scan(&reviewCount, &rating); err != nil {
		return fmt.Errorf("db addproductreview query reviewcount + rating: %v", err)
	}

	// insert new product review
	query = `insert into product_reviews(reviewer_id, reviewer_name, product_id, body, rating, created_at, updated_at)
			values(?,?,?,?,?,?,?)`
	if _, err := tx.ExecContext(ctx, query,
		pr.ReviewerID,
		pr.ReviewerName,
		pr.ProductID,
		pr.Body,
		pr.Rating,
		time.Now(),
		time.Now(),
	); err != nil {
		tx.Rollback()
		return fmt.Errorf("db addproductreview insert review: %v", err)
	}

	// update rating on product

	newRating := rating + (pr.Rating-rating)/float32(reviewCount)

	query = `UPDATE products SET rating = ? WHERE (id = ?);`
	if _, err := tx.ExecContext(ctx, query, newRating, pr.ProductID); err != nil {
		tx.Rollback()
		return fmt.Errorf("db addproductreview update rating: %v", err)
	}
	tx.Commit()

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

func (m *DBrepo) GetProductNextIndex() (int, error) {

	p := model.Product{}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `select id from products order by id desc limit 1`

	rows, err := m.QueryContext(ctx, query)
	if err != nil {
		return -1, err
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(
			&p.ID,
		)
		if err != nil {
			return -1, err
		}
	}
	if err := rows.Err(); err != nil {
		return -1, err
	}
	return p.ID + 1, nil
}
