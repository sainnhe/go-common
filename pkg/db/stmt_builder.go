// -------------------------------------------------------------------------------------------
// Copyright (c) Team Sorghum. All rights reserved.
// Licensed under the GPL v3 License. See LICENSE in the project root for license information.
// -------------------------------------------------------------------------------------------

package db

import (
	"fmt"
	"slices"
	"strings"
)

// KV is the key-value pair that can be used in [StmtBuilder].
type KV struct {
	Key string
	Val string
}

const (
	// Placeholder is the placeholder of an argument that can be used in [StmtBuilder].
	Placeholder = "?"

	pgBeginIndex = 1
)

/*
StmtBuilder builds SQL statements.

# SQL injection and placeholders

Since this builder will use string replacement to build SQL statements,
please make sure the values used here, for example the table name, column names and values,
are safe and won't lead to SQL injection.

If you want to build prepared statements, or use SDKs like sqlx to parsing the output statements,
you can use placeholders to allow for binding parameters to the statements.
Binding parameters is safe and won't lead to SQL injection.

The placeholders used in different databases are listed as follows:

  - MySQL: ?
  - PostgreSQL: $N, where N is the 1-based positional argument index.
  - SQLite: ?

If the given [Type] is [PostgreSQL], and the given value is "?",
this builder will automatically converts "?" to "$N" based placeholders.

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
	// BuildMappedInsertSQL builds mapped insert SQL. If the given cols is empty, an empty string will be returned.
	BuildMappedInsertSQL(cols []KV) string

	// BuildMappedQuerySQL builds mapped query SQL. If the given selectedCols is empty, ["*"] will be used.
	BuildMappedQuerySQL(selectedCols []string, conds []KV) string

	// BuildMappedUpdateSQL builds mapped update SQL. If the given cols is empty, an empty string will be returned.
	BuildMappedUpdateSQL(cols, conds []KV) string

	// BuildMappedDeleteSQL builds mapped delete SQL.
	BuildMappedDeleteSQL(conds []KV) string

	// BuildNamedInsertSQL builds named insert SQL. If the given cols is empty, an empty string will be returned.
	BuildNamedInsertSQL(cols []string) string

	// BuildNamedQuerySQL builds named query SQL. If the given selectedCols is empty, ["*"] will be used.
	BuildNamedQuerySQL(selectedCols, conds []string) string

	// BuildNamedUpdateSQL builds named update SQL. If the given cols is empty, an empty string will be returned.
	BuildNamedUpdateSQL(cols, conds []string) string

	// BuildNamedDeleteSQL builds named delete SQL.
	BuildNamedDeleteSQL(conds []string) string
}

type stmtBuilderImpl struct {
	tbl string
	typ Type
}

// NewStmtBuilder initializes a new [StmtBuilder], where tbl is the table name, and typ is the type of the database.
// Nil will be returned if one of the given arguments is invalid.
func NewStmtBuilder(tbl string, typ Type) StmtBuilder {
	if len(tbl) == 0 || typ < MySQL || typ > SQLite {
		return nil
	}
	return &stmtBuilderImpl{
		tbl,
		typ,
	}
}

func (s *stmtBuilderImpl) escapeColNames(colNames []string) {
	for i := range colNames {
		if colNames[i] == "*" {
			continue
		}
		switch s.typ {
		case MySQL:
			colNames[i] = fmt.Sprintf("`%s`", colNames[i])
		case PostgreSQL, SQLite:
			colNames[i] = fmt.Sprintf("%q", colNames[i])
		}
	}
}

func (s *stmtBuilderImpl) convertPlaceholders(beginIndex *int, vals []string) {
	if s.typ == MySQL || s.typ == SQLite {
		return
	}
	for i := range vals {
		if vals[i] == Placeholder {
			vals[i] = fmt.Sprintf("$%d", *beginIndex)
			(*beginIndex)++
		}
	}
}

func (s *stmtBuilderImpl) buildMappedConds(beginIndex *int, conds []KV) string {
	if len(conds) == 0 {
		return ""
	}
	eqs := make([]string, 0, len(conds))
	for _, kv := range conds {
		val := kv.Val
		if val == Placeholder && s.typ == PostgreSQL {
			val = fmt.Sprintf("$%d", *beginIndex)
			(*beginIndex)++
		}
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
		switch s.typ {
		case MySQL:
			eqs = append(eqs, fmt.Sprintf("`%s` = :%s", cond, cond))
		case PostgreSQL, SQLite:
			eqs = append(eqs, fmt.Sprintf("%q = :%s", cond, cond))
		}
	}
	return fmt.Sprintf(" WHERE %s", strings.Join(eqs, " AND "))
}

func (s *stmtBuilderImpl) BuildMappedInsertSQL(cols []KV) string {
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
	beginIndex := pgBeginIndex
	s.convertPlaceholders(&beginIndex, colVals)
	return fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)",
		s.tbl, strings.Join(colNames, ", "), strings.Join(colVals, ", "))
}

func (s *stmtBuilderImpl) BuildMappedQuerySQL(selectedCols []string, conds []KV) string {
	selectedCols = slices.Clone(selectedCols)
	if len(selectedCols) == 0 {
		selectedCols = []string{"*"}
	}
	s.escapeColNames(selectedCols)
	beginIndex := pgBeginIndex
	return fmt.Sprintf("SELECT %s FROM %s%s",
		strings.Join(selectedCols, ", "),
		s.tbl,
		s.buildMappedConds(&beginIndex, conds),
	)
}

func (s *stmtBuilderImpl) BuildMappedUpdateSQL(cols, conds []KV) string {
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
	beginIndex := pgBeginIndex
	s.convertPlaceholders(&beginIndex, colVals)
	colEqs := make([]string, 0, len(cols))
	for i := range cols {
		colEqs = append(colEqs, fmt.Sprintf("%s = %s", colNames[i], colVals[i]))
	}

	return fmt.Sprintf("UPDATE %s SET %s%s",
		s.tbl,
		strings.Join(colEqs, ", "),
		s.buildMappedConds(&beginIndex, conds),
	)
}

func (s *stmtBuilderImpl) BuildMappedDeleteSQL(conds []KV) string {
	beginIndex := pgBeginIndex
	return fmt.Sprintf("DELETE FROM %s%s",
		s.tbl,
		s.buildMappedConds(&beginIndex, conds),
	)
}

func (s *stmtBuilderImpl) BuildNamedInsertSQL(cols []string) string {
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

func (s *stmtBuilderImpl) BuildNamedQuerySQL(selectedCols, conds []string) string {
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

func (s *stmtBuilderImpl) BuildNamedUpdateSQL(cols, conds []string) string {
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

func (s *stmtBuilderImpl) BuildNamedDeleteSQL(conds []string) string {
	return fmt.Sprintf("DELETE FROM %s%s",
		s.tbl,
		s.buildNamedConds(conds),
	)
}
