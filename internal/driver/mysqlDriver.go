// Internal driver package to connect to postgresDB
package driver

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// DB struct for generic usage

const (
	maxOpenDbConn = 10
	maxIdleConn   = 5
	maxDbLifeTime = 5 * time.Second
)

// Connect creates database pool for MySQL
func Connect(dsn string, dialect string) (*sql.DB, error) {

	db, err := sql.Open(dialect, dsn)
	if err != nil {
		err = fmt.Errorf("ConnectSQL: %w", err)
		return nil, err
	}

	db.SetMaxOpenConns(maxOpenDbConn)
	db.SetMaxIdleConns(maxIdleConn)
	db.SetConnMaxLifetime(maxDbLifeTime)

	if err = testDB(db); err != nil {
		err = fmt.Errorf("ConnectSQL: %w", err)
		return nil, err
	}

	return db, nil
}

// testDB tries to ping the database
func testDB(db *sql.DB) error {
	if err := db.Ping(); err != nil {
		return fmt.Errorf("ping: %w", err)
	}
	return nil
}
