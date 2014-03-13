package gemini

import (
	"reflect"
	"testing"
)

const (
	STRING_CONST = ""
	INT_CONST    = 0
)

type dbFieldTest struct {
	in  interface{}
	out []dbField
}

type Struct2 struct {
	Id   int    `dbInfo:"primaryKey,autoIncr,notNull"`
	Name string `db:"DbName"`
}

var dbFieldTests = []*dbFieldTest{
	&dbFieldTest{
		in: Struct2{},
		out: []dbField{
			dbField{
				fieldName:    "Id",
				columnName:   "Id",
				joinsTo:      nil,
				goType:       reflect.TypeOf(INT_CONST),
				isPrimaryKey: true,
				isAutoIncr:   true,
				isNotNull:    true,
			},
			dbField{
				fieldName:    "Name",
				columnName:   "DbName",
				joinsTo:      nil,
				goType:       reflect.TypeOf(STRING_CONST),
				isPrimaryKey: false,
				isAutoIncr:   false,
				isNotNull:    false,
			},
		},
	},
}

func Test_getFieldsFor(t *testing.T) {
	tMap := TableMap{}
	for _, test := range dbFieldTests {
		fields := tMap.getFieldsFor(test.in)
		if len(fields) != len(test.out) {
			t.Errorf("expected %d fields, got %d",
				len(test.out),
				len(fields))
		}
		for i := 0; i < len(fields); i++ {
			if fields[i] != test.out[i] {
				t.Errorf("fields differ, expected: %#v, got: %#v",
					test.out[i],
					fields[i])
			}
		}
	}
}

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

type dbTest struct {
	in string
	isPrimaryKey,
	isAutoIncr,
	isNotNull bool
}

var dbTests = []*dbTest{
	&dbTest{
		"",
		false,
		false,
		false,
	},
	&dbTest{
		"primaryKey",
		true,
		false,
		false,
	},
	&dbTest{
		"primaryKey,autoIncr",
		true,
		true,
		false,
	},
	&dbTest{
		"primaryKey,notNull",
		true,
		false,
		true,
	},
	&dbTest{
		"primaryKey,autoIncr,notNull",
		true,
		true,
		true,
	},
	&dbTest{
		"autoIncr,notNull",
		false,
		true,
		true,
	},
	&dbTest{
		"autoIncr",
		false,
		true,
		false,
	},
	&dbTest{
		"notNull",
		false,
		false,
		true,
	},
}

func Test_dbInfo(t *testing.T) {
	for _, test := range dbTests {
		if tagIsPrimaryKey(test.in) != test.isPrimaryKey {
			t.Error(
				"tag was primary key: %v, was expected to be: %v",
				tagIsPrimaryKey(test.in),
				test.isPrimaryKey)
		}

		if tagIsAutoIncr(test.in) != test.isAutoIncr {
			t.Error(
				"tag was auto incr: %v, was expected to be: %v",
				tagIsAutoIncr(test.in),
				test.isAutoIncr)
		}

		if tagIsNotNull(test.in) != test.isNotNull {
			t.Error(
				"tag was not null: %v, was expected to be: %v",
				tagIsNotNull(test.in),
				test.isNotNull)
		}
	}
}
