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
				TableInfo     TableInfo `name:"differentName"`
				Text          string
				Id            int64
				Name          string `db:"differentName"`
				Valid         bool
				OtherId       int
				UnsignedId    uint
				StringPointer *string `db:"str_pntr"`
			}{},
			MySQL{},
			"CREATE TABLE differentName (" +
				"\n\tText varchar(255)," +
				"\n\tId bigint," +
				"\n\tdifferentName varchar(255)," +
				"\n\tValid boolean," +
				"\n\tOtherId int," +
				"\n\tUnsignedId int unsigned," +
				"\n\tstr_pntr varchar(255)" +
				"\n);",
		},
		testCreateQuery{
			struct {
				TableInfo     TableInfo `name:"differentName"`
				Text          string
				Id            int64
				Name          string `db:"differentName"`
				Valid         bool
				OtherId       int
				UnsignedId    uint
				StringPointer *string `db:"str_pntr"`
			}{},
			SqliteDialect{},
			"CREATE TABLE differentName (" +
				"\n\tText varchar(255)," +
				"\n\tId integer," +
				"\n\tdifferentName varchar(255)," +
				"\n\tValid integer," +
				"\n\tOtherId integer," +
				"\n\tUnsignedId varchar(255)," +
				"\n\tstr_pntr varchar(255)" +
				"\n);",
		},
		testCreateQuery{
			struct {
				TableInfo     TableInfo `name:"differentName"`
				Text          string
				Id            int64
				Name          string `db:"differentName"`
				Valid         bool
				OtherId       int
				UnsignedId    uint
				StringPointer *string `db:"str_pntr"`
			}{},
			PostgresDialect{},
			"CREATE TABLE differentName (" +
				"\n\tText text," +
				"\n\tId bigint," +
				"\n\tdifferentName text," +
				"\n\tValid boolean," +
				"\n\tOtherId integer," +
				"\n\tUnsignedId text," +
				"\n\tstr_pntr text" +
				"\n);",
		},
	}

	for _, test := range tests {
		query := CreateTableQueryFor(test.i, test.dialect)
		if query != test.out {
			t.Errorf("query %q != expected %q", query, test.out)
		}
	}
}
