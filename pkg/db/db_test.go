// -------------------------------------------------------------------------------------------
// Copyright (c) Team Sorghum. All rights reserved.
// Licensed under the GPL v3 License. See LICENSE in the project root for license information.
// -------------------------------------------------------------------------------------------

package db_test

import (
	"errors"
	"testing"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/teamsorghum/go-common/pkg/constant"
	"github.com/teamsorghum/go-common/pkg/db"
)

func TestNewPool(t *testing.T) {
	t.Parallel()

	t.Run("Success", func(t *testing.T) {
		t.Parallel()

		pool, cleanup, err := db.NewPool(&db.Config{
			Driver: "pgx",
			DSN:    "postgres://teamsorghum:teamsorghum@localhost:5432/test",
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
			DSN:    "postgres://teamsorghum:teamsorghum@localhost:5432/test",
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
