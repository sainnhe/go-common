// -------------------------------------------------------------------------------------------
// Copyright (c) Team Sorghum. All rights reserved.
// Licensed under the GPL v3 License. See LICENSE in the project root for license information.
// -------------------------------------------------------------------------------------------

package db_test

import (
	"testing"

	"github.com/teamsorghum/go-common/pkg/db"
)

func TestNewStmtBuilder(t *testing.T) {
	t.Parallel()

	if db.NewStmtBuilder("", db.MySQL) != nil {
		t.Fatalf("Expect nil")
	}

	if db.NewStmtBuilder("my_tbl", db.Type(3)) != nil {
		t.Fatalf("Expect nil")
	}
}

func TestBuildMappedInsertSQL(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		tbl            string
		cols           []db.KV
		wantMySQL      string
		wantPostgreSQL string
		wantSQLite     string
	}{
		{
			name: "Single column",
			tbl:  "products",
			cols: []db.KV{
				{"name", "'product'"},
			},
			wantMySQL:      "INSERT INTO products (`name`) VALUES ('product')",
			wantPostgreSQL: "INSERT INTO products (\"name\") VALUES ('product')",
			wantSQLite:     "INSERT INTO products (\"name\") VALUES ('product')",
		},
		{
			name: "Multiple columns",
			tbl:  "users",
			cols: []db.KV{
				{"email", "?"},
				{"age", "20"},
				{"username", "?"},
			},
			wantMySQL:      "INSERT INTO users (`email`, `age`, `username`) VALUES (?, 20, ?)",
			wantPostgreSQL: "INSERT INTO users (\"email\", \"age\", \"username\") VALUES ($1, 20, $2)",
			wantSQLite:     "INSERT INTO users (\"email\", \"age\", \"username\") VALUES (?, 20, ?)",
		},
		{
			name:           "Empty columns",
			tbl:            "test",
			cols:           []db.KV{},
			wantMySQL:      "",
			wantPostgreSQL: "",
			wantSQLite:     "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mysqlBuilder := db.NewStmtBuilder(tt.tbl, db.MySQL)
			postgresqlBuilder := db.NewStmtBuilder(tt.tbl, db.PostgreSQL)
			sqliteBuilder := db.NewStmtBuilder(tt.tbl, db.SQLite)

			if s := mysqlBuilder.BuildMappedInsertSQL(tt.cols); s != tt.wantMySQL {
				t.Fatalf("Want %s\nGot %s", tt.wantMySQL, s)
			}
			if s := postgresqlBuilder.BuildMappedInsertSQL(tt.cols); s != tt.wantPostgreSQL {
				t.Fatalf("Want %s\nGot %s", tt.wantPostgreSQL, s)
			}
			if s := sqliteBuilder.BuildMappedInsertSQL(tt.cols); s != tt.wantSQLite {
				t.Fatalf("Want %s\nGot %s", tt.wantSQLite, s)
			}
		})
	}
}

func TestBuildMappedQuerySQL(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		tbl            string
		selectedCols   []string
		conds          []db.KV
		wantMySQL      string
		wantPostgreSQL string
		wantSQLite     string
	}{
		{
			name:           "Single column and condition",
			tbl:            "products",
			selectedCols:   []string{"username"},
			conds:          []db.KV{{"id", "5"}},
			wantMySQL:      "SELECT `username` FROM products WHERE id = 5",
			wantPostgreSQL: "SELECT \"username\" FROM products WHERE id = 5",
			wantSQLite:     "SELECT \"username\" FROM products WHERE id = 5",
		},
		{
			name:         "Multiple columns and conditions",
			tbl:          "users",
			selectedCols: []string{"username", "nickname"},
			conds: []db.KV{
				{"name", "?"},
				{"age", "20"},
				{"gender", "?"},
			},
			wantMySQL:      "SELECT `username`, `nickname` FROM users WHERE name = ? AND age = 20 AND gender = ?",
			wantPostgreSQL: "SELECT \"username\", \"nickname\" FROM users WHERE name = $1 AND age = 20 AND gender = $2",
			wantSQLite:     "SELECT \"username\", \"nickname\" FROM users WHERE name = ? AND age = 20 AND gender = ?",
		},
		{
			name:           "No columns and conditions",
			tbl:            "orders",
			selectedCols:   []string{},
			conds:          []db.KV{},
			wantMySQL:      "SELECT * FROM orders",
			wantPostgreSQL: "SELECT * FROM orders",
			wantSQLite:     "SELECT * FROM orders",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mysqlBuilder := db.NewStmtBuilder(tt.tbl, db.MySQL)
			postgresqlBuilder := db.NewStmtBuilder(tt.tbl, db.PostgreSQL)
			sqliteBuilder := db.NewStmtBuilder(tt.tbl, db.SQLite)

			if s := mysqlBuilder.BuildMappedQuerySQL(tt.selectedCols, tt.conds); s != tt.wantMySQL {
				t.Fatalf("Want %s\nGot %s", tt.wantMySQL, s)
			}
			if s := postgresqlBuilder.BuildMappedQuerySQL(tt.selectedCols, tt.conds); s != tt.wantPostgreSQL {
				t.Fatalf("Want %s\nGot %s", tt.wantPostgreSQL, s)
			}
			if s := sqliteBuilder.BuildMappedQuerySQL(tt.selectedCols, tt.conds); s != tt.wantSQLite {
				t.Fatalf("Want %s\nGot %s", tt.wantSQLite, s)
			}
		})
	}
}

func TestBuildMappedUpdateSQL(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		tbl            string
		cols           []db.KV
		conds          []db.KV
		wantMySQL      string
		wantPostgreSQL string
		wantSQLite     string
	}{
		{
			name: "Multiple cols AND conds",
			tbl:  "users",
			cols: []db.KV{
				{"age", "20"},
				{"username", "?"},
				{"nickname", "?"},
			},
			conds: []db.KV{
				{"id", "?"},
				{"status", "'active'"},
			},
			wantMySQL:      "UPDATE users SET `age` = 20, `username` = ?, `nickname` = ? WHERE id = ? AND status = 'active'",
			wantPostgreSQL: "UPDATE users SET \"age\" = 20, \"username\" = $1, \"nickname\" = $2 WHERE id = $3 AND status = 'active'", // nolint:lll
			wantSQLite:     "UPDATE users SET \"age\" = 20, \"username\" = ?, \"nickname\" = ? WHERE id = ? AND status = 'active'",    // nolint:lll
		},
		{
			name:           "Empty cols",
			tbl:            "test",
			cols:           []db.KV{},
			conds:          []db.KV{{"id", "1"}},
			wantMySQL:      "",
			wantPostgreSQL: "",
			wantSQLite:     "",
		},
		{
			name:           "Empty conds",
			tbl:            "test",
			cols:           []db.KV{{"name", "'test'"}},
			conds:          []db.KV{},
			wantMySQL:      "UPDATE test SET `name` = 'test'",
			wantPostgreSQL: "UPDATE test SET \"name\" = 'test'",
			wantSQLite:     "UPDATE test SET \"name\" = 'test'",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mysqlBuilder := db.NewStmtBuilder(tt.tbl, db.MySQL)
			postgresqlBuilder := db.NewStmtBuilder(tt.tbl, db.PostgreSQL)
			sqliteBuilder := db.NewStmtBuilder(tt.tbl, db.SQLite)

			if s := mysqlBuilder.BuildMappedUpdateSQL(tt.cols, tt.conds); s != tt.wantMySQL {
				t.Fatalf("Want %s\nGot %s", tt.wantMySQL, s)
			}
			if s := postgresqlBuilder.BuildMappedUpdateSQL(tt.cols, tt.conds); s != tt.wantPostgreSQL {
				t.Fatalf("Want %s\nGot %s", tt.wantPostgreSQL, s)
			}
			if s := sqliteBuilder.BuildMappedUpdateSQL(tt.cols, tt.conds); s != tt.wantSQLite {
				t.Fatalf("Want %s\nGot %s", tt.wantSQLite, s)
			}
		})
	}
}

func TestBuildMappedDeleteSQL(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		tbl            string
		conds          []db.KV
		wantMySQL      string
		wantPostgreSQL string
		wantSQLite     string
	}{
		{
			name: "Multiple conditions",
			tbl:  "users",
			conds: []db.KV{
				{"id", "?"},
				{"status", "'inactive'"},
				{"age", "?"},
			},
			wantMySQL:      "DELETE FROM users WHERE id = ? AND status = 'inactive' AND age = ?",
			wantPostgreSQL: "DELETE FROM users WHERE id = $1 AND status = 'inactive' AND age = $2",
			wantSQLite:     "DELETE FROM users WHERE id = ? AND status = 'inactive' AND age = ?",
		},
		{
			name:           "No conditions",
			tbl:            "orders",
			conds:          []db.KV{},
			wantMySQL:      "DELETE FROM orders",
			wantPostgreSQL: "DELETE FROM orders",
			wantSQLite:     "DELETE FROM orders",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mysqlBuilder := db.NewStmtBuilder(tt.tbl, db.MySQL)
			postgresqlBuilder := db.NewStmtBuilder(tt.tbl, db.PostgreSQL)
			sqliteBuilder := db.NewStmtBuilder(tt.tbl, db.SQLite)

			if s := mysqlBuilder.BuildMappedDeleteSQL(tt.conds); s != tt.wantMySQL {
				t.Fatalf("Want %s\nGot %s", tt.wantMySQL, s)
			}
			if s := postgresqlBuilder.BuildMappedDeleteSQL(tt.conds); s != tt.wantPostgreSQL {
				t.Fatalf("Want %s\nGot %s", tt.wantPostgreSQL, s)
			}
			if s := sqliteBuilder.BuildMappedDeleteSQL(tt.conds); s != tt.wantSQLite {
				t.Fatalf("Want %s\nGot %s", tt.wantSQLite, s)
			}
		})
	}
}

func TestBuildNamedInsertSQL(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name            string
		tbl             string
		cols            []string
		wantMySQL       string
		wantPgAndSqlite string
	}{
		{
			name:            "Single column",
			tbl:             "products",
			cols:            []string{"name"},
			wantMySQL:       "INSERT INTO products (`name`) VALUES (:name)",
			wantPgAndSqlite: "INSERT INTO products (\"name\") VALUES (:name)",
		},
		{
			name:            "Multiple columns",
			tbl:             "users",
			cols:            []string{"username", "age"},
			wantMySQL:       "INSERT INTO users (`username`, `age`) VALUES (:username, :age)",
			wantPgAndSqlite: "INSERT INTO users (\"username\", \"age\") VALUES (:username, :age)",
		},
		{
			name:            "Empty columns",
			tbl:             "test",
			cols:            []string{},
			wantMySQL:       "",
			wantPgAndSqlite: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mysqlBuilder := db.NewStmtBuilder(tt.tbl, db.MySQL)
			postgresqlBuilder := db.NewStmtBuilder(tt.tbl, db.PostgreSQL)
			sqliteBuilder := db.NewStmtBuilder(tt.tbl, db.SQLite)

			if s := mysqlBuilder.BuildNamedInsertSQL(tt.cols); s != tt.wantMySQL {
				t.Fatalf("Want %s\nGot %s", tt.wantMySQL, s)
			}
			if s := postgresqlBuilder.BuildNamedInsertSQL(tt.cols); s != tt.wantPgAndSqlite {
				t.Fatalf("Want %s\nGot %s", tt.wantPgAndSqlite, s)
			}
			if s := sqliteBuilder.BuildNamedInsertSQL(tt.cols); s != tt.wantPgAndSqlite {
				t.Fatalf("Want %s\nGot %s", tt.wantPgAndSqlite, s)
			}
		})
	}
}

func TestBuildNamedQuerySQL(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name            string
		tbl             string
		selectedCols    []string
		conds           []string
		wantMySQL       string
		wantPgAndSqlite string
	}{
		{
			name:            "Multiple conditions and columns",
			tbl:             "users",
			selectedCols:    []string{"nickname", "gender"},
			conds:           []string{"username", "age"},
			wantMySQL:       "SELECT `nickname`, `gender` FROM users WHERE `username` = :username AND `age` = :age",
			wantPgAndSqlite: "SELECT \"nickname\", \"gender\" FROM users WHERE \"username\" = :username AND \"age\" = :age",
		},
		{
			name:            "No conditions and columns",
			tbl:             "test",
			selectedCols:    []string{},
			conds:           []string{},
			wantMySQL:       "SELECT * FROM test",
			wantPgAndSqlite: "SELECT * FROM test",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mysqlBuilder := db.NewStmtBuilder(tt.tbl, db.MySQL)
			postgresqlBuilder := db.NewStmtBuilder(tt.tbl, db.PostgreSQL)
			sqliteBuilder := db.NewStmtBuilder(tt.tbl, db.SQLite)

			if s := mysqlBuilder.BuildNamedQuerySQL(tt.selectedCols, tt.conds); s != tt.wantMySQL {
				t.Fatalf("Want %s\nGot %s", tt.wantMySQL, s)
			}
			if s := postgresqlBuilder.BuildNamedQuerySQL(tt.selectedCols, tt.conds); s != tt.wantPgAndSqlite {
				t.Fatalf("Want %s\nGot %s", tt.wantPgAndSqlite, s)
			}
			if s := sqliteBuilder.BuildNamedQuerySQL(tt.selectedCols, tt.conds); s != tt.wantPgAndSqlite {
				t.Fatalf("Want %s\nGot %s", tt.wantPgAndSqlite, s)
			}
		})
	}
}

func TestBuildNamedUpdateSQL(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name            string
		tbl             string
		cols            []string
		conds           []string
		wantMySQL       string
		wantPgAndSqlite string
	}{
		{
			name:            "Multiple cols AND conds",
			tbl:             "users",
			cols:            []string{"username", "age"},
			conds:           []string{"id", "status"},
			wantMySQL:       "UPDATE users SET `username` = :username, `age` = :age WHERE `id` = :id AND `status` = :status",
			wantPgAndSqlite: "UPDATE users SET \"username\" = :username, \"age\" = :age WHERE \"id\" = :id AND \"status\" = :status", // nolint:lll
		},
		{
			name:            "Empty cols",
			tbl:             "test",
			cols:            []string{},
			conds:           []string{"id"},
			wantMySQL:       "",
			wantPgAndSqlite: "",
		},
		{
			name:            "Empty conds",
			tbl:             "test",
			cols:            []string{"name"},
			conds:           []string{},
			wantMySQL:       "UPDATE test SET `name` = :name",
			wantPgAndSqlite: "UPDATE test SET \"name\" = :name",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mysqlBuilder := db.NewStmtBuilder(tt.tbl, db.MySQL)
			postgresqlBuilder := db.NewStmtBuilder(tt.tbl, db.PostgreSQL)
			sqliteBuilder := db.NewStmtBuilder(tt.tbl, db.SQLite)

			if s := mysqlBuilder.BuildNamedUpdateSQL(tt.cols, tt.conds); s != tt.wantMySQL {
				t.Fatalf("Want %s\nGot %s", tt.wantMySQL, s)
			}
			if s := postgresqlBuilder.BuildNamedUpdateSQL(tt.cols, tt.conds); s != tt.wantPgAndSqlite {
				t.Fatalf("Want %s\nGot %s", tt.wantPgAndSqlite, s)
			}
			if s := sqliteBuilder.BuildNamedUpdateSQL(tt.cols, tt.conds); s != tt.wantPgAndSqlite {
				t.Fatalf("Want %s\nGot %s", tt.wantPgAndSqlite, s)
			}
		})
	}
}

func TestBuildNamedDeleteSQL(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name            string
		tbl             string
		conds           []string
		wantMySQL       string
		wantPgAndSqlite string
	}{
		{
			name:            "Multiple conditions",
			tbl:             "users",
			conds:           []string{"id", "status"},
			wantMySQL:       "DELETE FROM users WHERE `id` = :id AND `status` = :status",
			wantPgAndSqlite: "DELETE FROM users WHERE \"id\" = :id AND \"status\" = :status",
		},
		{
			name:            "No conditions",
			tbl:             "test",
			conds:           []string{},
			wantMySQL:       "DELETE FROM test",
			wantPgAndSqlite: "DELETE FROM test",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mysqlBuilder := db.NewStmtBuilder(tt.tbl, db.MySQL)
			postgresqlBuilder := db.NewStmtBuilder(tt.tbl, db.PostgreSQL)
			sqliteBuilder := db.NewStmtBuilder(tt.tbl, db.SQLite)

			if s := mysqlBuilder.BuildNamedDeleteSQL(tt.conds); s != tt.wantMySQL {
				t.Fatalf("Want %s\nGot %s", tt.wantMySQL, s)
			}
			if s := postgresqlBuilder.BuildNamedDeleteSQL(tt.conds); s != tt.wantPgAndSqlite {
				t.Fatalf("Want %s\nGot %s", tt.wantPgAndSqlite, s)
			}
			if s := sqliteBuilder.BuildNamedDeleteSQL(tt.conds); s != tt.wantPgAndSqlite {
				t.Fatalf("Want %s\nGot %s", tt.wantPgAndSqlite, s)
			}
		})
	}
}
