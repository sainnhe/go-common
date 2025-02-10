package util

import (
	"fmt"
	"sort"
	"strings"

	"github.com/teamsorghum/go-common/pkg/log"
)

/*
BuildMappedInsertSQL is used to build mapped insert SQL statement.

Note that this function will use string replace, so make sure the values passed into this function is safe.

Params:
  - tbl string: The table name.
  - cols map[string]string: The column names and values. For example {"username": "$1", "nickname": "'foo'",
    "create_at": "now()"}.
  - logger log.Logger: The logger.

Returns:
  - string: The SQL statement.

Example:

	INSERT INTO mytbl (username, nickname, create_at) VALUES ($1, 'foo', now())
*/
func BuildMappedInsertSQL(tbl string, cols map[string]string, logger log.Logger) string {
	s1 := make([]string, 0, len(cols))
	s2 := make([]string, 0, len(cols))
	keys := make([]string, 0, len(cols))
	for k := range cols {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		s1 = append(s1, k)
		s2 = append(s2, cols[k])
	}
	sql := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)",
		tbl, strings.Join(s1, ", "), strings.Join(s2, ", "))
	if logger != nil {
		logger.Debugf("BuildMappedInsertSQL: %s", sql)
	}
	return sql
}

/*
BuildMappedQuerySQL is used to build mapped query SQL statement.

Note that this function will use string replace, so make sure the values passed into this function is safe.

Params:
  - tbl string: The table name.
  - conds map[string]string: The equal conditions. For example {"username": "$1", "nickname": "'foo'"}.
  - logger log.Logger: The logger.

Returns:
  - string: The SQL statement.

Example:

	SELECT * FROM mytbl WHERE username = $1 AND nickname = 'foo'
*/
func BuildMappedQuerySQL(tbl string, conds map[string]string, logger log.Logger) string {
	s := make([]string, 0, len(conds))
	keys := make([]string, 0, len(conds))
	for k := range conds {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		s = append(s, fmt.Sprintf("%s = %s", k, conds[k]))
	}
	sql := ""
	if len(conds) > 0 {
		sql = fmt.Sprintf("SELECT * FROM %s WHERE %s",
			tbl, strings.Join(s, " AND "))
	} else {
		sql = fmt.Sprintf("SELECT * FROM %s", tbl)
	}
	if logger != nil {
		logger.Debugf("BuildMappedQuerySQL: %s", sql)
	}
	return sql
}

/*
BuildMappedUpdateSQL is used to build mapped update SQL statement.

Note that this function will use string replace, so make sure the values passed into this function is safe.

Params:
  - tbl string: The table name.
  - cols map[string]string: The column names and values. For example {"username": "$1", "nickname": "'foo'",
    "create_at": "now()"}.
  - conds map[string]string: The equal conditions. For example {"username": "$2", "nickname": "'bar'"}.
  - logger log.Logger: The logger.

Returns:
  - string: The SQL statement.

Example:

	UPDATE mytbl SET username = $1, nickname = 'foo' WHERE username = $2 AND nickname = 'bar'
*/
func BuildMappedUpdateSQL(tbl string, cols, conds map[string]string, logger log.Logger) string {
	colSlice := make([]string, 0, len(cols))
	colKeys := make([]string, 0, len(cols))
	for k := range cols {
		colKeys = append(colKeys, k)
	}
	sort.Strings(colKeys)
	for _, k := range colKeys {
		colSlice = append(colSlice, fmt.Sprintf("%s = %s", k, cols[k]))
	}
	condSlice := make([]string, 0, len(conds))
	condKeys := make([]string, 0, len(conds))
	for k := range conds {
		condKeys = append(condKeys, k)
	}
	sort.Strings(condKeys)
	for _, k := range condKeys {
		condSlice = append(condSlice, fmt.Sprintf("%s = %s", k, conds[k]))
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
	if logger != nil {
		logger.Debugf("BuildMappedUpdateSQL: %s", sql)
	}
	return sql
}

/*
BuildMappedDeleteSQL is used to build mapped delete SQL statement.

Note that this function will use string replace, so make sure the values passed into this function is safe.

Params:
  - tbl string: The table name.
  - conds map[string]string: The equal conditions. For example {"username": "$1", "nickname": "'foo'"}.
  - logger log.Logger: The logger.

Returns:
  - string: The SQL statement.

Example:

	DELETE FROM mytbl WHERE username = $1 AND nickname = 'foo'
*/
func BuildMappedDeleteSQL(tbl string, conds map[string]string, logger log.Logger) string {
	s := make([]string, 0, len(conds))
	keys := make([]string, 0, len(conds))
	for k := range conds {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		s = append(s, fmt.Sprintf("%s = %s", k, conds[k]))
	}
	sql := ""
	if len(conds) > 0 {
		sql = fmt.Sprintf("DELETE FROM %s WHERE %s",
			tbl, strings.Join(s, " AND "))
	} else {
		sql = fmt.Sprintf("DELETE FROM %s", tbl)
	}
	if logger != nil {
		logger.Debugf("BuildMappedDeleteSQL: %s", sql)
	}
	return sql
}

/*
BuildNamedInsertSQL is used to build named insert SQL statement.

Note that this function will use string replace, so make sure the values passed into this function is safe.

Params:
  - tbl string: The table name.
  - cols []string: The column names. For example ["username", "nickname"].
  - logger log.Logger: The logger.

Returns:
  - string: The SQL statement.

Example:

	INSERT INTO mytbl (username, nickname) VALUES (:username, :nickname)
*/
func BuildNamedInsertSQL(tbl string, cols []string, logger log.Logger) string {
	s1 := make([]string, 0, len(cols))
	s2 := make([]string, 0, len(cols))
	sort.Strings(cols)
	for _, col := range cols {
		s1 = append(s1, col)
		s2 = append(s2, ":"+col)
	}
	sql := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)",
		tbl, strings.Join(s1, ", "), strings.Join(s2, ", "))
	if logger != nil {
		logger.Debugf("BuildNamedInsertSQL: %s", sql)
	}
	return sql
}

/*
BuildNamedQuerySQL is used to build named query SQL statement.

Note that this function will use string replace, so make sure the values passed into this function is safe.

Params:
  - tbl string: The table name.
  - conds []string: The equal conditions. For example ["username", "nickname"].
  - logger log.Logger: The logger.

Returns:
  - string: The SQL statement.

Example:

	SELECT * FROM mytbl WHERE username = :username AND nickname = :nickname
*/
func BuildNamedQuerySQL(tbl string, conds []string, logger log.Logger) string {
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
	if logger != nil {
		logger.Debugf("BuildNamedQuerySQL: %s", sql)
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
  - logger log.Logger: The logger.

Returns:
  - string: The SQL statement.

Example:

	UPDATE mytbl SET username = :username, nickname = :nickname WHERE age = :age AND gender = :gender
*/
func BuildNamedUpdateSQL(tbl string, cols, conds []string, logger log.Logger) string {
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
	if logger != nil {
		logger.Debugf("BuildNamedUpdateSQL: %s", sql)
	}
	return sql
}

/*
BuildNamedDeleteSQL is used to build named delete SQL statement.

Note that this function will use string replace, so make sure the values passed into this function is safe.

Params:
  - tbl string: The table name.
  - conds []string: The equal conditions. For example ["username", "nickname"].
  - logger log.Logger: The logger.

Returns:
  - string: The SQL statement.

Example:

	DELETE FROM mytbl WHERE username = :username AND nickname = :nickname
*/
func BuildNamedDeleteSQL(tbl string, conds []string, logger log.Logger) string {
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
	if logger != nil {
		logger.Debugf("BuildNamedDeleteSQL: %s", sql)
	}
	return sql
}
