package storage

import "database/sql"

type Storage struct {
	DB *sql.DB
}

func New(db *sql.DB) *Storage {
	return &Storage{
		DB: db,
	}
}
