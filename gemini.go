package gemini

import (
	"database/sql"
	"fmt"
	"reflect"
)

// The sole purpose of TableInfo is to tag extra information you would like.
// For now the main use is to specify a table name different than the struct name -
// this can be done by setting, `name:"tableName"`
type TableInfo struct{}

type geminiMode string

type Gemini struct {
	Dbs map[string]*sql.DB

	StructsMap map[string]TableMap

	// unexported fields
	runInMemory bool

	// map of function name to function pointer
	data map[string][]interface{}

	// next iteration, if there are keys on table, let's make data
	// be a map to a (map) on one, ([]map) some of those keys
}

func NewGemini(dbs []*sql.DB) *Gemini {
	if len(dbs) == 0 {
		return &Gemini{
			runInMemory: true,
			data:        make(map[string][]interface{}),
		}
	}
	return &Gemini{}
}

func (g *Gemini) AddTable(i interface{}) *Gemini {
	g.AddTableWithName(i, reflect.TypeOf(i).Name())
	return g
}

func (g *Gemini) AddTableWithName(i interface{}, tableName string) *Gemini {

	return g
}

func (g *Gemini) CreateTableFor(i interface{}) error {
	// need to know how to pass in which db to interact with, or just type?
	query := CreateTableQueryFor(i, MySQL{})
	fmt.Println(query)
	return nil
}

func CreateTableQueryFor(i interface{}, dialect Dialect) string {
	query := "CREATE TABLE "
	val := reflect.ValueOf(i)
	t := val.Type()
	tableName := val.Type().Name()
	if v, ok := val.Type().FieldByName("TableInfo"); ok {
		if realName := v.Tag.Get("name"); realName != "" {
			tableName = realName
		}
	}
	query += tableName + " (\n"

	n := t.NumField()
	// loop through fields and add to query
	for i := 0; i < n; i++ {
		f := t.Field(i)
		fieldName := f.Name
		if tagName := f.Tag.Get("db"); tagName != "" {
			fieldName = tagName
		}

		// switch to switch or use const for TableInfo{}
		if f.Type == reflect.TypeOf(TableInfo{}) {
			continue
		}
		fmt.Println(fieldName)
		fieldType := f.Type
		// check for db field type in tag

		lineEnding := ","
		if i == n-1 {
			lineEnding = ""
		}

		// query tag for other info like maxsize
		query += fmt.Sprintf(
			"\t%s %s%s\n",
			fieldName,
			dialect.ToSqlType(fieldType, 0, false),
			lineEnding,
		)

	}

	// what about engine, auto inc start charset?
	// put them on tableInfo?
	return query + ");"
}
