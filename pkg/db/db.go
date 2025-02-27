//go:generate mockgen -write_package_comment=false -source=db.go -destination=db_mock.go -package db

// Package db implements database related operations.
package db

import (
	"context"
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"
)

// Repo defines interface for common database operations.
type Repo[T any] interface {
	// Insert inserts a record and updates the ID field of the given data object based on returned ID.
	Insert(ctx context.Context, do *T) error
	// QueryByID queries record by ID. If no record is found, return [sql.ErrNoRows], otherwise it will return an error
	// that may occur during execution or nil.
	QueryByID(ctx context.Context, id int64) (*T, error)
	// Update updates a record.
	Update(ctx context.Context, do *T) error
	// Delete deletes a record.
	Delete(ctx context.Context, do *T) error
	// BeginTx begins a transaction.
	BeginTx(ctx context.Context, opts *sql.TxOptions) (*sqlx.Tx, error)
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

// NewPool initializes a new database connection pool.
func NewPool(driver, dsn string) (pool *sqlx.DB, cleanup func() error, err error) {
	pool, err = sqlx.Open(driver, dsn)
	if err != nil {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(3)*time.Second) // nolint:mnd
	defer cancel()
	err = pool.PingContext(ctx)
	cleanup = pool.Close
	return
}
