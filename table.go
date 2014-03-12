package gemini

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type TableMap struct {
	TableName  string
	StructName string
	Fields     []ColumnMapping
}

type ColumnMapping struct {
	structFieldName string
	columnName      string

	goType       reflect.Type
	isPrimaryKey bool
	isAutoIncr   bool
	isNotNull    bool
}

func TableMapFromStruct(i interface{}, tableName string) *TableMap {
	return &TableMap{
		TableName:  tableName,
		StructName: reflect.TypeOf(i).Name(),
		Fields:     getColumnsFor(i),
	}
}

func getColumnsFor(i interface{}) (cols []ColumnMapping) {
	for _, field := range getFieldsFor(i) {
		cols = append(cols,
			ColumnMapping{
				structFieldName: field.fieldName,
				columnName:      field.columnName,
				goType:          field.goType,
				isPrimaryKey:    field.isPrimaryKey,
				isAutoIncr:      field.isAutoIncr,
				isNotNull:       field.isNotNull,
			})
	}
	return
}

func (t *TableMap) HasPrimaryKey() bool {
	// TODO(ttacon)
	return false
}

func (t *TableMap) PrimaryKey() reflect.Value {
	// TODO(ttacon)
	return reflect.ValueOf(nil)
}

type dbField struct {
	fieldName  string
	columnName string

	joinsTo *joinInfo

	goType reflect.Type

	isPrimaryKey,
	isAutoIncr,
	isNotNull bool
}

type joinInfo struct {
	tableName  string
	columnName string
	eager      bool
	// how to know what field to fill in?
}

func getFieldsFor(i interface{}) []dbField {
	fields := make(map[string]dbField)

	t := reflect.TypeOf(i)
	n := t.NumField()

	for i := 0; i < n; i++ {
		f := t.Field(i)
		field := dbField{
			fieldName:    f.Name,
			columnName:   getColName(f),
			joinsTo:      parseJoinInfo(f.Tag.Get("joinInfo")),
			goType:       f.Type,
			isPrimaryKey: tagIsPrimaryKey(f.Tag.Get("dbInfo")),
			isAutoIncr:   tagIsAutoIncr(f.Tag.Get("dbInfo")),
			isNotNull:    tagIsNotNull(f.Tag.Get("dbInfo")),
		}
		fields[field.fieldName] = field
	}

	var toReturn []dbField
	// TODO(ttacon): change to use length of map and iteration positions
	for _, v := range fields {
		toReturn = append(toReturn, v)
	}
	return toReturn
}

func getColName(col reflect.StructField) string {
	if colName := col.Tag.Get("db"); colName != "" {
		return colName
	}
	return col.Name
}

func parseJoinInfo(j string) *joinInfo {
	var (
		table  string
		column string
		eager  = false
	)

	if len(j) == 0 {
		return nil
	}

	pieces := strings.Split(j, ",")
	switch len(pieces) {
	case 3:
		eager, _ = strconv.ParseBool(pieces[2])
		fallthrough
	case 2:
		column = pieces[1]
		table = pieces[0]
	default:
		panic(JoinInfoError{
			fmt.Sprintf("invalid number of arguments found: %d", len(pieces)),
		})
	}
	return &joinInfo{
		tableName:  table,
		columnName: column,
		eager:      eager,
	}
}

type JoinInfoError struct {
	reason string
}

func (j JoinInfoError) Error() string { return j.reason }

func tagIsPrimaryKey(dbInfo string) bool {
	return strings.Contains(dbInfo, "primaryKey")
}

func tagIsAutoIncr(dbInfo string) bool {
	return strings.Contains(dbInfo, "autoIncr")
}

func tagIsNotNull(dbInfo string) bool {
	return strings.Contains(dbInfo, "notNull")
}
