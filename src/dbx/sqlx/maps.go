package sqlx

import (
	"github.com/surdeus/godat/src/mapx"
)

var (
	MysqlStringMapKeyType = map[string] KeyType {
		"" : NotKeyType,
		"PRI" : PrimaryKeyType,
		"UNI" : UniqueKeyType,
		"MUL" : ForeignKeyType,
	}
	MysqlKeyTypeMapString = mapx.Reverse(
		MysqlStringMapKeyType,
	)
	MysqlColumnTypeMapString = map[ColumnVarType] string {
		IntColumnVarType : "int",
		SmallintColumnVarType : "smallint",

		FloatColumnVarType : "float",
		DoubleColumnVarType : "double",

		BigintColumnVarType : "bigint",
		BitColumnVarType : "bit",
		TinyintColumnVarType : "tinyint",

		VarcharColumnVarType : "varchar",
		NvarcharColumnVarType : "nvarchar",

		CharColumnVarType : "char",
		NcharColumnVarType : "nchar",

		TextColumnVarType : "text",
		NtextColumnVarType : "ntext",

		DateColumnVarType : "date",
		TimeColumnVarType : "time",
		TimestampColumnVarType : "timestamp",
		DatetimeColumnVarType : "datetime",
		YearColumnVarType : "year",

		BinaryColumnVarType : "binary",
		VarbinaryColumnVarType : "varbinary",

		ImageColumnVarType : "image",

		ClobColumnVarType : "clob",
		BlobColumnVarType : "blob",
		XmlColumnVarType : "xml",
		JsonColumnVarType : "json",
	}

	MysqlStringMapColumnType = mapx.Reverse(
		MysqlColumnTypeMapString,
	)

	queryFormatMap = map[QueryType] QueryFormatFunc {
		SelectQueryType : selectQuery,
		InsertQueryType : insertQuery,
		RenameTableQueryType : renameTable,
		RenameColumnQueryType : renameColumn,
		CreateTableQueryType : createTable,
		AlterColumnTypeQueryType : alterColumnType,
		DropPrimaryKeyQueryType : dropPrimaryKey,
	}

	VarTypeMapSqlType = map[ColumnVarType] SqlType {
		BoolColumnVarType : BoolSqlType,
		BitColumnVarType : BoolSqlType,

		IntColumnVarType : Int32SqlType,
		SmallintColumnVarType : Int16SqlType,
		TinyintColumnVarType : ByteSqlType,
		BigintColumnVarType : Int64SqlType,

		DoubleColumnVarType : Float64SqlType,
		FloatColumnVarType : Float64SqlType,

		VarcharColumnVarType : StringSqlType,
		NvarcharColumnVarType : StringSqlType,

		CharColumnVarType : StringSqlType,
		NcharColumnVarType : StringSqlType,

		TextColumnVarType : StringSqlType,
		NtextColumnVarType : StringSqlType,

		DateColumnVarType : TimeSqlType,
		TimeColumnVarType : TimeSqlType,
		TimestampColumnVarType : TimeSqlType,
		DatetimeColumnVarType : TimeSqlType,
		YearColumnVarType : TimeSqlType,

		BinaryColumnVarType : RawBytesSqlType,
		VarbinaryColumnVarType : RawBytesSqlType,

		ImageColumnVarType : RawBytesSqlType,

		ClobColumnVarType : RawBytesSqlType,
		BlobColumnVarType : RawBytesSqlType,
		XmlColumnVarType : RawBytesSqlType,
		JsonColumnVarType : RawBytesSqlType,
	}
)

