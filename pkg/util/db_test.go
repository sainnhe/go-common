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
		cols map[string]string
		want string
	}{
		{
			name: "Multiple columns with sorting",
			tbl:  "users",
			cols: map[string]string{
				"username": "$1",
				"age":      "20",
				"email":    "'test@example.com'",
			},
			want: "insert into users (age, email, username) values (20, 'test@example.com', $1) returning id",
		},
		{
			name: "Empty columns",
			tbl:  "test",
			cols: map[string]string{},
			want: "insert into test () values () returning id",
		},
		{
			name: "Single column",
			tbl:  "products",
			cols: map[string]string{
				"name": "'product'",
			},
			want: "insert into products (name) values ('product') returning id",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := util.BuildMappedInsertSQL(tt.tbl, tt.cols, nil)
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
		conds map[string]string
		want  string
	}{
		{
			name:  "Multiple conditions sorted",
			tbl:   "users",
			conds: map[string]string{"username": "$1", "age": "20"},
			want:  "select * from users where age = 20 and username = $1",
		},
		{
			name:  "Single condition",
			tbl:   "products",
			conds: map[string]string{"id": "5"},
			want:  "select * from products where id = 5",
		},
		{
			name:  "No conditions",
			tbl:   "orders",
			conds: map[string]string{},
			want:  "select * from orders where ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := util.BuildMappedQuerySQL(tt.tbl, tt.conds, nil)
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
		cols  map[string]string
		conds map[string]string
		want  string
	}{
		{
			name: "Multiple cols and conds",
			tbl:  "users",
			cols: map[string]string{
				"username": "$1",
				"age":      "20",
			},
			conds: map[string]string{
				"id":     "5",
				"status": "'active'",
			},
			want: "update users set age = 20, username = $1 where id = 5 and status = 'active'",
		},
		{
			name:  "Empty cols",
			tbl:   "test",
			cols:  map[string]string{},
			conds: map[string]string{"id": "1"},
			want:  "update test set  where id = 1",
		},
		{
			name:  "Empty conds",
			tbl:   "test",
			cols:  map[string]string{"name": "'test'"},
			conds: map[string]string{},
			want:  "update test set name = 'test' where ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := util.BuildMappedUpdateSQL(tt.tbl, tt.cols, tt.conds, nil)
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
		conds map[string]string
		want  string
	}{
		{
			name:  "Multiple conditions",
			tbl:   "users",
			conds: map[string]string{"id": "5", "status": "'inactive'"},
			want:  "delete from users where id = 5 and status = 'inactive'",
		},
		{
			name:  "No conditions",
			tbl:   "orders",
			conds: map[string]string{},
			want:  "delete from orders where ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := util.BuildMappedDeleteSQL(tt.tbl, tt.conds, nil)
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
			want: "insert into users (age, username) values (:age, :username) returning id",
		},
		{
			name: "Empty columns",
			tbl:  "test",
			cols: []string{},
			want: "insert into test () values () returning id",
		},
		{
			name: "Single column",
			tbl:  "products",
			cols: []string{"name"},
			want: "insert into products (name) values (:name) returning id",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := util.BuildNamedInsertSQL(tt.tbl, tt.cols, nil)
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
			want:  "select * from users where age = :age and username = :username",
		},
		{
			name:  "No conditions",
			tbl:   "test",
			conds: []string{},
			want:  "select * from test where ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := util.BuildNamedQuerySQL(tt.tbl, tt.conds, nil)
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
			name:  "Sorted cols and conds",
			tbl:   "users",
			cols:  []string{"username", "age"},
			conds: []string{"id", "status"},
			want:  "update users set age = :age, username = :username where id = :id and status = :status",
		},
		{
			name:  "Empty cols",
			tbl:   "test",
			cols:  []string{},
			conds: []string{"id"},
			want:  "update test set  where id = :id",
		},
		{
			name:  "Empty conds",
			tbl:   "test",
			cols:  []string{"name"},
			conds: []string{},
			want:  "update test set name = :name where ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := util.BuildNamedUpdateSQL(tt.tbl, tt.cols, tt.conds, nil)
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
			want:  "delete from users where id = :id and status = :status",
		},
		{
			name:  "No conditions",
			tbl:   "test",
			conds: []string{},
			want:  "delete from test where ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := util.BuildNamedDeleteSQL(tt.tbl, tt.conds, nil)
			if got != tt.want {
				t.Errorf("Want %q, got %q", tt.want, got)
			}
		})
	}
}
