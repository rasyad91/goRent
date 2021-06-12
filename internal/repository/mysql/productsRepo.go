package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"goRent/internal/model"
	"time"

	"golang.org/x/sync/errgroup"
)

func (m *DBrepo) GetProductByID(ctx context.Context, id int) (model.Product, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	p := model.Product{}
	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
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
			if err == sql.ErrNoRows {
				return err
			}
			return fmt.Errorf("db getproductbyid: %v", err)
		}
		query = `select username from users where id = ?`
		if err := m.DB.QueryRowContext(ctx, query, p.OwnerID).Scan(&p.OwnerName); err != nil {
			return err
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			return nil
		}
	})

	// get reviews from reviews table
	g.Go(func() error {
		query := `select id, reviewer_id, reviewer_name, product_id, body, rating, created_at, updated_at
		from product_reviews where product_id = ?`
		rows, err := m.DB.QueryContext(ctx, query, id)

		if err != nil {
			return err
		}
		defer rows.Close()
		for rows.Next() {
			r := model.ProductReview{}
			if err := rows.Scan(
				&r.ID,
				&r.ReviewerID,
				&r.ReviewerName,
				&r.ProductID,
				&r.Body,
				&r.Rating,
				&r.CreatedAt,
				&r.UpdatedAt,
			); err != nil {
				if err == sql.ErrNoRows {
					return err
				}
				return fmt.Errorf("db getproductbyid: %v", err)
			}
			p.Reviews = append(p.Reviews, r)
		}
		if err := rows.Err(); err != nil {
			return err
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			return nil
		}
	})

	g.Go(func() error {
		query := `select url from images where product_id = ?`
		rows, err := m.DB.QueryContext(ctx, query, id)

		if err != nil {
			return err
		}
		defer rows.Close()
		for rows.Next() {
			var imgUrl string
			if err := rows.Scan(
				&imgUrl,
			); err != nil {
				if err == sql.ErrNoRows {
					return err
				}
				return fmt.Errorf("db getproductbyid: %v", err)
			}
			p.Images = append(p.Images, imgUrl)

		}
		if err := rows.Err(); err != nil {
			return err
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			return nil
		}
	})

	return p, g.Wait()
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
		tx.Rollback()
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

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `select id from products order by id desc limit 1`

	var id int
	row := m.QueryRowContext(ctx, query)
	if err := row.Scan(&id); err != nil {
		return -1, err
	}

	return id + 1, nil
}

func (m *DBrepo) InsertProduct(p model.Product) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err := m.ExecContext(ctx, "INSERT INTO goRent.products (id,owner_id,brand,category,title,rating,description,price,created_at,updated_at) VALUES (?,?,?,?,?,?,?,?,?,?);",
		p.ID, p.OwnerID, p.Brand, p.Category, p.Title, p.Rating, p.Description, p.Price, p.CreatedAt, p.UpdatedAt)

	for _, v := range p.Images {

		err := m.InsertProductImages(p.ID, v)
		if err != nil {
			return err
		}
	}

	if err != nil {
		return fmt.Errorf("db InsertProduct: %v", err)
	}
	return nil
}

func (m *DBrepo) InsertProductImages(i int, s string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err := m.ExecContext(ctx, "INSERT INTO goRent.images (product_id, url, created_at, updated_at) VALUES (?,?,?,?);",
		i, s, time.Now(), time.Now())
	if err != nil {
		return fmt.Errorf("db InsertProductImages: %v", err)
	}
	return nil
}

// func (m *DBrepo) GetProductByID2(ctx context.Context, id int) (model.Product, error) {
// 	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
// 	defer cancel()

// 	p := model.Product{}
// 	g, ctx := errgroup.WithContext(ctx)

// 	g.Go(func() error {
// 		query := `select
// 						p.id, p.owner_id, p.brand, p.category, p.title, p.rating, p.description, p.price, p.created_at, p.updated_at
// 					from
// 						products p where id = ?
// 	`
// 		if err := m.DB.QueryRowContext(ctx, query, id).Scan(
// 			&p.ID,
// 			&p.OwnerID,
// 			&p.Brand,
// 			&p.Category,
// 			&p.Title,
// 			&p.Rating,
// 			&p.Description,
// 			&p.Price,
// 			&p.CreatedAt,
// 			&p.UpdatedAt,
// 		); err != nil {
// 			if err == sql.ErrNoRows {
// 				return err
// 			}
// 			return fmt.Errorf("db getproductbyid: %v", err)
// 		}
// 		query = `select username from users where id = ?`
// 		if err := m.DB.QueryRowContext(ctx, query, p.OwnerID).Scan(&p.OwnerName); err != nil {
// 			return err
// 		}
// 		select {
// 		case <-ctx.Done():
// 			return ctx.Err()
// 		default:
// 			return nil
// 		}
// 	})

// 	// get reviews from reviews table
// 	g.Go(func() error {
// 		query := `select id, reviewer_id, reviewer_name, product_id, body, rating, created_at, updated_at
// 		from product_reviews where product_id = ?`
// 		rows, err := m.DB.QueryContext(ctx, query, id)

// 		if err != nil {
// 			return err
// 		}
// 		defer rows.Close()
// 		for rows.Next() {
// 			r := model.ProductReview{}
// 			if err := rows.Scan(
// 				&r.ID,
// 				&r.ReviewerID,
// 				&r.ReviewerName,
// 				&r.ProductID,
// 				&r.Body,
// 				&r.Rating,
// 				&r.CreatedAt,
// 				&r.UpdatedAt,
// 			); err != nil {
// 				if err == sql.ErrNoRows {
// 					return err
// 				}
// 				return fmt.Errorf("db getproductbyid: %v", err)
// 			}
// 			p.Reviews = append(p.Reviews, r)
// 		}
// 		if err := rows.Err(); err != nil {
// 			return err
// 		}
// 		select {
// 		case <-ctx.Done():
// 			return ctx.Err()
// 		default:
// 			return nil
// 		}
// 	})

// 	g.Go(func() error {
// 		query := `select url from images where product_id = ?`
// 		rows, err := m.DB.QueryContext(ctx, query, id)

// 		if err != nil {
// 			return err
// 		}
// 		defer rows.Close()
// 		for rows.Next() {
// 			var imgUrl string
// 			if err := rows.Scan(
// 				&imgUrl,
// 			); err != nil {
// 				if err == sql.ErrNoRows {
// 					return err
// 				}
// 				return fmt.Errorf("db getproductbyid: %v", err)
// 			}
// 			p.Images = append(p.Images, imgUrl)

// 		}
// 		if err := rows.Err(); err != nil {
// 			return err
// 		}
// 		select {
// 		case <-ctx.Done():
// 			return ctx.Err()
// 		default:
// 			return nil
// 		}
// 	})

// 	return p, g.Wait()
// }

func (m *DBrepo) UpdateProducts(p model.Product) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// query := `UPDATE goRent.products set processed = true, updated_at = ? where id = ?`
	query := `UPDATE goRent.products set brand=?, title=?, description=?, price=?, created_at=? where id = ?`

	// _, err := m.ExecContext(ctx, "UPDATE goRent.products (brand,title,description,price, created_at) where id =? VALUES (?,?,?,?,?);",
	_, err := m.ExecContext(ctx, query, p.Brand, p.Title, p.Description, p.Price, time.Now(), p.ID)

	// for _, v := range p.Images {

	// 	err := m.InsertProductImages(p.ID, v)
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	if err != nil {
		return fmt.Errorf("db InsertProduct: %v", err)
	}
	return nil
}
