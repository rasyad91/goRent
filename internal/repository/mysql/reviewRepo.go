package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"goRent/internal/model"
	"time"
)

func (m *dbRepo) CreateProductReview(pr model.ProductReview) (float32, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// create new transaction
	tx, err := m.DB.BeginTx(ctx, nil)
	if err != nil {
		return 0, fmt.Errorf("db addproductreview: %v", err)
	}

	defer func() {
		if err := recover(); err != nil {
			tx.Rollback()
		}
	}()

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
		if err == sql.ErrNoRows {
			reviewCount = 0
			rating = pr.Rating
		} else {
			tx.Rollback()
			return 0, fmt.Errorf("db addproductreview query reviewcount + rating: %v", err)
		}
		fmt.Println("db error:", err)
	}

	fmt.Println("rating:", rating)
	fmt.Println("PRrating:", pr.Rating)

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
		return 0, fmt.Errorf("db addproductreview insert review: %v", err)
	}

	// update rating on product
	//check
	var newRating float32
	if reviewCount == 0 {
		newRating = rating
	} else {
		newRating = rating + (pr.Rating-rating)/float32(reviewCount+1)
	}

	query = `UPDATE products SET rating = ? WHERE (id = ?);`
	if _, err := tx.ExecContext(ctx, query, newRating, pr.ProductID); err != nil {
		tx.Rollback()
		return 0, fmt.Errorf("db addproductreview update rating: %v", err)
	}
	tx.Commit()

	return newRating, nil
}
