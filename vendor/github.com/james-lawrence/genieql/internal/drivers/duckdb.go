package drivers

import (
	"github.com/james-lawrence/genieql"
	"github.com/james-lawrence/genieql/internal/errorsx"
)

const DuckDB = "github.com/marcboeker/go-duckdb"

// implements the duckdb driver https://github.com/marcboeker/go-duckdb
func init() {
	// genieql.DebugColumnDefinitions(ddb...)
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
		Type:       "VARCHAR",
		ColumnType: "sql.NullString",
		Native:     stringExprString,
		Decode:     StdlibDecodeString,
		Encode:     StdlibEncodeString,
	},
	{
		DBTypeName: "BOOLEAN",
		Type:       "BOOLEAN",
		ColumnType: "sql.NullBool",
		Native:     boolExprString,
		Decode:     StdlibDecodeBool,
		Encode:     StdlibEncodeBool,
	},
	{
		DBTypeName: "BIGINT",
		Type:       "BIGINT",
		ColumnType: "sql.NullInt64",
		Native:     int64ExprString,
		Decode:     StdlibDecodeInt64,
		Encode:     StdlibEncodeInt64,
	},
	{
		DBTypeName: "INTEGER",
		Type:       "INTEGER",
		ColumnType: "sql.NullInt32",
		Native:     int32ExprString,
		Decode:     StdlibDecodeInt32,
		Encode:     StdlibEncodeInt32,
	},
	{
		DBTypeName: "UINTEGER",
		Type:       "UINTEGER",
		ColumnType: "sql.NullInt64",
		Native:     uint32ExprString,
		Decode:     StdlibDecodeInt64,
		Encode:     StdlibEncodeInt64,
	},
	{
		DBTypeName: "UBIGINT",
		Type:       "UBIGINT",
		ColumnType: "sql.NullInt64",
		Native:     uint64ExprString,
		Decode:     StdlibDecodeInt64,
		Encode:     StdlibEncodeInt64,
	},
	{
		DBTypeName: "SMALLINT",
		Type:       "SMALLINT",
		ColumnType: "sql.NullInt16",
		Native:     int16ExprString,
		Decode:     StdlibDecodeInt16,
		Encode:     StdlibEncodeInt16,
	},
	{
		DBTypeName: "USMALLINT",
		Type:       "USMALLINT",
		ColumnType: "sql.NullInt32",
		Native:     uint16ExprString,
		Decode:     StdlibDecodeInt32,
		Encode:     StdlibEncodeInt32,
	},
	{
		DBTypeName: "FLOAT",
		Type:       "FLOAT",
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
		Type:       "TIMESTAMPZ",
		ColumnType: "sql.NullTime",
		Native:     timeExprString,
		Decode:     StdlibDecodeTime,
		Encode:     StdlibEncodeTime,
	},
}
