package gemini

import (
	"database/sql"
	"reflect"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"

	_ "github.com/lib/pq"

	_ "github.com/mattn/go-sqlite3"

	_ "github.com/ziutek/mymysql/godrv"
)

type TestCreateTableForStruct struct {
	TableInfo TableInfo `name:"differentName"`
}

type testCreateQuery struct {
	i       interface{}
	dialect Dialect
	out     string
}

func TestNewGemini(t *testing.T) {
	g := NewGemini(nil)
	if g == nil {
		t.Error("ruh roh")
	}
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
				Bytes         []byte
				S8            int8
				U8            uint8
				S16           int16
				U16           uint16
				U32           uint32
				U64           uint64
				F32           float32
				F64           float64
				NI64          sql.NullInt64
				NF64          sql.NullFloat64
				NB            sql.NullBool
				Tm            time.Time
			}{},
			MySQL{},
			"CREATE TABLE differentName (" +
				"\n\tText varchar(255)," +
				"\n\tId bigint," +
				"\n\tdifferentName varchar(255)," +
				"\n\tValid boolean," +
				"\n\tOtherId int," +
				"\n\tUnsignedId int unsigned," +
				"\n\tstr_pntr varchar(255)," +
				"\n\tBytes mediumblob," +
				"\n\tS8 tinyint," +
				"\n\tU8 tinyint unsigned," +
				"\n\tS16 smallint," +
				"\n\tU16 smallint unsigned," +
				"\n\tU32 int unsigned," +
				"\n\tU64 bigint unsigned," +
				"\n\tF32 double," +
				"\n\tF64 double," +
				"\n\tNI64 bigint," +
				"\n\tNF64 double," +
				"\n\tNB tinyint," +
				"\n\tTm datetime" +
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
				Bytes         []byte
				S8            int8
				U8            uint8
				S16           int16
				U16           uint16
				U32           uint32
				U64           uint64
				F32           float32
				F64           float64
				NI64          sql.NullInt64
				NF64          sql.NullFloat64
				NB            sql.NullBool
				Tm            time.Time
			}{},
			SqliteDialect{},
			"CREATE TABLE differentName (" +
				"\n\tText varchar(255)," +
				"\n\tId integer," +
				"\n\tdifferentName varchar(255)," +
				"\n\tValid integer," +
				"\n\tOtherId integer," +
				"\n\tUnsignedId varchar(255)," +
				"\n\tstr_pntr varchar(255)," +
				"\n\tBytes blob," +
				"\n\tS8 integer," +
				"\n\tU8 integer," +
				"\n\tS16 integer," +
				"\n\tU16 integer," +
				"\n\tU32 integer," +
				"\n\tU64 integer," +
				"\n\tF32 real," +
				"\n\tF64 real," +
				"\n\tNI64 integer," +
				"\n\tNF64 real," +
				"\n\tNB integer," +
				"\n\tTm datetime" +
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
				Bytes         []byte
				S8            int8
				U8            uint8
				S16           int16
				U16           uint16
				U32           uint32
				U64           uint64
				F32           float32
				F64           float64
				NI64          sql.NullInt64
				NF64          sql.NullFloat64
				NB            sql.NullBool
				Tm            time.Time
			}{},
			PostgresDialect{},
			"CREATE TABLE differentName (" +
				"\n\tText text," +
				"\n\tId bigint," +
				"\n\tdifferentName text," +
				"\n\tValid boolean," +
				"\n\tOtherId integer," +
				"\n\tUnsignedId text," +
				"\n\tstr_pntr text," +
				"\n\tBytes bytea," +
				"\n\tS8 integer," +
				"\n\tU8 integer," +
				"\n\tS16 integer," +
				"\n\tU16 integer," +
				"\n\tU32 integer," +
				"\n\tU64 bigint," +
				"\n\tF32 real," +
				"\n\tF64 double precision," +
				"\n\tNI64 bigint," +
				"\n\tNF64 double precision," +
				"\n\tNB boolean," +
				"\n\tTm timestamp with time zone" +
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

func TestCreateTableFor(t *testing.T) {
	g := &Gemini{}
	if g.CreateTableFor(TestCreateTableForStruct{}) != nil {
		t.Error("At current skeleton, CreateTableFor should not return errors, but it did...")
	}
}

type genTest struct {
	in  interface{}
	out string
}

type TableNameStruct1 struct{}
type TableNameStruct2 struct {
	Field1 string
	Field2 int
}
type TableNameStruct3 struct {
	TableInfo TableInfo `name:"yoloStruct"`
	Field1    string
	Field2    int
}

func Test_tableNameForStruct(t *testing.T) {
	var tests = []genTest{
		genTest{
			TableNameStruct1{},
			"TableNameStruct1",
		},
		genTest{
			TableNameStruct2{},
			"TableNameStruct2",
		},
		genTest{
			TableNameStruct3{},
			"yoloStruct",
		},
	}

	for i, test := range tests {
		if tableNameForStruct(reflect.TypeOf(test.in)) != test.out {
			// TODO(ttacon): remove i from loop, and use test.name (once added to struct)
			t.Errorf("test %d failed in Test_tableNameForStruct", i)
		}
	}
}

type addTableTest struct {
	dbs                  []*sql.DB
	structs              []interface{}
	expectedErr          error
	expectedDbForStructs map[reflect.Type]*sql.DB
}

type ATS1 struct{}
type ATS2 struct {
	ID int
}
type ATS3 struct {
	Name string
}
type ATS4 struct {
	TableInfo `name:"yoloFrontend"`
}
type ATS5 struct {
	TableInfo TableInfo `name:"yoloBackend"`
	ID        uint8
	Name      string
}

func TestAddTable(t *testing.T) {
	// TODO(ttacon): move these to helper
	db, err := sql.Open("sqlite3", "/tmp/gorptest.bin")
	if err != nil {
		t.Errorf("failed to connect to sqlite3, err: %v", err)
	}

	var tests = []addTableTest{
		addTableTest{
			structs: []interface{}{
				ATS1{},
			},
			expectedErr: NoDbSpecified,
		},
		addTableTest{
			dbs: []*sql.DB{db},
			structs: []interface{}{
				ATS1{},
				ATS2{},
				ATS3{},
				ATS4{},
				ATS5{},
			},
			expectedDbForStructs: map[reflect.Type]*sql.DB{
				reflect.TypeOf(ATS1{}): db,
				reflect.TypeOf(ATS2{}): db,
				reflect.TypeOf(ATS3{}): db,
				reflect.TypeOf(ATS4{}): db,
				reflect.TypeOf(ATS5{}): db,
			},
		},
	}

	for i, test := range tests {
		g := NewGemini(test.dbs)
		var err error
		for _, str := range test.structs {
			err = g.AddTable(str)
			if err != nil {
				break
			}
		}

		if err != nil {
			if test.expectedErr != nil && test.expectedErr == err {
				continue
			}
			t.Errorf("test %d failed, err: %v", i, err)
			continue
		}

		for _, str := range test.structs {
			dbFound, err := g.dbForStruct(str)
			if err != nil {
				t.Error(err)
				continue
			}

			if dbFound != test.expectedDbForStructs[reflect.TypeOf(str)] {
				t.Errorf("database that %v is tied was not the expected one", str)
			}
		}
	}
}
