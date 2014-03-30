package gemini

import (
	"database/sql"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"

	"labix.org/v2/mgo"
)

type InsertStruct struct {
	Id   int `dbInfo:"autoIncr,primaryKey"`
	Name string
}

func TestGet_OnlyMysql(t *testing.T) {
	mysqlDb, err := sql.Open(gomysqlConnInfo.Driver, gomysqlConnInfo.DSN)
	if err != nil {
		t.Errorf("failed to connect to gomysql, err: %v", err)
	}

	g := NewGemini([]*DbInfo{
		&DbInfo{
			Dialect: MySQL{},
			Db:      mysqlDb,
			DbName:  "geminitest",
		},
	})

	g.AddTable(InsertStruct{})

	i := InsertStruct{}

	err = g.Get(&i, 1)
	if err != nil {
		t.Error(err)
	}

	if i.Id != 1 {
		t.Errorf("Expected i's id to be 1, got %d (%#v)", i.Id, i)
	}
}

func TestInsert_OnlyMysql(t *testing.T) {
	mysqlDb, err := sql.Open(gomysqlConnInfo.Driver, gomysqlConnInfo.DSN)
	if err != nil {
		t.Errorf("failed to connect to gomysql, err: %v", err)
	}

	g := NewGemini([]*DbInfo{
		&DbInfo{
			Dialect: MySQL{},
			Db:      mysqlDb,
			DbName:  "geminitest",
		},
	})

	i := InsertStruct{
		Name: "yolo",
	}
	g.AddTable(InsertStruct{})
	// TODO(ttacon): create function to set up dbs/tables prior to tests running
	g.CreateTableFor(InsertStruct{}, MySQL{})

	if err = g.Insert(&i); err != nil {
		t.Error(err)
	}
}

func TestInsert_OnlyMongo(t *testing.T) {
	mongodb, err := mgo.DialWithTimeout(mongodbConnInfo.DSN, time.Second)
	if err != nil {
		t.Errorf("failed to connect to gomysql, err: %v", err)
	}

	g := NewGemini([]*DbInfo{
		&DbInfo{
			Dialect:   MongoDB{},
			MongoSesh: mongodb,
			DbName:    "geminitest",
		},
	})

	i := InsertStruct{
		Name: "yolo",
	}
	g.AddTable(InsertStruct{})
	// TODO(ttacon): create function to set up dbs/tables prior to tests running
	//g.CreateTableFor(InsertStruct{}, MongoDB{})

	if err = g.Insert(&i); err != nil {
		t.Error(err)
	}
}

func TestInsert_OnlySqlite(t *testing.T) {
	sqliteDb, err := sql.Open(sqlite3ConnInfo.Driver, sqlite3ConnInfo.DSN)
	if err != nil {
		t.Errorf("failed to connect to sqlite, err: %v", err)
	}

	g := NewGemini([]*DbInfo{
		&DbInfo{
			Dialect: SqliteDialect{},
			Db:      sqliteDb,
			DbName:  "geminitest",
		},
	})

	i := InsertStruct{
		Name: "yolo",
	}
	g.AddTable(InsertStruct{})
	// TODO(ttacon): create function to set up dbs/tables prior to tests running
	g.CreateTableFor(InsertStruct{}, SqliteDialect{})

	if err = g.Insert(&i); err != nil {
		t.Error(err)
	}
}

func TestInsert_OnlyPostgres(t *testing.T) {
	postgresDb, err := sql.Open(postgresConnInfo.Driver, postgresConnInfo.DSN)
	if err != nil {
		t.Errorf("failed to connect to postgres, err: %v", err)
	}

	g := NewGemini([]*DbInfo{
		&DbInfo{
			Dialect: PostgresDialect{},
			Db:      postgresDb,
			DbName:  "geminitest",
		},
	})

	i := InsertStruct{
		Name: "yolo",
	}
	g.AddTable(InsertStruct{})
	// TODO(ttacon): create function to set up dbs/tables prior to tests running
	g.CreateTableFor(InsertStruct{}, PostgresDialect{})

	if err = g.Insert(&i); err != nil {
		t.Error(err)
	}
}
