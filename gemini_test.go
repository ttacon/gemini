package gemini

import "testing"

type TestCreateTableForStruct struct {
	TableInfo TableInfo `name:"differentName"`
}

type testCreateQuery struct {
	i       interface{}
	dialect Dialect
	out     string
}

func Test_CreateTableQueryFor(t *testing.T) {
	var tests = []testCreateQuery{
		testCreateQuery{
			struct {
				TableInfo TableInfo `name:"differentName"`
				Text      string
				Id        int64
			}{},
			MySQL{},
			"CREATE TABLE differentName (\n\tText varchar(255),\n\tId bigint\n);",
		},
		testCreateQuery{
			struct {
				TableInfo TableInfo `name:"differentName"`
				Text      string
				Id        int64
			}{},
			SqliteDialect{},
			"CREATE TABLE differentName (\n\tText varchar(255),\n\tId integer\n);",
		},
		testCreateQuery{
			struct {
				TableInfo TableInfo `name:"differentName"`
				Text      string
				Id        int64
			}{},
			PostgresDialect{},
			"CREATE TABLE differentName (\n\tText text,\n\tId bigint\n);",
		},
	}

	for _, test := range tests {
		query := CreateTableQueryFor(test.i, test.dialect)
		if query != test.out {
			t.Errorf("query %q != expected %q", query, test.out)
		}
	}
}
