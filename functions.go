package gemini

import (
	"fmt"
	"reflect"
)

func (g *Gemini) Get(i interface{}, key interface{}) error {
	var (
		table     *TableMap
		tableName string
	)

	val := reflect.ValueOf(i)
	if !val.IsValid() {
		return fmt.Errorf("invalid struct value")
	}

	keyVal := reflect.ValueOf(i)
	if !keyVal.IsValid() {
		return fmt.Errorf("invalid key value")
	}

	if v, ok := val.Type().FieldByName("TableInfo"); ok && v.Tag.Get("name") != "" {
		tableName = v.Tag.Get("name")
	} else {
		tableName = val.Type().Name()
	}

	// see if struct exists in table map
	if tMap, ok := g.StructsMap[tableName]; ok {
		table = tMap
	} else {
		// if not, well thanks a lot, let's have some fun
		// and see if the struct defines a primary key on it
		table = TableMapFromStruct(i, tableName)
	}

	// NOTE(ttacon): for now we won't support composite primary keys
	// if we have no primary key at this point, let's
	// let them know that Get() was a silly idea
	if !table.HasPrimaryKey() {
		return NoPrimaryKey
	}

	return g.getItFrom(i, key, table)
}

func (g *Gemini) getItFrom(i interface{}, key interface{}, table *TableMap) error {
	// TODO(ttacon)
	return nil
}
