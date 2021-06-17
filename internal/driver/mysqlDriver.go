// Internal driver package to connect to postgresDB
package driver

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// DB struct for generic usage
type DB struct {
	SQL *sql.DB
}

const (
	maxOpenDbConn = 10
	maxIdleConn   = 5
	maxDbLifeTime = 5 * time.Second
)

// Connect creates database pool for MySQL
func Connect(dsn string) (*DB, error) {
	dbConn := &DB{}

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		err = fmt.Errorf("ConnectSQL: %w", err)
		return dbConn, err
	}

	db.SetMaxOpenConns(maxOpenDbConn)
	db.SetMaxIdleConns(maxIdleConn)
	db.SetConnMaxLifetime(maxDbLifeTime)

	dbConn.SQL = db
	if err = testDB(dbConn.SQL); err != nil {
		err = fmt.Errorf("ConnectSQL: %w", err)
		return dbConn, err
	}

	return dbConn, nil
}

// testDB tries to ping the database
func testDB(db *sql.DB) error {
	if err := db.Ping(); err != nil {
		return fmt.Errorf("ping: %w", err)
	}
	return nil
}
