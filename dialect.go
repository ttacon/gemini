package gemini

import (
	"fmt"
	"reflect"
)

type DbDialect int

/**
MySQL
Postgres
Sqlite
MongoDb

To add:
- redis
- FoundationDB
- MariaDB (specifically GIS functions?)
*/

var TestMode = true

type Dialect interface {
	ToSqlType(val reflect.Type, maxsize int, isAutoIncr bool) string
	NextPlaceholder(curr int) (int, string)
}

func insertQueryAndArgs(i interface{}, t *TableMap, dialect Dialect) (string, []interface{}) {
	var (
		v         reflect.Value
		query     string
		valString string
		args      []interface{}
	)

	// TODO(ttacon): make reflect.Value{} constant
	if reflect.TypeOf(i) == reflect.TypeOf(reflect.Value{}) {
		v = i.(reflect.Value)
	} else {
		v = reflect.ValueOf(i)
	}

	// TODO(ttacon): change to use bytes.Buffer and WriteString
	currVal := 1
	query += "insert into " + t.TableName + " ("
	returningQuery := ""
	for _, field := range t.Fields {
		// we ignore table names of course
		if field.goType == reflect.TypeOf(TableInfo{}) {
			continue
		}

		if len(valString) != 0 {
			query += ", "
			valString += ", "
		}

		if field.isAutoIncr {
			// TODO(ttacon): really we should be ignoring this
			// unless that column has a value, so we need a test case for this
			//query += field.columnName
			//valString += "null"
			// TODO(ttacon): make this a function/constant
			if reflect.TypeOf(dialect) == reflect.TypeOf(PostgresDialect{}) {
				returningQuery = fmt.Sprintf(" returning %s", field.columnName)
			}
		} else {
			// TODO(ttacon): need to ignore ignored fields (-), omitempty fields and lazy joins?
			query += field.columnName
			// TODO(ttacon): eventually use placeholders here
			nextCurr, placeholder := dialect.NextPlaceholder(currVal)
			valString += placeholder
			currVal = nextCurr
			// valString += "?"
			args = append(args, v.FieldByName(field.structFieldName).Interface())
		}
	}

	return query + ") values (" + valString + ")" + returningQuery, args
}

func deleteQueryAndArgs(i interface{}, t *TableMap, dialect Dialect) (string, []interface{}) {
	var (
		v     reflect.Value
		query string
		args  []interface{}
	)

	// TODO(ttacon): make reflect.Value{} constant
	if reflect.TypeOf(i) == reflect.TypeOf(reflect.Value{}) {
		v = i.(reflect.Value)
	} else {
		v = reflect.ValueOf(i)
	}

	// TODO(ttacon): change to use bytes.Buffer and WriteString
	currVal := 1
	query += fmt.Sprintf("delete from %s where ", t.TableName)
	for i, field := range t.primaryKeys {
		if i != 0 {
			query += " and "
		}

		// TODO(ttacon): eventually use placeholders here
		nextCurr, placeholder := dialect.NextPlaceholder(currVal)
		// TODO(ttacon): need to change primary key slice to
		// deal with field name and column name
		query += fmt.Sprintf("%s = %s", field.Name, placeholder)
		currVal = nextCurr

		args = append(args, v.FieldByName(field.Name).Interface())
	}

	return query, args
}

type MySQL struct{}

func (m MySQL) NextPlaceholder(curr int) (int, string) {
	return curr, "?"
}

func (m MySQL) ToSqlType(val reflect.Type, maxsize int, isAutoIncr bool) string {
	switch val.Kind() {
	case reflect.Ptr:
		return m.ToSqlType(val.Elem(), maxsize, isAutoIncr)
	case reflect.Bool:
		return "boolean"
	case reflect.Int8:
		return "tinyint"
	case reflect.Uint8:
		return "tinyint unsigned"
	case reflect.Int16:
		return "smallint"
	case reflect.Uint16:
		return "smallint unsigned"
	case reflect.Int, reflect.Int32:
		return "int"
	case reflect.Uint, reflect.Uint32:
		return "int unsigned"
	case reflect.Int64:
		return "bigint"
	case reflect.Uint64:
		return "bigint unsigned"
	case reflect.Float64, reflect.Float32:
		return "double"
	case reflect.Slice:
		if val.Elem().Kind() == reflect.Uint8 {
			return "mediumblob"
		}
	}

	switch val.Name() {
	case "NullInt64":
		return "bigint"
	case "NullFloat64":
		return "double"
	case "NullBool":
		return "tinyint"
	case "Time":
		return "datetime"
	}

	if maxsize < 1 {
		maxsize = 255
	}
	return fmt.Sprintf("varchar(%d)", maxsize)
}

type SqliteDialect struct{}

func (s SqliteDialect) NextPlaceholder(curr int) (int, string) {
	return curr, "?"
}

func (d SqliteDialect) ToSqlType(val reflect.Type, maxsize int, isAutoIncr bool) string {
	switch val.Kind() {
	case reflect.Ptr:
		return d.ToSqlType(val.Elem(), maxsize, isAutoIncr)
	case reflect.Bool:
		return "integer"
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return "integer"
	case reflect.Float64, reflect.Float32:
		return "real"
	case reflect.Slice:
		if val.Elem().Kind() == reflect.Uint8 {
			return "blob"
		}
	}

	switch val.Name() {
	case "NullInt64":
		return "integer"
	case "NullFloat64":
		return "real"
	case "NullBool":
		return "integer"
	case "Time":
		return "datetime"
	}

	if maxsize < 1 {
		maxsize = 255
	}
	return fmt.Sprintf("varchar(%d)", maxsize)
}

type PostgresDialect struct{}

func (p PostgresDialect) NextPlaceholder(curr int) (int, string) {
	return curr + 1, fmt.Sprintf("$%d", curr)
}

func (d PostgresDialect) ToSqlType(val reflect.Type, maxsize int, isAutoIncr bool) string {
	switch val.Kind() {
	case reflect.Ptr:
		return d.ToSqlType(val.Elem(), maxsize, isAutoIncr)
	case reflect.Bool:
		return "boolean"
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Uint8, reflect.Uint16, reflect.Uint32:
		if isAutoIncr {
			return "serial"
		}
		return "integer"
	case reflect.Int64, reflect.Uint64:
		if isAutoIncr {
			return "bigserial"
		}
		return "bigint"
	case reflect.Float64:
		return "double precision"
	case reflect.Float32:
		return "real"
	case reflect.Slice:
		if val.Elem().Kind() == reflect.Uint8 {
			return "bytea"
		}
	}

	switch val.Name() {
	case "NullInt64":
		return "bigint"
	case "NullFloat64":
		return "double precision"
	case "NullBool":
		return "boolean"
	case "Time":
		return "timestamp with time zone"
	}

	if maxsize > 0 {
		return fmt.Sprintf("varchar(%d)", maxsize)
	} else {
		return "text"
	}

}

//////////// MongoDB ////////////
type MongoDB struct{}

func (m MongoDB) NextPlaceholder(curr int) (int, string) {
	return curr, "yolo"
}

func (m MongoDB) ToSqlType(val reflect.Type, maxsize int, isAutoIncr bool) string {
	// TODO(ttacon): do it
	return "yolo"
}
