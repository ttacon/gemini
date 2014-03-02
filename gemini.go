package gemini

import (
	"database/sql"
	"reflect"
)

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
