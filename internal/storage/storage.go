package storage

import (
	_ "embed"
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Storage struct {
	pg *pgxpool.Pool
}

var (
	ErrNotFound = errors.New("not found")
)

func New(pg *pgxpool.Pool) *Storage {
	s := &Storage{
		pg: pg,
	}

	return s
}
