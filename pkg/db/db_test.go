package db_test

import (
	"errors"
	"testing"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/sainnhe/go-common/pkg/constant"
	"github.com/sainnhe/go-common/pkg/db"
)

func TestNewPool(t *testing.T) {
	t.Parallel()

	t.Run("Success", func(t *testing.T) {
		t.Parallel()

		pool, cleanup, err := db.NewPool(&db.Config{
			Driver: "pgx",
			DSN:    "postgres://sainnhe:sainnhe@localhost:5432/test",
		})
		defer cleanup()
		if pool == nil {
			t.Fatal("pool is nil")
		}
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("Nil config", func(t *testing.T) {
		t.Parallel()

		_, _, err := db.NewPool(nil)
		if !errors.Is(err, constant.ErrNilDeps) {
			t.Fatalf("Expect error %+v, got %+v", constant.ErrNilDeps, err)
		}
	})

	t.Run("Driver doesn't exist.", func(t *testing.T) {
		t.Parallel()

		_, _, err := db.NewPool(&db.Config{
			Driver: "pg",
			DSN:    "postgres://sainnhe:sainnhe@localhost:5432/test",
		})
		if err == nil {
			t.Fatal("Expect error, got nil.")
		}
	})

	t.Run("DSN invalid", func(t *testing.T) {
		t.Parallel()

		_, _, err := db.NewPool(&db.Config{
			Driver: "pg",
			DSN:    "postgres://foo:bar@localhost:5432/test",
		})
		if err == nil {
			t.Fatal("Expect error, got nil.")
		}
	})
}
