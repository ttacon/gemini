package gemini

import "testing"

type TestCreateTableForStruct struct {
	TableInfo TableInfo `name:"differentName"`
}

type testCreateQuery struct {
	i   interface{}
	out string
}

func Test_CreateTableQueryFor(t *testing.T) {
	var tests = []testCreateQuery{
		testCreateQuery{
			struct {
				TableInfo TableInfo `name:"differentName"`
			}{},
			"CREATE TABLE differentName (\n);",
		},
	}
	g := &Gemini{}

	for _, test := range tests {
		query := g.CreateTableQueryFor(test.i)
		if query != test.out {
			t.Errorf("query %q != expected %q", query, test.out)
		}
	}
}
