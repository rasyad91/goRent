//Users - all the functions for the users
package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"goRent/internal/model"
	"time"

	"golang.org/x/sync/errgroup"
)

//check to see if email exist in the system
func (m *dbRepo) EmailExist(e string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var email string
	err := m.DB.QueryRowContext(ctx, "SELECT * FROM gorent.users where email=?", e).Scan(
		&email,
	)
	if err == sql.ErrNoRows {
		return nil
	}
	if err != nil {
		return err
	}
	return nil
}

//fill in User struct
func (m *dbRepo) GetUser(username string) (model.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	u := model.User{}
	if err := m.DB.QueryRowContext(ctx, "SELECT id,username,email,image_url,password,access_level,rating,postal_code,street_name, block,unit_number,created_at, updated_at FROM gorent.users where username=?", username).
		Scan(
			&u.ID,
			&u.Username,
			&u.Email,
			&u.Image_URL,
			&u.Password,
			&u.AccessLevel,
			&u.Rating,
			&u.Address.PostalCode,
			&u.Address.StreetName,
			&u.Address.Block,
			&u.Address.UnitNumber,
			&u.CreatedAt,
			&u.UpdatedAt,
		); err != nil {
		if err == sql.ErrNoRows {
			return model.User{}, err
		}
		return model.User{}, fmt.Errorf("db GetUser: %v", err)
	}
	// add in concurrency here - 1 get Rent 1 get Product
	// get rents
	rent_query := `select
			r.id, r.owner_id, r.renter_id, r.product_id, r.restriction_id, r.processed, r.start_date, r.end_date, r.duration, r.total_cost, r.created_at, r.updated_at,
			p.id, p.owner_id, p.brand, p.category, p.title, p.rating, p.description, p.price, p.created_at, p.updated_at, i.url
		from
			rents r
		left join
			products p on (p.id = r.product_id)
		left join (select product_id,min(url) url from images group by 1) i on p.id = i.product_id
		where
			r.renter_id = ?
		order by r.product_id asc`
	product_query := `select p.*, i.url from products p 
		left join (select product_id,min(url) url from images group by 1) i on p.id = i.product_id 
		where p.owner_id = ? order by id asc`
	booking_query := `select
		r.id, r.owner_id, r.renter_id, r.product_id, r.restriction_id, r.processed, r.start_date, r.end_date, r.duration, r.total_cost, r.created_at, r.updated_at,
		p.id, p.owner_id, p.brand, p.category, p.title, p.rating, p.description, p.price, p.created_at, p.updated_at, i.url
	from
		rents r
	left join
		products p on (p.id = r.product_id)
	left join (select product_id,min(url) url from images group by 1) i on p.id = i.product_id
	where
		r.owner_id = ?
	order by r.product_id asc`
	//initializing concurrency // linear - 9.791375ms, concurrent - 7.357559ms
	//timing prior to concurrency
	x, ctx := errgroup.WithContext(ctx)
	x.Go(func() error {
		if err := m.runQuery(ctx, &u, rent_query, "rent"); err != nil {
			return err
		}
		return nil
	})
	x.Go(func() error {
		if err := m.runQuery(ctx, &u, product_query, "product"); err != nil {
			return err
		}
		return nil
	})
	x.Go(func() error {
		if err := m.runQuery(ctx, &u, booking_query, "booking"); err != nil {
			return err
		}
		return nil
	})

	return u, x.Wait()
}

//for GetUser to pass in query and fill in the User struct
func (m *dbRepo) runQuery(ctx context.Context, user *model.User, query string, structType string) error {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)

	defer cancel()
	rows, err := m.DB.QueryContext(ctx, query, user.ID)
	if err != nil {
		return fmt.Errorf("db GetUser %s: %v", structType, err)
	}
	defer rows.Close()
	var urlString string
	for rows.Next() {
		if structType == "rent" {
			r := model.Rent{}
			if err := rows.Scan(
				&r.ID,
				&r.OwnerID,
				&r.RenterID,
				&r.ProductID,
				&r.RestrictionID,
				&r.Processed,
				&r.StartDate,
				&r.EndDate,
				&r.Duration,
				&r.TotalCost,
				&r.CreatedAt,
				&r.UpdatedAt,
				&r.Product.ID,
				&r.Product.OwnerID,
				&r.Product.Brand,
				&r.Product.Category,
				&r.Product.Title,
				&r.Product.Rating,
				&r.Product.Description,
				&r.Product.Price,
				&r.Product.CreatedAt,
				&r.Product.UpdatedAt,
				&urlString,
			); err != nil {
				if err == sql.ErrNoRows {
					return err
				}
				return fmt.Errorf("db GetUser %s: %v", structType, err)
			}
			r.Product.Images = append(r.Product.Images, urlString)
			user.Rents = append(user.Rents, r)
		} else if structType == "booking" {
			r := model.Rent{}
			if err := rows.Scan(
				&r.ID,
				&r.OwnerID,
				&r.RenterID,
				&r.ProductID,
				&r.RestrictionID,
				&r.Processed,
				&r.StartDate,
				&r.EndDate,
				&r.Duration,
				&r.TotalCost,
				&r.CreatedAt,
				&r.UpdatedAt,
				&r.Product.ID,
				&r.Product.OwnerID,
				&r.Product.Brand,
				&r.Product.Category,
				&r.Product.Title,
				&r.Product.Rating,
				&r.Product.Description,
				&r.Product.Price,
				&r.Product.CreatedAt,
				&r.Product.UpdatedAt,
				&urlString,
			); err != nil {
				if err == sql.ErrNoRows {
					return err
				}
				return fmt.Errorf("db GetUser %s: %v", structType, err)
			}
			r.Product.Images = append(r.Product.Images, urlString)
			user.Bookings = append(user.Bookings, r)
		} else {
			r := model.Product{}
			if err := rows.Scan(
				&r.ID,
				&r.OwnerID,
				&r.Brand,
				&r.Category,
				&r.Title,
				&r.Rating,
				&r.Description,
				&r.Price,
				// &p.Images,
				&r.CreatedAt,
				&r.UpdatedAt,
				&urlString,
			); err != nil {
				if err == sql.ErrNoRows {
					return err
				}
				return fmt.Errorf("db GetUser %s: %v", structType, err)
			}
			r.Images = append(r.Images, urlString)
			user.Products = append(user.Products, r)
		}
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		return nil
	}
}

//Insert User into the DB
func (m *dbRepo) InsertUser(u model.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err := m.ExecContext(ctx, "INSERT INTO gorent.users (username,email,password,image_url,postal_code,street_name,block,unit_number,created_at,updated_at) VALUES (?,?,?,?,?,?,?,?,?,?);",
		u.Username, u.Email, u.Password, u.Image_URL,
		u.Address.PostalCode, u.Address.StreetName, u.Address.Block, u.Address.UnitNumber,
		time.Now(), time.Now())
	if err != nil {
		return fmt.Errorf("db InsertUser: %v", err)
	}
	return nil
}

//Edit Profile base on which form returns
func (m *dbRepo) EditUser(u model.User, editType string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := error(nil)
	if editType == "address" {
		_, err = m.ExecContext(ctx, "UPDATE gorent.users SET block = ?, street_name = ?, unit_number = ?, postal_code = ? WHERE id = ?", u.Address.Block, u.Address.StreetName, u.Address.UnitNumber, u.Address.PostalCode, u.ID)
		fmt.Println("AddressChange test:", u)
	} else if editType == "profile" {
		_, err = m.ExecContext(ctx, "UPDATE gorent.users SET username = ?, email = ? WHERE id = ?", u.Username, u.Email, u.ID)
		fmt.Println("ProfileChange test:", u)
	} else if editType == "profileImage" {
		_, err = m.ExecContext(ctx, "UPDATE gorent.users SET image_url = ? WHERE id = ?", u.Image_URL, u.ID)
		fmt.Println("ProfileImage test:", u)
	} else {
		_, err = m.ExecContext(ctx, "UPDATE gorent.users SET password = ? WHERE id = ?", u.Password, u.ID)
		fmt.Println("PassWordChange test:", u)
	}
	if err != nil {
		return fmt.Errorf("db EditUser: %v", err)
	}
	return nil
}
