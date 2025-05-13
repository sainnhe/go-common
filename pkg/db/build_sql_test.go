// -------------------------------------------------------------------------------------------
// Copyright (c) Team Sorghum. All rights reserved.
// Licensed under the GPL v3 License. See LICENSE in the project root for license information.
// -------------------------------------------------------------------------------------------

package db_test

import (
	"testing"

	"github.com/teamsorghum/go-common/pkg/db"
)

func TestBuildMappedInsertSQL(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		tbl  string
		cols []db.KV
		want string
	}{
		{
			name: "Single column",
			tbl:  "products",
			cols: []db.KV{
				{"name", "'product'"},
			},
			want: "INSERT INTO products (name) VALUES ('product')",
		},
		{
			name: "Multiple columns",
			tbl:  "users",
			cols: []db.KV{
				{"email", "'test@example.com'"},
				{"age", "20"},
				{"username", "$1"},
			},
			want: "INSERT INTO users (email, age, username) VALUES ('test@example.com', 20, $1)",
		},
		{
			name: "Empty columns",
			tbl:  "test",
			cols: []db.KV{},
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := db.BuildMappedInsertSQL(tt.tbl, tt.cols)
			if got != tt.want {
				t.Errorf("Want %q, got %q", tt.want, got)
			}
		})
	}
}

func TestBuildMappedQuerySQL(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		tbl   string
		cols  []string
		conds []db.KV
		want  string
	}{
		{
			name:  "Single column and condition",
			tbl:   "products",
			cols:  []string{"username"},
			conds: []db.KV{{"id", "5"}},
			want:  "SELECT username FROM products WHERE id = 5",
		},
		{
			name:  "Multiple columns and conditions",
			tbl:   "users",
			cols:  []string{"username", "nickname"},
			conds: []db.KV{{"age", "20"}, {"gender", "'male'"}},
			want:  "SELECT username, nickname FROM users WHERE age = 20 AND gender = 'male'",
		},
		{
			name:  "No columns and conditions",
			tbl:   "orders",
			cols:  []string{},
			conds: []db.KV{},
			want:  "SELECT * FROM orders",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := db.BuildMappedQuerySQL(tt.tbl, tt.cols, tt.conds)
			if got != tt.want {
				t.Errorf("Want %q, got %q", tt.want, got)
			}
		})
	}
}

func TestBuildMappedUpdateSQL(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		tbl   string
		cols  []db.KV
		conds []db.KV
		want  string
	}{
		{
			name: "Multiple cols AND conds",
			tbl:  "users",
			cols: []db.KV{
				{"age", "20"},
				{"username", "$1"},
			},
			conds: []db.KV{
				{"id", "5"},
				{"status", "'active'"},
			},
			want: "UPDATE users SET age = 20, username = $1 WHERE id = 5 AND status = 'active'",
		},
		{
			name:  "Empty cols",
			tbl:   "test",
			cols:  []db.KV{},
			conds: []db.KV{{"id", "1"}},
			want:  "",
		},
		{
			name:  "Empty conds",
			tbl:   "test",
			cols:  []db.KV{{"name", "'test'"}},
			conds: []db.KV{},
			want:  "UPDATE test SET name = 'test'",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := db.BuildMappedUpdateSQL(tt.tbl, tt.cols, tt.conds)
			if got != tt.want {
				t.Errorf("Want %q, got %q", tt.want, got)
			}
		})
	}
}

func TestBuildMappedDeleteSQL(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		tbl   string
		conds []db.KV
		want  string
	}{
		{
			name:  "Multiple conditions",
			tbl:   "users",
			conds: []db.KV{{"id", "5"}, {"status", "'inactive'"}},
			want:  "DELETE FROM users WHERE id = 5 AND status = 'inactive'",
		},
		{
			name:  "No conditions",
			tbl:   "orders",
			conds: []db.KV{},
			want:  "DELETE FROM orders",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := db.BuildMappedDeleteSQL(tt.tbl, tt.conds)
			if got != tt.want {
				t.Errorf("Want %q, got %q", tt.want, got)
			}
		})
	}
}

func TestBuildNamedInsertSQL(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		tbl  string
		cols []string
		want string
	}{
		{
			name: "Multiple columns",
			tbl:  "users",
			cols: []string{"username", "age"},
			want: "INSERT INTO users (username, age) VALUES (:username, :age)",
		},
		{
			name: "Empty columns",
			tbl:  "test",
			cols: []string{},
			want: "",
		},
		{
			name: "Single column",
			tbl:  "products",
			cols: []string{"name"},
			want: "INSERT INTO products (name) VALUES (:name)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := db.BuildNamedInsertSQL(tt.tbl, tt.cols)
			if got != tt.want {
				t.Errorf("Want %q, got %q", tt.want, got)
			}
		})
	}
}

func TestBuildNamedQuerySQL(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		tbl   string
		cols  []string
		conds []string
		want  string
	}{
		{
			name:  "Multiple conditions and columns",
			tbl:   "users",
			cols:  []string{"nickname", "gender"},
			conds: []string{"username", "age"},
			want:  "SELECT nickname, gender FROM users WHERE username = :username AND age = :age",
		},
		{
			name:  "No conditions and columns",
			tbl:   "test",
			cols:  []string{},
			conds: []string{},
			want:  "SELECT * FROM test",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := db.BuildNamedQuerySQL(tt.tbl, tt.cols, tt.conds)
			if got != tt.want {
				t.Errorf("Want %q, got %q", tt.want, got)
			}
		})
	}
}

func TestBuildNamedUpdateSQL(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		tbl   string
		cols  []string
		conds []string
		want  string
	}{
		{
			name:  "Multiple cols AND conds",
			tbl:   "users",
			cols:  []string{"username", "age"},
			conds: []string{"id", "status"},
			want:  "UPDATE users SET username = :username, age = :age WHERE id = :id AND status = :status",
		},
		{
			name:  "Empty cols",
			tbl:   "test",
			cols:  []string{},
			conds: []string{"id"},
			want:  "",
		},
		{
			name:  "Empty conds",
			tbl:   "test",
			cols:  []string{"name"},
			conds: []string{},
			want:  "UPDATE test SET name = :name",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := db.BuildNamedUpdateSQL(tt.tbl, tt.cols, tt.conds)
			if got != tt.want {
				t.Errorf("Want %q, got %q", tt.want, got)
			}
		})
	}
}

func TestBuildNamedDeleteSQL(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		tbl   string
		conds []string
		want  string
	}{
		{
			name:  "Multiple conditions",
			tbl:   "users",
			conds: []string{"id", "status"},
			want:  "DELETE FROM users WHERE id = :id AND status = :status",
		},
		{
			name:  "No conditions",
			tbl:   "test",
			conds: []string{},
			want:  "DELETE FROM test",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := db.BuildNamedDeleteSQL(tt.tbl, tt.conds)
			if got != tt.want {
				t.Errorf("Want %q, got %q", tt.want, got)
			}
		})
	}
}
