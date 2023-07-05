package saver

import (
	"context"
	"database/sql"
	entitiesDB "github.com/akrillis/nakamatest/internal/entities/db"
)

type DB struct {
	ctx context.Context
	db  *sql.DB
}

func NewDB(ctx context.Context, db *sql.DB) *DB {
	return &DB{
		ctx: ctx,
		db:  db,
	}
}

func (d *DB) Save(path, content string) error {
	_, err := d.db.ExecContext(d.ctx, entitiesDB.StmtInsert, path, content)
	if err != nil {
		return err
	}

	return nil
}
