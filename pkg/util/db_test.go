package util_test

import (
	"testing"

	"github.com/sainnhe/go-common/pkg/util"
)

func TestBuildMappedInsertSQL(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		tbl      string
		cols     map[string]string
		expected string
	}{
		{
			name: "multiple columns with sorting",
			tbl:  "users",
			cols: map[string]string{
				"username": "$1",
				"age":      "20",
				"email":    "'test@example.com'",
			},
			expected: "insert into users (age, email, username) values (20, 'test@example.com', $1) returning id",
		},
		{
			name:     "empty columns",
			tbl:      "test",
			cols:     map[string]string{},
			expected: "insert into test () values () returning id",
		},
		{
			name: "single column",
			tbl:  "products",
			cols: map[string]string{
				"name": "'product'",
			},
			expected: "insert into products (name) values ('product') returning id",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := util.BuildMappedInsertSQL(tt.tbl, tt.cols)
			if actual != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, actual)
			}
		})
	}
}

func TestBuildMappedQuerySQL(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		tbl      string
		conds    map[string]string
		expected string
	}{
		{
			name:     "multiple conditions sorted",
			tbl:      "users",
			conds:    map[string]string{"username": "$1", "age": "20"},
			expected: "select * from users where age = 20 and username = $1",
		},
		{
			name:     "single condition",
			tbl:      "products",
			conds:    map[string]string{"id": "5"},
			expected: "select * from products where id = 5",
		},
		{
			name:     "no conditions",
			tbl:      "orders",
			conds:    map[string]string{},
			expected: "select * from orders where ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := util.BuildMappedQuerySQL(tt.tbl, tt.conds)
			if actual != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, actual)
			}
		})
	}
}

func TestBuildMappedUpdateSQL(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		tbl      string
		cols     map[string]string
		conds    map[string]string
		expected string
	}{
		{
			name: "multiple cols and conds",
			tbl:  "users",
			cols: map[string]string{
				"username": "$1",
				"age":      "20",
			},
			conds: map[string]string{
				"id":     "5",
				"status": "'active'",
			},
			expected: "update users set age = 20, username = $1 where id = 5 and status = 'active'",
		},
		{
			name:     "empty cols",
			tbl:      "test",
			cols:     map[string]string{},
			conds:    map[string]string{"id": "1"},
			expected: "update test set  where id = 1",
		},
		{
			name:     "empty conds",
			tbl:      "test",
			cols:     map[string]string{"name": "'test'"},
			conds:    map[string]string{},
			expected: "update test set name = 'test' where ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := util.BuildMappedUpdateSQL(tt.tbl, tt.cols, tt.conds)
			if actual != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, actual)
			}
		})
	}
}

func TestBuildMappedDeleteSQL(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		tbl      string
		conds    map[string]string
		expected string
	}{
		{
			name:     "multiple conditions",
			tbl:      "users",
			conds:    map[string]string{"id": "5", "status": "'inactive'"},
			expected: "delete from users where id = 5 and status = 'inactive'",
		},
		{
			name:     "no conditions",
			tbl:      "orders",
			conds:    map[string]string{},
			expected: "delete from orders where ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := util.BuildMappedDeleteSQL(tt.tbl, tt.conds)
			if actual != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, actual)
			}
		})
	}
}

func TestBuildNamedInsertSQL(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		tbl      string
		cols     []string
		expected string
	}{
		{
			name:     "sorted columns",
			tbl:      "users",
			cols:     []string{"username", "age"},
			expected: "insert into users (age, username) values (:age, :username) returning id",
		},
		{
			name:     "empty columns",
			tbl:      "test",
			cols:     []string{},
			expected: "insert into test () values () returning id",
		},
		{
			name:     "single column",
			tbl:      "products",
			cols:     []string{"name"},
			expected: "insert into products (name) values (:name) returning id",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := util.BuildNamedInsertSQL(tt.tbl, tt.cols)
			if actual != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, actual)
			}
		})
	}
}

func TestBuildNamedQuerySQL(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		tbl      string
		conds    []string
		expected string
	}{
		{
			name:     "sorted conditions",
			tbl:      "users",
			conds:    []string{"username", "age"},
			expected: "select * from users where age = :age and username = :username",
		},
		{
			name:     "no conditions",
			tbl:      "test",
			conds:    []string{},
			expected: "select * from test where ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := util.BuildNamedQuerySQL(tt.tbl, tt.conds)
			if actual != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, actual)
			}
		})
	}
}

func TestBuildNamedUpdateSQL(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		tbl      string
		cols     []string
		conds    []string
		expected string
	}{
		{
			name:     "sorted cols and conds",
			tbl:      "users",
			cols:     []string{"username", "age"},
			conds:    []string{"id", "status"},
			expected: "update users set age = :age, username = :username where id = :id and status = :status",
		},
		{
			name:     "empty cols",
			tbl:      "test",
			cols:     []string{},
			conds:    []string{"id"},
			expected: "update test set  where id = :id",
		},
		{
			name:     "empty conds",
			tbl:      "test",
			cols:     []string{"name"},
			conds:    []string{},
			expected: "update test set name = :name where ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := util.BuildNamedUpdateSQL(tt.tbl, tt.cols, tt.conds)
			if actual != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, actual)
			}
		})
	}
}

func TestBuildNamedDeleteSQL(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		tbl      string
		conds    []string
		expected string
	}{
		{
			name:     "sorted conditions",
			tbl:      "users",
			conds:    []string{"id", "status"},
			expected: "delete from users where id = :id and status = :status",
		},
		{
			name:     "no conditions",
			tbl:      "test",
			conds:    []string{},
			expected: "delete from test where ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := util.BuildNamedDeleteSQL(tt.tbl, tt.conds)
			if actual != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, actual)
			}
		})
	}
}
