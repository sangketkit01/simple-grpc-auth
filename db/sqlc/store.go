package db

import (
	"database/sql"
	"time"
)

type Store struct {
	Queries
	db *sql.DB
}

func NewStore(db *sql.DB) (*Store, error) {
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)
	db.SetConnMaxIdleTime(1 * time.Minute)
	return &Store{
		db: db,
		Queries: *New(db),
	}, nil
}
