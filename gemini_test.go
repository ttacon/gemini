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
			}{},
			MySQL{},
			"CREATE TABLE differentName (\n\tText varchar(255)\n);",
		},
	}

	for _, test := range tests {
		query := CreateTableQueryFor(test.i, test.dialect)
		if query != test.out {
			t.Errorf("query %q != expected %q", query, test.out)
		}
	}
}
