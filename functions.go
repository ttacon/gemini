package gemini

import (
	"fmt"
	"reflect"
)

func (g *Gemini) Get(i interface{}, key interface{}) error {
	// TODO(ttacon): really the key param should be variadic (for composite primary keys)
	val := reflect.ValueOf(i)
	if !val.IsValid() {
		return fmt.Errorf("invalid struct value")
	}

	keyVal := reflect.ValueOf(i)
	if !keyVal.IsValid() {
		return fmt.Errorf("invalid key value")
	}

	table := g.tableFor(i)

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

func (g *Gemini) Insert(i interface{}) error {
	if reflect.TypeOf(i).Kind() != reflect.Ptr {
		return fmt.Errorf("cannot insert non pointer type")
	}

	e := reflect.ValueOf(i).Elem()
	tableName := tableNameForStruct(e.Type())

	// TODO(ttacon): perhaps we should be smart and try to just insert if there is only one db
	// even if we don't have a mapping from table to db
	db, ok := g.TablesToDb[tableName]
	if !ok {
		return fmt.Errorf("table %s is not specified to interact with any db", tableName)
	}

	tMap, ok := g.StructsMap[tableName]
	if !ok {
		return fmt.Errorf("table %s does not have a table map", tableName)
	}

	if dbInfo, ok := g.TableToDatabaseInfo[tableName]; !ok {
		// TODO(ttacon): we need a better mapping, also, this needs to be changed to
		// be dealt with by the dialect, and not checked like this.
		return fmt.Errorf("no dialect found for db")
	} else {
		if reflect.TypeOf(dbInfo.Dialect) == reflect.TypeOf(MongoDB{}) {
			//TODO(ttacon):todo
			return dbInfo.MongoSesh.DB(dbInfo.DbName).C(tableName).Insert(i)
		}
	}

	// TODO(ttacon): make smart mapping of table name to db driver and dialect
	query, args := insertQueryAndArgs(e, tMap, g.DbToDriver[db])
	// TODO(ttacon): use result (the underscored place)?
	result, err := db.Exec(query, args...)
	if err != nil {
		return err
	}

	if tMap.autoIncrField != nil {
		autoIncrVal, err := result.LastInsertId()
		if err != nil {
			return err
		}

		fieldVal := e.FieldByName(tMap.autoIncrField.Name)
		k := fieldVal.Kind()

		if (k == reflect.Int) || (k == reflect.Int16) || (k == reflect.Int32) || (k == reflect.Int64) {
			fieldVal.SetInt(autoIncrVal)
		} else if (k == reflect.Uint16) || (k == reflect.Uint32) || (k == reflect.Uint64) {
			fieldVal.SetUint(uint64(autoIncrVal))
		}
	}

	return nil
}

func (g *Gemini) Delete(i interface{}) error {
	// TODO(ttacon)
	return nil
}

func (g *Gemini) Update(i interface{}) error {
	// TODO(ttacon)
	return nil
}

func (g *Gemini) Select(i interface{}, query string, args ...interface{}) error {
	// TODO(ttacon)
	return nil
}

func (g *Gemini) Exec(i interface{}, query string, args ...interface{}) error {
	// TODO(ttacon)
	return nil
}

func (g *Gemini) tableFor(i interface{}) *TableMap {
	var (
		tableName string
		val       = reflect.ValueOf(i)
	)

	if v, ok := val.Type().FieldByName("TableInfo"); ok && v.Tag.Get("name") != "" {
		tableName = v.Tag.Get("name")
	} else {
		tableName = val.Type().Name()
	}

	// see if struct exists in table map
	if tMap, ok := g.StructsMap[tableName]; ok {
		return tMap
	}

	return TableMapFromStruct(i, tableName)
}
