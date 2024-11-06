package sql

import (
	"time"

	mssql "github.com/denisenkom/go-mssqldb"
	"github.com/golang-sql/civil"
)

func applyMSSQLDataType(arg any, column string, schema map[string]any) (any, error) {
	if len(schema) == 0 {
		return arg, nil
	}
	fs, found := schema[column]
	if !found {
		return arg, nil
	}
	fieldSchema := fs.(map[string]any)

	switch fieldSchema["type"].(string) {
	case "VARCHAR":
		arg = mssql.VarChar(arg.(string))
	case "DATETIME":
		datetime := fieldSchema["datetime"].(map[string]any)
		t, err := time.Parse(datetime["format"].(string), arg.(string))
		if err != nil {
			return arg, err
		}
		arg = mssql.DateTime1(t)
	case "DATETIME_OFFSET":
		datetimeOffset := fieldSchema["datetime_offset"].(map[string]any)
		t, err := time.Parse(datetimeOffset["format"].(string), arg.(string))
		if err != nil {
			return arg, err
		}
		arg = mssql.DateTimeOffset(t)
	case "DATE":
		date := fieldSchema["date"].(map[string]any)
		t, err := time.Parse(date["format"].(string), arg.(string))
		if err != nil {
			return arg, err
		}
		arg = civil.DateOf(t)
	}
	return arg, nil
}
