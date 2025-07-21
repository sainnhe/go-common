package db

import (
	"fmt"
	"slices"
	"strings"

	"github.com/jmoiron/sqlx"
)

// KV is the key-value pair that can be used in [StmtBuilder].
type KV struct {
	Key string
	Val string
}

// Placeholder is the placeholder of an argument that can be used in [StmtBuilder].
const Placeholder = "?"

/*
StmtBuilder builds SQL statements.

# SQL injection and placeholders

This builder will use string replacement to build SQL statements,
so please make sure the values used here, for example the table name, column names and values,
are safe and won't lead to SQL injection.

If the input argument is at risk of SQL injection, you pass a placeholder and bind values to them instead.

Although different databases use different placeholders (e.g. MySQL and SQLite use ?, and PostgreSQL uses $N, where N is
1-based positional argument index), you should always use ? or [Placeholder] in StmtBuilder.

StmtBuilder will rebind them according to the driver name specified during initialization.

# Mapped and named statements

There are two ways to build SQL statements, namely mapped and named.

Mapped building takes [KV] to build a SQL statement. It maps a value to a key, for example "WHERE `name` = ?".

Mapped building will directly use the value of [KV] in the SQL statement, that means you can also use numbers, strings
and functions as the value. For example:

  - { Key: "age", Val: "20" }
  - { Key: "name", Val: "'my name'" }
  - { Key: "create_at", Val: "NOW()" }

Named building however, only needs the column names. It binds columns based on the "db" struct tag.
For example "WHERE `name` = :name" will bind the value of a struct that has struct field `db:"name"`.

As a rule of thumb, use mapped building when you are targeting at a set of specific columns,
and use named building when you are targeting at a set of specific columns or all columns.
*/
type StmtBuilder interface {
	// GetTbl returns the table name used in this builder.
	GetTbl() string

	// GetDri returns the driver name used in this builder.
	GetDri() string

	// BuildMappedInsertStmt builds mapped insert statement.
	// If the given cols is empty, an empty string will be returned.
	BuildMappedInsertStmt(cols []KV) string

	// BuildMappedQueryStmt builds mapped query statement.
	// If the given selectedCols is empty, ["*"] will be used.
	BuildMappedQueryStmt(selectedCols []string, conds []KV) string

	// BuildMappedUpdateStmt builds mapped update statement.
	// If the given cols is empty, an empty string will be returned.
	BuildMappedUpdateStmt(cols, conds []KV) string

	// BuildMappedDeleteStmt builds mapped delete statement.
	BuildMappedDeleteStmt(conds []KV) string

	// BuildNamedInsertStmt builds named insert statement.
	// If the given cols is empty, an empty string will be returned.
	BuildNamedInsertStmt(cols []string) string

	// BuildNamedQueryStmt builds named query statement.
	// If the given selectedCols is empty, ["*"] will be used.
	BuildNamedQueryStmt(selectedCols, conds []string) string

	// BuildNamedUpdateStmt builds named update statement.
	// If the given cols is empty, an empty string will be returned.
	BuildNamedUpdateStmt(cols, conds []string) string

	// BuildNamedDeleteStmt builds named delete statement.
	BuildNamedDeleteStmt(conds []string) string
}

type stmtBuilderImpl struct {
	tbl string
	dri string
}

// NewStmtBuilder initializes a new [StmtBuilder], where tbl is the table name, and dri is the driver name.
// Nil will be returned if one of the given arguments is invalid.
func NewStmtBuilder(tbl string, dri string) StmtBuilder {
	if len(tbl) == 0 || sqlx.BindType(dri) == sqlx.UNKNOWN {
		return nil
	}
	return &stmtBuilderImpl{
		tbl,
		dri,
	}
}

func (s *stmtBuilderImpl) GetTbl() string {
	return s.tbl
}

func (s *stmtBuilderImpl) GetDri() string {
	return s.dri
}

func (s *stmtBuilderImpl) escapeColNames(colNames []string) {
	for i := range colNames {
		if colNames[i] == "*" {
			continue
		}
		switch s.dri {
		case "mysql":
			colNames[i] = fmt.Sprintf("`%s`", colNames[i])
		case "postgres", "pgx", "sqlite3":
			colNames[i] = fmt.Sprintf("%q", colNames[i])
		}
	}
}

func (s *stmtBuilderImpl) buildMappedConds(conds []KV) string {
	if len(conds) == 0 {
		return ""
	}
	eqs := make([]string, 0, len(conds))
	for _, kv := range conds {
		val := kv.Val
		eqs = append(eqs, fmt.Sprintf("%s = %s", kv.Key, val))
	}
	return fmt.Sprintf(" WHERE %s", strings.Join(eqs, " AND "))
}

func (s *stmtBuilderImpl) buildNamedConds(conds []string) string {
	if len(conds) == 0 {
		return ""
	}
	eqs := make([]string, 0, len(conds))
	for _, cond := range conds {
		switch s.dri {
		case "mysql":
			eqs = append(eqs, fmt.Sprintf("`%s` = :%s", cond, cond))
		case "postgres", "pgx", "sqlite3":
			eqs = append(eqs, fmt.Sprintf("%q = :%s", cond, cond))
		}
	}
	return fmt.Sprintf(" WHERE %s", strings.Join(eqs, " AND "))
}

func (s *stmtBuilderImpl) BuildMappedInsertStmt(cols []KV) string {
	if len(cols) == 0 {
		return ""
	}
	colNames := make([]string, 0, len(cols))
	colVals := make([]string, 0, len(cols))
	for _, col := range cols {
		colNames = append(colNames, col.Key)
		colVals = append(colVals, col.Val)
	}
	s.escapeColNames(colNames)
	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)",
		s.tbl, strings.Join(colNames, ", "), strings.Join(colVals, ", "))
	return sqlx.Rebind(sqlx.BindType(s.dri), query)
}

func (s *stmtBuilderImpl) BuildMappedQueryStmt(selectedCols []string, conds []KV) string {
	selectedCols = slices.Clone(selectedCols)
	if len(selectedCols) == 0 {
		selectedCols = []string{"*"}
	}
	s.escapeColNames(selectedCols)
	query := fmt.Sprintf("SELECT %s FROM %s%s",
		strings.Join(selectedCols, ", "),
		s.tbl,
		s.buildMappedConds(conds),
	)
	return sqlx.Rebind(sqlx.BindType(s.dri), query)
}

func (s *stmtBuilderImpl) BuildMappedUpdateStmt(cols, conds []KV) string {
	if len(cols) == 0 {
		return ""
	}

	// Build columns
	colNames := make([]string, 0, len(cols))
	colVals := make([]string, 0, len(cols))
	for _, col := range cols {
		colNames = append(colNames, col.Key)
		colVals = append(colVals, col.Val)
	}
	s.escapeColNames(colNames)
	colEqs := make([]string, 0, len(cols))
	for i := range cols {
		colEqs = append(colEqs, fmt.Sprintf("%s = %s", colNames[i], colVals[i]))
	}

	query := fmt.Sprintf("UPDATE %s SET %s%s",
		s.tbl,
		strings.Join(colEqs, ", "),
		s.buildMappedConds(conds),
	)
	return sqlx.Rebind(sqlx.BindType(s.dri), query)
}

func (s *stmtBuilderImpl) BuildMappedDeleteStmt(conds []KV) string {
	query := fmt.Sprintf("DELETE FROM %s%s",
		s.tbl,
		s.buildMappedConds(conds),
	)
	return sqlx.Rebind(sqlx.BindType(s.dri), query)
}

func (s *stmtBuilderImpl) BuildNamedInsertStmt(cols []string) string {
	if len(cols) == 0 {
		return ""
	}
	colNames := make([]string, 0, len(cols))
	colVals := make([]string, 0, len(cols))
	for _, col := range cols {
		colNames = append(colNames, col)
		colVals = append(colVals, fmt.Sprintf(":%s", col))
	}
	s.escapeColNames(colNames)
	return fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)",
		s.tbl, strings.Join(colNames, ", "), strings.Join(colVals, ", "))
}

func (s *stmtBuilderImpl) BuildNamedQueryStmt(selectedCols, conds []string) string {
	selectedCols = slices.Clone(selectedCols)
	if len(selectedCols) == 0 {
		selectedCols = []string{"*"}
	}
	s.escapeColNames(selectedCols)
	return fmt.Sprintf("SELECT %s FROM %s%s",
		strings.Join(selectedCols, ", "),
		s.tbl,
		s.buildNamedConds(conds),
	)
}

func (s *stmtBuilderImpl) BuildNamedUpdateStmt(cols, conds []string) string {
	if len(cols) == 0 {
		return ""
	}

	// Build columns
	colNames := make([]string, 0, len(cols))
	colVals := make([]string, 0, len(cols))
	for _, col := range cols {
		colNames = append(colNames, col)
		colVals = append(colVals, fmt.Sprintf(":%s", col))
	}
	s.escapeColNames(colNames)
	colEqs := make([]string, 0, len(cols))
	for i := range cols {
		colEqs = append(colEqs, fmt.Sprintf("%s = %s", colNames[i], colVals[i]))
	}

	return fmt.Sprintf("UPDATE %s SET %s%s",
		s.tbl,
		strings.Join(colEqs, ", "),
		s.buildNamedConds(conds),
	)
}

func (s *stmtBuilderImpl) BuildNamedDeleteStmt(conds []string) string {
	return fmt.Sprintf("DELETE FROM %s%s",
		s.tbl,
		s.buildNamedConds(conds),
	)
}
