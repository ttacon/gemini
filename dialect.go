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
*/

var TestMode = true

type Dialect interface {
	ToSqlType(val reflect.Type, maxsize int) string
}

type MySQL struct{}

func (m MySQL) ToSqlType(val reflect.Type, maxsize int) string {
	switch val.Kind() {
	case reflect.Ptr:
		return m.ToSqlType(val.Elem(), maxsize)
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
