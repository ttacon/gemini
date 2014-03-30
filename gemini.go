package gemini

import (
	"database/sql"
	"fmt"
	"labix.org/v2/mgo"
	"reflect"
	"strings"
)

// The sole purpose of TableInfo is to tag extra information you would like.
// For now the main use is to specify a table name different than the struct name -
// this can be done by setting, `name:"tableName"`
type TableInfo struct{}

type geminiMode string

type DbInfo struct {
	Dialect   Dialect
	Db        *sql.DB
	MongoSesh *mgo.Session
	DbName    string
}

type Gemini struct {
	Dbs []*sql.DB

	StructsMap map[string]*TableMap

	TablesToDb map[string]*sql.DB

	DbToDriver map[*sql.DB]Dialect

	// ugh adding these for the moment, table name -> db info
	TableToDatabaseInfo map[string]*DbInfo

	DatabaseInfo []*DbInfo

	// unexported fields
	runInMemory bool

	// map of function name to function pointer
	data map[string][]interface{}

	// next iteration, if there are keys on table, let's make data
	// be a map to a (map) on one, ([]map) some of those keys
}

func NewGemini(dbsInfo []*DbInfo) *Gemini {
	if len(dbsInfo) == 0 {
		return &Gemini{
			runInMemory:         true,
			StructsMap:          make(map[string]*TableMap),
			TablesToDb:          make(map[string]*sql.DB),
			DbToDriver:          make(map[*sql.DB]Dialect),
			TableToDatabaseInfo: make(map[string]*DbInfo),
			data:                make(map[string][]interface{}),
		}
	}

	dbs := make([]*sql.DB, len(dbsInfo))
	for i, dbInfo := range dbsInfo {
		dbs[i] = dbInfo.Db
	}

	g := &Gemini{
		Dbs:                 dbs,
		StructsMap:          make(map[string]*TableMap),
		TablesToDb:          make(map[string]*sql.DB),
		data:                make(map[string][]interface{}),
		DatabaseInfo:        dbsInfo,
		TableToDatabaseInfo: make(map[string]*DbInfo),
	}

	return g
}

func (g *Gemini) AddTable(i interface{}) error {
	if len(g.DatabaseInfo) != 1 {
		return NoDbSpecified
	}

	g.AddTableWithNameToDb(i, tableNameForStruct(reflect.TypeOf(i)), g.DatabaseInfo[0])
	return nil
}

func (g *Gemini) AddTableWithName(i interface{}, tableName string) error {
	if len(g.DatabaseInfo) != 1 {
		return NoDbSpecified
	}

	g.AddTableWithNameToDb(i, tableName, g.DatabaseInfo[0])
	return nil
}

func (g *Gemini) AddTableToDb(i interface{}, dbInfo *DbInfo) *Gemini {
	g.AddTableWithNameToDb(i, tableNameForStruct(reflect.TypeOf(i)), dbInfo)
	return g
}

func (g *Gemini) AddTableWithNameToDb(
	i interface{},
	tableName string,
	dbInfo *DbInfo) *Gemini {

	g.StructsMap[tableName] = TableMapFromStruct(i, tableName)
	g.TablesToDb[tableName] = dbInfo.Db
	g.TableToDatabaseInfo[tableName] = dbInfo
	return g
}

func (g *Gemini) dbForStruct(i interface{}) (*sql.DB, error) {
	db, ok := g.TablesToDb[tableNameForStruct(reflect.TypeOf(i))]
	if !ok {
		return nil, NoDbForStruct
	}
	return db, nil
}

func (g *Gemini) CreateTableFor(i interface{}, d Dialect) error {
	// need to know how to pass in which db to interact with, or just type?
	tableName := tableNameForStruct(reflect.TypeOf(i))
	db, ok := g.TablesToDb[tableName]
	if !ok {
		return NoDbSpecified
	}

	// TODO(ttacon): the following functions all need to be pulled out into helpers
	query := CreateTableQueryFor(i, d)

	// TODO(ttacon): should have bool about transaction mode
	// TODO(ttacon): don't ignore result
	_, err := db.Exec(query)
	if err != nil {
		return err
	}

	return nil
}

// returns the name of the table the struct should refer to
// if the TableInfo field is not on the struct, we assume
// the table name to be the name of the struct
func tableNameForStruct(t reflect.Type) string {
	tableName := t.Name()
	if v, ok := t.FieldByName("TableInfo"); ok {
		if realName := v.Tag.Get("name"); realName != "" {
			tableName = realName
		}
	}
	return tableName
}

func CreateTableQueryFor(i interface{}, dialect Dialect) string {
	// TODO/NOTE(ttacon): should we be nice and return an error if the struct has no fields?
	query := "CREATE TABLE "
	val := reflect.ValueOf(i)
	t := val.Type()
	query += tableNameForStruct(t) + " (\n"

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
			dialect.ToSqlType(
				fieldType,
				0,
				strings.Contains(f.Tag.Get("dbInfo"), "autoIncr")),
			lineEnding,
		)

	}

	// what about engine, auto inc start charset?
	// put them on tableInfo?
	return query + ");"
}
