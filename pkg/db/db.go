//go:generate mockgen -write_package_comment=false -source=db.go -destination=db_mock.go -package db

// Package db defines database related APIs.
package db

import (
	"context"
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"
)

// Repo defines a interface for common database operations, where DO is the struct of data object.
type Repo[DO any] interface {
	// Insert inserts a record and updates the ID field of the given data object based on returned ID.
	Insert(ctx context.Context, do *DO) error

	// QueryByID queries record by ID. If no record is found, return [sql.ErrNoRows], otherwise it will return an error
	// that may occur during execution or nil.
	QueryByID(ctx context.Context, id int64) (*DO, error)

	// Update updates a record.
	Update(ctx context.Context, do *DO) error

	// Delete deletes a record.
	Delete(ctx context.Context, do *DO) error

	// BeginTx begins a transaction.
	BeginTx(ctx context.Context, opts *sql.TxOptions) (*sqlx.Tx, error)
}

// DO defines a common data object. You should embed this struct in your own data object.
type DO struct {
	// ID is the primary key.
	ID int64 `db:"id"`

	// CreateTime is the create time of a record.
	CreateTime time.Time `db:"create_time"`

	// UpdateTime is the update time of a record.
	UpdateTime time.Time `db:"update_time"`

	// Ext is the extension field. If the table structure is hard to modify after running for a period of time and you
	// want to add a new column, you can use ext as a temporary alternative. The ext field should be of type json in the
	// database.
	Ext string `db:"ext"`
}

// DOCols contains the column names of [DO].
var DOCols = []string{
	"id",
	"create_time",
	"update_time",
	"ext",
}

// NewPool initializes a new database connection pool. The driver and DSN (Data Source Name) can be found in your SQL
// driver documentation, for example [github.com/go-sql-driver/mysql].
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
