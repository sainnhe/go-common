// -------------------------------------------------------------------------------------------
// Copyright (c) Team Sorghum. All rights reserved.
// Licensed under the GPL v3 License. See LICENSE in the project root for license information.
// -------------------------------------------------------------------------------------------

package db_test

import (
	"testing"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/teamsorghum/go-common/pkg/db"
)

func TestNewPool(t *testing.T) {
	t.Parallel()

	t.Run("Success", func(t *testing.T) {
		t.Parallel()

		pool, cleanup, err := db.NewPool("pgx", "postgres://teamsorghum:teamsorghum@localhost:5432/test")
		defer cleanup()
		if pool == nil {
			t.Fatal("pool is nil")
		}
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("Driver doesn't exist.", func(t *testing.T) {
		t.Parallel()

		_, _, err := db.NewPool("pg", "postgres://teamsorghum:teamsorghum@localhost:5432/test")
		if err == nil {
			t.Fatal("Expect error, got nil.")
		}
	})

	t.Run("Ping failed", func(t *testing.T) {
		t.Parallel()

		_, _, err := db.NewPool("pgx", "postgres://foo:bar@localhost:5432/test")
		if err == nil {
			t.Fatal("Expect error, got nil.")
		}
	})
}
