//go:generate mockgen -write_package_comment=false -source=db.go -destination=db_mock.go -package db

// Package db implements database related operations.
package db

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/jmoiron/sqlx"
)

// Repo defines interface for common database operations.
type Repo[T any] interface {
	/*
		Insert inserts a record into database and updates the ID field of the data object based on returned ID.

		Params:
			- ctx context.Context: The context of request.
			- do *T: The data object to be inserted.

		Returns:
			- error: The error.
	*/
	Insert(ctx context.Context, do *T) error
	/*
		QueryByID queries record by ID.

		Params:
			- ctx context.Context: The context of request.
			- id int64: The ID to be queried by.

		Returns:
			- *T: Query result. Nil will be returned if the expected record is not found.
			- error: If no record is found, return ErrNoRows, otherwise it will return an error that may occur during
			execution or nil.
	*/
	QueryByID(ctx context.Context, id int64) (*T, error)
	/*
		Update updates a record.

		Params:
			- ctx context.Context: The context of request.
			- do *T: The data object to be updated.

		Returns:
			- error: The error.
	*/
	Update(ctx context.Context, do *T) error
	/*
		Delete deletes a record.

		Params:
			- ctx context.Context: The context of request.
			- do *T: The data object to be deleted.

		Returns:
			- error: The error.
	*/
	Delete(ctx context.Context, do *T) error
}

// DO defines a common data object. You should embed this struct in your own data object.
type DO struct {
	// ID is the primary key.
	ID int64 `db:"id"`
	// CreateTime is the create time of record.
	CreateTime time.Time `db:"create_time"`
	// UpdateTime is the update time of record.
	UpdateTime time.Time `db:"update_time"`
	// Ext is the extension of record, it should have a json type in db.
	Ext string `db:"ext"`
}

// NewPool initializes a new database connection pool, and returns a connection pool, a map that contains prepared
// statements, a cleanup function and an error. The map should have a key of type string, and a value of type
// [*sqlx.Stmts].
func NewPool(driver, dsn string) (pool *sqlx.DB, stmts *sync.Map, cleanup func() error, err error) {
	pool, err = sqlx.Open(driver, dsn)
	if err != nil {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(3)*time.Second) // nolint:mnd
	defer cancel()
	err = pool.PingContext(ctx)
	cleanup = func() error {
		e := error(nil)
		for _, stmt := range stmts.Range {
			s, ok := stmt.(*sqlx.Stmt)
			if !ok {
				return errors.New("wrong stmt type")
			}
			e = errors.Join(e, s.Close())
		}
		return errors.Join(e, pool.Close())
	}
	return
}
