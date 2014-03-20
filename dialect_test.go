package gemini

import "testing"

type insertQueryAndArgsTest struct {
	strct           interface{}
	tbleName        string
	dialect         Dialect
	out             string
	expectedArgsLen int
}

func Test_insertQueryAndArgs(t *testing.T) {
	var tests = []*insertQueryAndArgsTest{
		&insertQueryAndArgsTest{
			TestCreateTableForStruct{},
			"differentName",
			MySQL{},
			"insert into differentName (ID, Name) values (?, ?)",
			2,
		},
	}

	for i, test := range tests {
		query, args := insertQueryAndArgs(
			test.strct,
			TableMapFromStruct(test.strct, test.tbleName),
			test.dialect,
		)

		if len(args) != test.expectedArgsLen {
			t.Errorf("[test %d] expected %d args back, got %d", i, test.expectedArgsLen, len(args))
		}

		if query != test.out {
			t.Errorf("[test %d] expected query %q, got %q", i, test.out, query)
		}
	}
}
