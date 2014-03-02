package gemini

import (
	"reflect"
	"testing"
)

type Struct1 struct {
	field1         interface{} `db:"field1" expected:"field1"`
	field3         interface{} `expected:"field3"`
	fieldWithNoNum interface{} `db:"field4" expected:"field4"`
}

func Test_getColName(t *testing.T) {
	ty := reflect.TypeOf(Struct1{})
	n := ty.NumField()

	for i := 0; i < n; i++ {
		f := ty.Field(i)
		expected := f.Tag.Get("expected")
		got := getColName(f)
		if got != expected {
			t.Errorf(
				"getColName differs from expected, got: %s, expected %s",
				got,
				expected)
		}
	}
}

type JoinStruct struct {
	Join1 Struct1 `joinInfo:"table1,col1"`
	Join2 Struct1
	Join3 Struct1 `joinInfo:"table3,col3,true"`
}

type joinTest struct {
	in          JoinStruct
	out         *joinInfo
	fieldToTest int
}

var joinTests = []*joinTest{
	&joinTest{
		in: JoinStruct{},
		out: &joinInfo{
			tableName:  "table1",
			columnName: "col1",
			eager:      false,
		},
		fieldToTest: 0,
	},
	&joinTest{
		in:          JoinStruct{},
		out:         nil,
		fieldToTest: 1,
	},
	&joinTest{
		in: JoinStruct{},
		out: &joinInfo{
			tableName:  "table3",
			columnName: "col3",
			eager:      true,
		},
		fieldToTest: 2,
	},
}

func Test_parseJoinInfo(t *testing.T) {
	for _, test := range joinTests {
		ty := reflect.TypeOf(test.in)
		j := parseJoinInfo(ty.Field(test.fieldToTest).Tag.Get("joinInfo"))

		if j != nil && test.out != nil && *test.out != *j {
			t.Errorf(
				"incorrectly parsed join info, got: %v, expected: %v",
				j,
				test.out)
		}
	}
}
