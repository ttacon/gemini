package gemini

import (
	"database/sql"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

type InsertStruct struct {
	Id   int `dbInfo:"autoIncr"`
	Name string
}

func TestInsert_OnlyMysql(t *testing.T) {
	mysqlDb, err := sql.Open(gomysqlConnInfo.Driver, gomysqlConnInfo.DSN)
	if err != nil {
		t.Errorf("failed to connect to gomysql, err: %v", err)
	}

	g := NewGemini([]*sql.DB{mysqlDb})
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
