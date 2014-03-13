package gemini

import (
	"database/sql"
	"testing"
	"time"
)

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
