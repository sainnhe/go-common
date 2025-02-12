package util

import (
	"fmt"
	"sort"
	"strings"
)

// KV is the key value pair.
type KV struct {
	Key   string
	Value string
}

/*
BuildMappedInsertSQL is used to build mapped insert SQL statement.

Note that this function will use string replace, so make sure the values passed into this function is safe.

Params:
  - tbl string: The table name.
  - cols []KV: The column names and values. For example [{"username", "$1"}, {"nickname", "'foo'"}, {"create_at",
    "NOW()"}].

Returns:
  - string: The SQL statement.

Example:

	INSERT INTO mytbl (username, nickname, create_at) VALUES ($1, 'foo', NOW())
*/
func BuildMappedInsertSQL(tbl string, cols []KV) string {
	s1 := make([]string, 0, len(cols))
	s2 := make([]string, 0, len(cols))
	for _, col := range cols {
		s1 = append(s1, col.Key)
		s2 = append(s2, col.Value)
	}
	sql := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)",
		tbl, strings.Join(s1, ", "), strings.Join(s2, ", "))
	return sql
}

/*
BuildMappedQuerySQL is used to build mapped query SQL statement.

Note that this function will use string replace, so make sure the values passed into this function is safe.

Params:
  - tbl string: The table name.
  - conds []KV: The equal conditions. For example [{"username", "$1"}, {"nickname", "'foo'"}].

Returns:
  - string: The SQL statement.

Example:

	SELECT * FROM mytbl WHERE username = $1 AND nickname = 'foo'
*/
func BuildMappedQuerySQL(tbl string, conds []KV) string {
	s := make([]string, 0, len(conds))
	for _, cond := range conds {
		s = append(s, fmt.Sprintf("%s = %s", cond.Key, cond.Value))
	}
	sql := ""
	if len(conds) > 0 {
		sql = fmt.Sprintf("SELECT * FROM %s WHERE %s",
			tbl, strings.Join(s, " AND "))
	} else {
		sql = fmt.Sprintf("SELECT * FROM %s", tbl)
	}
	return sql
}

/*
BuildMappedUpdateSQL is used to build mapped update SQL statement.

Note that this function will use string replace, so make sure the values passed into this function is safe.

Params:
  - tbl string: The table name.
  - cols []KV: The column names and values. For example [{"username", "$1"}, {"nickname", "'foo'"}].
  - conds []KV: The equal conditions. For example [{"username", "$2"}, {"nickname", "'bar'"}].

Returns:
  - string: The SQL statement.

Example:

	UPDATE mytbl SET username = $1, nickname = 'foo' WHERE username = $2 AND nickname = 'bar'
*/
func BuildMappedUpdateSQL(tbl string, cols, conds []KV) string {
	colSlice := make([]string, 0, len(cols))
	for _, col := range cols {
		colSlice = append(colSlice, fmt.Sprintf("%s = %s", col.Key, col.Value))
	}
	condSlice := make([]string, 0, len(conds))
	for _, cond := range conds {
		condSlice = append(condSlice, fmt.Sprintf("%s = %s", cond.Key, cond.Value))
	}
	sql := ""
	if len(conds) > 0 {
		sql = fmt.Sprintf("UPDATE %s SET %s WHERE %s", tbl,
			strings.Join(colSlice, ", "),
			strings.Join(condSlice, " AND "))
	} else {
		sql = fmt.Sprintf("UPDATE %s SET %s", tbl,
			strings.Join(colSlice, ", "))
	}
	return sql
}

/*
BuildMappedDeleteSQL is used to build mapped delete SQL statement.

Note that this function will use string replace, so make sure the values passed into this function is safe.

Params:
  - tbl string: The table name.
  - conds []KV: The equal conditions. For example [{"username", "$1"}, {"nickname", "'foo'"}].

Returns:
  - string: The SQL statement.

Example:

	DELETE FROM mytbl WHERE username = $1 AND nickname = 'foo'
*/
func BuildMappedDeleteSQL(tbl string, conds []KV) string {
	s := make([]string, 0, len(conds))
	for _, cond := range conds {
		s = append(s, fmt.Sprintf("%s = %s", cond.Key, cond.Value))
	}
	sql := ""
	if len(conds) > 0 {
		sql = fmt.Sprintf("DELETE FROM %s WHERE %s",
			tbl, strings.Join(s, " AND "))
	} else {
		sql = fmt.Sprintf("DELETE FROM %s", tbl)
	}
	return sql
}

/*
BuildNamedInsertSQL is used to build named insert SQL statement.

Note that this function will use string replace, so make sure the values passed into this function is safe.

Params:
  - tbl string: The table name.
  - cols []string: The column names. For example ["username", "nickname"].

Returns:
  - string: The SQL statement.

Example:

	INSERT INTO mytbl (username, nickname) VALUES (:username, :nickname)
*/
func BuildNamedInsertSQL(tbl string, cols []string) string {
	s1 := make([]string, 0, len(cols))
	s2 := make([]string, 0, len(cols))
	sort.Strings(cols)
	for _, col := range cols {
		s1 = append(s1, col)
		s2 = append(s2, ":"+col)
	}
	sql := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)",
		tbl, strings.Join(s1, ", "), strings.Join(s2, ", "))
	return sql
}

/*
BuildNamedQuerySQL is used to build named query SQL statement.

Note that this function will use string replace, so make sure the values passed into this function is safe.

Params:
  - tbl string: The table name.
  - conds []string: The equal conditions. For example ["username", "nickname"].

Returns:
  - string: The SQL statement.

Example:

	SELECT * FROM mytbl WHERE username = :username AND nickname = :nickname
*/
func BuildNamedQuerySQL(tbl string, conds []string) string {
	s := make([]string, 0, len(conds))
	sort.Strings(conds)
	for _, cond := range conds {
		s = append(s, fmt.Sprintf("%s = :%s", cond, cond))
	}
	sql := ""
	if len(conds) > 0 {
		sql = fmt.Sprintf("SELECT * FROM %s WHERE %s",
			tbl, strings.Join(s, " AND "))
	} else {
		sql = fmt.Sprintf("SELECT * FROM %s", tbl)
	}
	return sql
}

/*
BuildNamedUpdateSQL is used to build named update SQL statement.

Note that this function will use string replace, so make sure the values passed into this function is safe.

Params:
  - tbl string: The table name.
  - cols []string: The column names. For example ["username", "nickname"].
  - conds []string: The equal conditions. For example ["age", "gender"].

Returns:
  - string: The SQL statement.

Example:

	UPDATE mytbl SET username = :username, nickname = :nickname WHERE age = :age AND gender = :gender
*/
func BuildNamedUpdateSQL(tbl string, cols, conds []string) string {
	colSlice := make([]string, 0, len(cols))
	sort.Strings(cols)
	for _, col := range cols {
		colSlice = append(colSlice, fmt.Sprintf("%s = :%s", col, col))
	}
	condSlice := make([]string, 0, len(conds))
	sort.Strings(conds)
	for _, cond := range conds {
		condSlice = append(condSlice, fmt.Sprintf("%s = :%s", cond, cond))
	}
	sql := ""
	if len(conds) > 0 {
		sql = fmt.Sprintf("UPDATE %s SET %s WHERE %s", tbl,
			strings.Join(colSlice, ", "),
			strings.Join(condSlice, " AND "))
	} else {
		sql = fmt.Sprintf("UPDATE %s SET %s", tbl,
			strings.Join(colSlice, ", "))
	}
	return sql
}

/*
BuildNamedDeleteSQL is used to build named delete SQL statement.

Note that this function will use string replace, so make sure the values passed into this function is safe.

Params:
  - tbl string: The table name.
  - conds []string: The equal conditions. For example ["username", "nickname"].

Returns:
  - string: The SQL statement.

Example:

	DELETE FROM mytbl WHERE username = :username AND nickname = :nickname
*/
func BuildNamedDeleteSQL(tbl string, conds []string) string {
	s := make([]string, 0, len(conds))
	sort.Strings(conds)
	for _, cond := range conds {
		s = append(s, fmt.Sprintf("%s = :%s", cond, cond))
	}
	sql := ""
	if len(conds) > 0 {
		sql = fmt.Sprintf("DELETE FROM %s WHERE %s",
			tbl, strings.Join(s, " AND "))
	} else {
		sql = fmt.Sprintf("DELETE FROM %s", tbl)
	}
	return sql
}
