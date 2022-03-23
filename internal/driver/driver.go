package driver

import (
	"database/sql"
	"time"
)

type DB struct {
	SQL *sql.DB
}

var dbConn = &DB{}

//membuat database pool
func ConnectSQL(dsn string) (*DB, error) {
	d, err := NewDatabase(dsn)
	if err != nil {
		panic(err)
	}

	d.SetMaxOpenConns(10)
	d.SetMaxIdleConns(5)
	d.SetConnMaxLifetime(5 * time.Minute)

	dbConn.SQL = d

	err = testDB(d)
	if err != nil {
		return nil, err
	}

	return dbConn, nil
}

//membuat database baru untuk app
func NewDatabase(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)

	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

//ping database
func testDB(d *sql.DB) error {
	err := d.Ping()

	if err != nil {
		return err
	}

	return nil
}
