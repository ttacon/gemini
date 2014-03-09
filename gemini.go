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
	query := g.CreateTableQueryFor(i)
	fmt.Println(query)
	return nil
}

func (g *Gemini) CreateTableQueryFor(i interface{}) string {
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
		// remove the Gemini reciever and pass in dialect as a parameter
		// so we can query it for the type
		fmt.Println(fieldName)
		//fieldType := gf.Type

	}

	// what about engine, auto inc start charset?
	// put them on tableInfo?
	return query + ");"
}
