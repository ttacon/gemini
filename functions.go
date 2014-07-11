package gemini

import (
	"database/sql"
	"fmt"
	"reflect"
)

func (g *Gemini) Get(i interface{}, keys ...interface{}) error {
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
	if !table.HasPrimaryKey() {
		return NoPrimaryKey
	}

	return g.getItFrom(i, keys, table)
}

func (g *Gemini) getItFrom(i interface{}, keys []interface{}, table *TableMap) error {
	// TODO(ttacon)
	primaryKeys := table.PrimaryKey()
	if len(primaryKeys) != len(keys) {
		return fmt.Errorf(
			"to use get, must provide correct number of primary keys (expected %d, got %d",
			len(primaryKeys),
			len(keys),
		)
	}

	queryString := "select "
	for i, field := range table.Fields {
		if i != 0 {
			queryString += ", "
		}
		queryString += field.columnName
	}

	queryString += " from " + table.TableName + " where "

	for i, key := range primaryKeys {
		if i != 0 {
			queryString += " and "
		}
		// TODO(ttacon): right now this doesn't deal with struct tag names nor
		// ensuring the value at key[i] is a decent value (not struct or pointer)
		// also, what if that field is a Struct, we need to know how to set the id
		// correctly
		queryString += fmt.Sprintf("%s = %v", key.Name, keys[i])
	}

	// this currently won't work for MongoDB
	db, ok := g.TableToDatabaseInfo[table.TableName]
	if !ok {
		return fmt.Errorf("no database info for table %q", table.TableName)
	}

	rows, err := db.Db.Query(queryString)
	if err != nil {
		return err
	}

	cols, err := rows.Columns()
	if err != nil {
		return err
	}

	if reflect.TypeOf(i).Kind() == reflect.Slice {

		for rows.Next() {
			if rows.Err() != nil {
				return rows.Err()
			}

			v := reflect.ValueOf(i)
			if v.Kind() == reflect.Ptr {
				v = v.Elem()
			}

			target := make([]interface{}, len(cols))
			for i, col := range cols {
				// TODO(ttacon): go through evern column here
				// TODO(ttacon): need to make sure this is all safe
				f := v.FieldByName(table.ColumnNameToMapping[col].structFieldName)
				target[i] = f.Addr().Interface()
			}
		}
	}
	if !rows.Next() {
		if rows.Err() != nil {
			return rows.Err()
		}
	}

	v := reflect.ValueOf(i)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	target := make([]interface{}, len(cols))
	for i, col := range cols {
		// TODO(ttacon): go through evern column here
		// TODO(ttacon): need to make sure this is all safe
		f := v.FieldByName(table.ColumnNameToMapping[col].structFieldName)
		target[i] = f.Addr().Interface()
	}

	return rows.Scan(target...)
}

func (g *Gemini) Insert(i interface{}) error {
	if reflect.TypeOf(i).Kind() != reflect.Ptr {
		return fmt.Errorf("cannot insert non pointer type")
	}

	e := reflect.ValueOf(i).Elem()
	tableName := tableNameForStruct(e.Type())

	// TODO(ttacon): perhaps we should be smart and try to just insert if there is only one db
	// even if we don't have a mapping from table to db
	dbInfo, ok := g.TableToDatabaseInfo[tableName]
	if !ok {
		return fmt.Errorf("table %s is not specified to interact with any db", tableName)
	}
	db := dbInfo.Db

	tMap, ok := g.StructsMap[tableName]
	if !ok {
		return fmt.Errorf("table %s does not have a table map", tableName)
	}

	if reflect.TypeOf(dbInfo.Dialect) == reflect.TypeOf(MongoDB{}) {
		return dbInfo.MongoSesh.DB(dbInfo.DbName).C(tableName).Insert(i)
	}

	// TODO(ttacon): make smart mapping of table name to db driver and dialect
	query, args := insertQueryAndArgs(e, tMap, dbInfo.Dialect)
	// TODO(ttacon): use result (the underscored place)?
	var autoIncrId int64
	if reflect.TypeOf(dbInfo.Dialect) == reflect.TypeOf(PostgresDialect{}) {
		rows := db.QueryRow(query, args...)
		if tMap.autoIncrField != nil {
			err := rows.Scan(&autoIncrId)
			if err != nil {
				return err
			}
		}
	} else {
		result, err := db.Exec(query, args...)
		if err != nil {
			return err
		}
		if tMap.autoIncrField != nil {
			autoIncrId, err = result.LastInsertId()
			if err != nil {
				return err
			}
		}
	}

	if tMap.autoIncrField != nil {
		fieldVal := e.FieldByName(tMap.autoIncrField.Name)
		k := fieldVal.Kind()

		if (k == reflect.Int) || (k == reflect.Int16) || (k == reflect.Int32) || (k == reflect.Int64) {
			fieldVal.SetInt(autoIncrId)
		} else if (k == reflect.Uint16) || (k == reflect.Uint32) || (k == reflect.Uint64) {
			fieldVal.SetUint(uint64(autoIncrId))
		}
	}

	return nil
}

func (g *Gemini) Delete(i interface{}) error {
	// TODO(ttacon)
	e := reflect.ValueOf(i).Elem()
	tableName := tableNameForStruct(e.Type())

	// TODO(ttacon): perhaps we should be smart and try to just insert if there is only one db
	// even if we don't have a mapping from table to db
	dbInfo, ok := g.TableToDatabaseInfo[tableName]
	if !ok {
		return fmt.Errorf("table %s is not specified to interact with any db", tableName)
	}
	db := dbInfo.Db

	tMap, ok := g.StructsMap[tableName]
	if !ok {
		return fmt.Errorf("table %s does not have a table map", tableName)
	}

	if len(tMap.primaryKeys) == 0 {
		return fmt.Errorf("table %s does not have a primary key registered and "+
			"so cannot delete from it by object", tableName)
	}

	// TODO(ttacon): ensure primary keys for table exist on struct passed in

	if reflect.TypeOf(dbInfo.Dialect) == reflect.TypeOf(MongoDB{}) {
		// TODO(ttacon): do it
		//return dbInfo.MongoSesh.DB(dbInfo.DbName).C(tableName).Insert(i)
		return nil
	}

	// TODO(ttacon): make smart mapping of table name to db driver and dialect
	query, args := deleteQueryAndArgs(e, tMap, dbInfo.Dialect)
	_, err := db.Exec(query, args...)
	return err
}

func (g *Gemini) Update(i interface{}) error {
	// TODO(ttacon)
	return nil
}

func (g *Gemini) Select(i interface{}, query string, args ...interface{}) error {
	val := reflect.ValueOf(i)
	if !val.IsValid() {
		return fmt.Errorf("invalid struct value")
	}

	keyVal := reflect.ValueOf(i)
	if !keyVal.IsValid() {
		return fmt.Errorf("invalid key value")
	}

	table := g.tableFor(i)

	dbi, ok := g.TableToDatabaseInfo[table.TableName]
	if !ok {
		return fmt.Errorf("no database info for table %q", table.TableName)
	}
	rows, err := dbi.Db.Query(query, args...)
	if err != nil {
		return err
	}

	cols, err := rows.Columns()
	if err != nil {
		return err
	}

	if !rows.Next() {
		if rows.Err() != nil {
			return rows.Err()
		}
	}

	v := reflect.ValueOf(i)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	target := make([]interface{}, len(cols))
	for i, col := range cols {
		// TODO(ttacon): go through evern column here
		// TODO(ttacon): need to make sure this is all safe
		f := v.FieldByName(table.ColumnNameToMapping[col].structFieldName)
		target[i] = f.Addr().Interface()
	}

	return rows.Scan(target...)
}

func (g *Gemini) Exec(query string, args ...interface{}) (sql.Result, error) {
	if len(g.Dbs) == 1 {
		return g.Dbs[0].Exec(query, args...)
	}
	return nil, NoDbSpecified
}

func (g *Gemini) ExecWithInfo(query string, info *DbInfo, args ...interface{}) (sql.Result, error) {
	// TODO(ttacon): allow users to attach db name to DbInfo so they don't
	// have to hold onto the db info
	return info.Db.Exec(query, args...)
}

func (g *Gemini) tableFor(i interface{}) *TableMap {
	var (
		tableName string
		val       = reflect.ValueOf(i)
	)

	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

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
