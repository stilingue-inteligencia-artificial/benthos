package sql

import (
	"time"

	mssql "github.com/denisenkom/go-mssqldb"
	"github.com/golang-sql/civil"
)

func applyMSSQLDataType(arg any, column string, dataTypes map[string]any) (any, error) {
	if len(dataTypes) == 0 {
		return arg, nil
	}
	fdt, found := dataTypes[column]
	if !found {
		return arg, nil
	}
	fieldDataType := fdt.(map[string]any)

	switch fieldDataType["type"].(string) {
	case "VARCHAR":
		arg = mssql.VarChar(arg.(string))
	case "DATETIME":
		datetime := fieldDataType["datetime"].(map[string]any)
		t, err := time.Parse(datetime["format"].(string), arg.(string))
		if err != nil {
			return arg, err
		}
		arg = mssql.DateTime1(t)
	case "DATETIME_OFFSET":
		datetimeOffset := fieldDataType["datetime_offset"].(map[string]any)
		t, err := time.Parse(datetimeOffset["format"].(string), arg.(string))
		if err != nil {
			return arg, err
		}
		arg = mssql.DateTimeOffset(t)
	case "DATE":
		date := fieldDataType["date"].(map[string]any)
		t, err := time.Parse(date["format"].(string), arg.(string))
		if err != nil {
			return arg, err
		}
		arg = civil.DateOf(t)
	}
	return arg, nil
}
