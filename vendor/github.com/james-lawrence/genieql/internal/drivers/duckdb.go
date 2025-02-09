package drivers

import (
	"github.com/james-lawrence/genieql"
	"github.com/james-lawrence/genieql/internal/errorsx"
)

const DuckDB = "github.com/marcboeker/go-duckdb"

// implements the duckdb driver https://github.com/marcboeker/go-duckdb
func init() {
	errorsx.MaybePanic(genieql.RegisterDriver(DuckDB, NewDriver(DuckDB, ddb...)))
}

const ddbDefaultEncode = `func() {}`

const ddbDefaultDecode = `func() {
	if err := {{ .From | expr }}.Scan({{.To | autoreference | expr}}); err != nil {
		return err
	}
}`

const ddbDecodeUUID = `func() {
	if {{ .From | expr }}.Valid {
		if uid, err := uuid.FromBytes([]byte({{ .From | expr }}.String)); err != nil {
			return err
		} else {
			{{ .To | autodereference | expr }} = uid.String()
		}
	}
}`

const ddbEncodeUUID = `func() {
	if {{ .From | expr }}.Valid {
		tmp := {{ .Type | expr }}({{ .From | expr }}.String)
		{{ .To | autodereference | expr }} = tmp
	}
}`

var ddb = []genieql.ColumnDefinition{
	{
		DBTypeName: "VARCHAR",
		Type:       "sql.NullString",
		ColumnType: "sql.NullString",
		Native:     stringExprString,
		Decode:     StdlibDecodeString,
		Encode:     StdlibEncodeString,
	},
	{
		DBTypeName: "BOOLEAN",
		Type:       "sql.NullBool",
		ColumnType: "sql.NullBool",
		Native:     boolExprString,
		Decode:     StdlibDecodeBool,
		Encode:     StdlibEncodeBool,
	},
	{
		DBTypeName: "BIGINT",
		Type:       "sql.NullInt64",
		ColumnType: "sql.NullInt64",
		Native:     int64ExprString,
		Decode:     StdlibDecodeInt64,
		Encode:     StdlibEncodeInt64,
	},
	{
		DBTypeName: "INTEGER",
		Type:       "sql.NullInt32",
		ColumnType: "sql.NullInt32",
		Native:     int32ExprString,
		Decode:     StdlibDecodeInt32,
		Encode:     StdlibEncodeInt32,
	},
	{
		DBTypeName: "SMALLINT",
		Type:       "sql.NullInt16",
		ColumnType: "sql.NullInt16",
		Native:     int16ExprString,
		Decode:     StdlibDecodeInt16,
		Encode:     StdlibEncodeInt16,
	},
	{
		DBTypeName: "FLOAT",
		Type:       "sql.NullFloat64",
		ColumnType: "sql.NullFloat64",
		Native:     float64ExprString,
		Decode:     StdlibDecodeFloat64,
		Encode:     StdlibEncodeFloat64,
	},
	{
		DBTypeName: "UUID",
		Type:       "UUID",
		ColumnType: "sql.NullString",
		Native:     stringExprString,
		Decode:     ddbDecodeUUID,
		Encode:     StdlibEncodeString,
	},
	{
		DBTypeName: "TIMESTAMPZ",
		Type:       "sql.NullTime",
		ColumnType: "sql.NullTime",
		Native:     timeExprString,
		Decode:     StdlibDecodeTime,
		Encode:     StdlibEncodeTime,
	},
}
