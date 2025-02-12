package util_test

import (
	"testing"

	"github.com/teamsorghum/go-common/pkg/util"
)

func TestBuildMappedInsertSQL(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		tbl  string
		cols []util.KV
		want string
	}{
		{
			name: "Multiple columns with sorting",
			tbl:  "users",
			cols: []util.KV{
				{"age", "20"},
				{"email", "'test@example.com'"},
				{"username", "$1"},
			},
			want: "INSERT INTO users (age, email, username) VALUES (20, 'test@example.com', $1)",
		},
		{
			name: "Empty columns",
			tbl:  "test",
			cols: []util.KV{},
			want: "INSERT INTO test () VALUES ()",
		},
		{
			name: "Single column",
			tbl:  "products",
			cols: []util.KV{
				{"name", "'product'"},
			},
			want: "INSERT INTO products (name) VALUES ('product')",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := util.BuildMappedInsertSQL(tt.tbl, tt.cols)
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
		conds []util.KV
		want  string
	}{
		{
			name:  "Multiple conditions sorted",
			tbl:   "users",
			conds: []util.KV{{"age", "20"}, {"username", "$1"}},
			want:  "SELECT * FROM users WHERE age = 20 AND username = $1",
		},
		{
			name:  "Single condition",
			tbl:   "products",
			conds: []util.KV{{"id", "5"}},
			want:  "SELECT * FROM products WHERE id = 5",
		},
		{
			name:  "No conditions",
			tbl:   "orders",
			conds: []util.KV{},
			want:  "SELECT * FROM orders",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := util.BuildMappedQuerySQL(tt.tbl, tt.conds)
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
		cols  []util.KV
		conds []util.KV
		want  string
	}{
		{
			name: "Multiple cols AND conds",
			tbl:  "users",
			cols: []util.KV{
				{"age", "20"},
				{"username", "$1"},
			},
			conds: []util.KV{
				{"id", "5"},
				{"status", "'active'"},
			},
			want: "UPDATE users SET age = 20, username = $1 WHERE id = 5 AND status = 'active'",
		},
		{
			name:  "Empty cols",
			tbl:   "test",
			cols:  []util.KV{},
			conds: []util.KV{{"id", "1"}},
			want:  "UPDATE test SET  WHERE id = 1",
		},
		{
			name:  "Empty conds",
			tbl:   "test",
			cols:  []util.KV{{"name", "'test'"}},
			conds: []util.KV{},
			want:  "UPDATE test SET name = 'test'",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := util.BuildMappedUpdateSQL(tt.tbl, tt.cols, tt.conds)
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
		conds []util.KV
		want  string
	}{
		{
			name:  "Multiple conditions",
			tbl:   "users",
			conds: []util.KV{{"id", "5"}, {"status", "'inactive'"}},
			want:  "DELETE FROM users WHERE id = 5 AND status = 'inactive'",
		},
		{
			name:  "No conditions",
			tbl:   "orders",
			conds: []util.KV{},
			want:  "DELETE FROM orders",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := util.BuildMappedDeleteSQL(tt.tbl, tt.conds)
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
			name: "Sorted columns",
			tbl:  "users",
			cols: []string{"username", "age"},
			want: "INSERT INTO users (age, username) VALUES (:age, :username)",
		},
		{
			name: "Empty columns",
			tbl:  "test",
			cols: []string{},
			want: "INSERT INTO test () VALUES ()",
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

			got := util.BuildNamedInsertSQL(tt.tbl, tt.cols)
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
		conds []string
		want  string
	}{
		{
			name:  "Sorted conditions",
			tbl:   "users",
			conds: []string{"username", "age"},
			want:  "SELECT * FROM users WHERE age = :age AND username = :username",
		},
		{
			name:  "No conditions",
			tbl:   "test",
			conds: []string{},
			want:  "SELECT * FROM test",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := util.BuildNamedQuerySQL(tt.tbl, tt.conds)
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
			name:  "Sorted cols AND conds",
			tbl:   "users",
			cols:  []string{"username", "age"},
			conds: []string{"id", "status"},
			want:  "UPDATE users SET age = :age, username = :username WHERE id = :id AND status = :status",
		},
		{
			name:  "Empty cols",
			tbl:   "test",
			cols:  []string{},
			conds: []string{"id"},
			want:  "UPDATE test SET  WHERE id = :id",
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

			got := util.BuildNamedUpdateSQL(tt.tbl, tt.cols, tt.conds)
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
			name:  "Sorted conditions",
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

			got := util.BuildNamedDeleteSQL(tt.tbl, tt.conds)
			if got != tt.want {
				t.Errorf("Want %q, got %q", tt.want, got)
			}
		})
	}
}
