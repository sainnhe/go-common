package util

import (
	"fmt"
	"sort"
	"strings"
)

/*
BuildMappedInsertSQL is used to build mapped insert SQL statement.

Note that this function will use string replace, so make sure the values passed into this function is safe.

Params:
  - tbl string: The table name.
  - cols map[string]string: The column names and values. For example {"username": "$1", "nickname": "'foo'",
    "create_at": "now()"}.

Returns:
  - string: The SQL statement.
*/
func BuildMappedInsertSQL(tbl string, cols map[string]string) string {
	s1 := []string{}
	s2 := []string{}
	keys := make([]string, 0, len(cols))
	for k := range cols {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		s1 = append(s1, k)
		s2 = append(s2, cols[k])
	}
	return fmt.Sprintf("insert into %s (%s) values (%s) returning id",
		tbl, strings.Join(s1, ", "), strings.Join(s2, ", "))
}

/*
BuildMappedQuerySQL is used to build mapped query SQL statement.

Note that this function will use string replace, so make sure the values passed into this function is safe.

Params:
  - tbl string: The table name.
  - conds map[string]string: The equal conditions. For example {"username": "$1", "nickname": "'foo'"} will be built
    into "username = $1 and nickname = 'foo'".

Returns:
  - string: The SQL statement.
*/
func BuildMappedQuerySQL(tbl string, conds map[string]string) string {
	s := []string{}
	keys := make([]string, 0, len(conds))
	for k := range conds {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		s = append(s, fmt.Sprintf("%s = %s", k, conds[k]))
	}
	return fmt.Sprintf("select * from %s where %s",
		tbl, strings.Join(s, " and "))
}

/*
BuildMappedUpdateSQL is used to build mapped update SQL statement.

Note that this function will use string replace, so make sure the values passed into this function is safe.

Params:
  - tbl string: The table name.
  - cols map[string]string: The column names and values. For example {"username": "$1", "nickname": "'foo'",
    "create_at": "now()"}.
  - conds map[string]string: The equal conditions. For example {"username": "$1", "nickname": "'foo'"} will be built
    into "username = $1 and nickname = 'foo'".

Returns:
  - string: The SQL statement.
*/
func BuildMappedUpdateSQL(tbl string, cols, conds map[string]string) string {
	colSlice := []string{}
	colKeys := make([]string, 0, len(cols))
	for k := range cols {
		colKeys = append(colKeys, k)
	}
	sort.Strings(colKeys)
	for _, k := range colKeys {
		colSlice = append(colSlice, fmt.Sprintf("%s = %s", k, cols[k]))
	}
	condSlice := []string{}
	condKeys := make([]string, 0, len(conds))
	for k := range conds {
		condKeys = append(condKeys, k)
	}
	sort.Strings(condKeys)
	for _, k := range condKeys {
		condSlice = append(condSlice, fmt.Sprintf("%s = %s", k, conds[k]))
	}
	return fmt.Sprintf("update %s set %s where %s", tbl,
		strings.Join(colSlice, ", "),
		strings.Join(condSlice, " and "))
}

/*
BuildMappedDeleteSQL is used to build mapped delete SQL statement.

Note that this function will use string replace, so make sure the values passed into this function is safe.

Params:
  - tbl string: The table name.
  - conds map[string]string: The equal conditions. For example {"username": "$1", "nickname": "'foo'"} will be built
    into "username = $1 and nickname = 'foo'".

Returns:
  - string: The SQL statement.
*/
func BuildMappedDeleteSQL(tbl string, conds map[string]string) string {
	s := []string{}
	keys := make([]string, 0, len(conds))
	for k := range conds {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		s = append(s, fmt.Sprintf("%s = %s", k, conds[k]))
	}
	return fmt.Sprintf("delete from %s where %s",
		tbl, strings.Join(s, " and "))
}

/*
BuildNamedInsertSQL is used to build named insert SQL statement.

Note that this function will use string replace, so make sure the values passed into this function is safe.

Params:
  - tbl string: The table name.
  - cols []string: The column names.

Returns:
  - string: The SQL statement.
*/
func BuildNamedInsertSQL(tbl string, cols []string) string {
	s1 := []string{}
	s2 := []string{}
	sort.Strings(cols)
	for _, col := range cols {
		s1 = append(s1, col)
		s2 = append(s2, ":"+col)
	}
	return fmt.Sprintf("insert into %s (%s) values (%s) returning id",
		tbl, strings.Join(s1, ", "), strings.Join(s2, ", "))
}

/*
BuildNamedQuerySQL is used to build named query SQL statement.

Note that this function will use string replace, so make sure the values passed into this function is safe.

Params:
  - tbl string: The table name.
  - conds []string: The equal conditions. For example ["username", "nickname"] will be built into
    "username = :username and nickname = :nickname".

Returns:
  - string: The SQL statement.
*/
func BuildNamedQuerySQL(tbl string, conds []string) string {
	s := []string{}
	sort.Strings(conds)
	for _, cond := range conds {
		s = append(s, fmt.Sprintf("%s = :%s", cond, cond))
	}
	return fmt.Sprintf("select * from %s where %s",
		tbl, strings.Join(s, " and "))
}

/*
BuildNamedUpdateSQL is used to build named update SQL statement.

Note that this function will use string replace, so make sure the values passed into this function is safe.

Params:
  - tbl string: The table name.
  - cols []string: The column names.
  - conds []string: The equal conditions. For example ["username", "nickname"] will be built into
    "username = :username and nickname = :nickname".

Returns:
  - string: The SQL statement.
*/
func BuildNamedUpdateSQL(tbl string, cols, conds []string) string {
	colSlice := []string{}
	sort.Strings(cols)
	for _, col := range cols {
		colSlice = append(colSlice, fmt.Sprintf("%s = :%s", col, col))
	}
	condSlice := []string{}
	sort.Strings(conds)
	for _, cond := range conds {
		condSlice = append(condSlice, fmt.Sprintf("%s = :%s", cond, cond))
	}
	return fmt.Sprintf("update %s set %s where %s", tbl,
		strings.Join(colSlice, ", "),
		strings.Join(condSlice, " and "))
}

/*
BuildNamedDeleteSQL is used to build named delete SQL statement.

Note that this function will use string replace, so make sure the values passed into this function is safe.

Params:
  - tbl string: The table name.
  - conds []string: The equal conditions. For example ["username", "nickname"] will be built into
    "username = :username and nickname = :nickname".

Returns:
  - string: The SQL statement.
*/
func BuildNamedDeleteSQL(tbl string, conds []string) string {
	s := []string{}
	sort.Strings(conds)
	for _, cond := range conds {
		s = append(s, fmt.Sprintf("%s = :%s", cond, cond))
	}
	return fmt.Sprintf("delete from %s where %s",
		tbl, strings.Join(s, " and "))
}
