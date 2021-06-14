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

	u := model.User{}
	if err := m.DB.QueryRowContext(ctx, "SELECT id,username,image_urlaccess_level FROM gorent.users").
		Scan(
			&u.ID,
			&u.Username,
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
}
