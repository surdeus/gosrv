package sqlx

import (
	"errors"
	"fmt"
	"database/sql"
)

type ColumnVarType int
type ColumnType struct {
	VarType ColumnVarType
	Args Valuers
}

const (
	NoColumnVarType = iota
	BoolColumnVaryType
	IntColumnVarType
	SmallintColumnVarType

	BigintColumnVarType

	BitColumnVarType
	TinyintColumnVarType

	DoubleColumnVarType
	FloatColumnVarType

	VarcharColumnVarType
	NvarcharColumnVarType

	CharColumnVarType
	NcharColumnVarType

	TextColumnVarType
	NtextColumnVarType

	DateColumnVarType
	TimeColumnVarType
	TimestampColumnVarType
	DatetimeColumnVarType
	YearColumnVarType

	BinaryColumnVarType
	VarbinaryColumnVarType

	ImageColumnVarType

	ClobColumnVarType
	BlobColumnVarType
	XmlColumnVarType
	JsonColumnVarType
)

var (
	UnknownColumnType = errors.New("unknown column type")
)

func CT() ColumnType {
	return ColumnType{}
}

func (ct ColumnType)SqlRaw(db *Db) (Raw, error) {
	ret := ""
	t, ok := MysqlColumnTypeMapString[ct.VarType]
	if !ok {
		return "", UnknownColumnType
	}

	ret = fmt.Sprintf("%s%s", t, db.TupleBuf(ct.Args))

	return Raw(ret), nil
}

func (ct ColumnType)Varchar(n int32) ColumnType {
	ct.VarType = VarcharColumnVarType
	ct.Args = Valuers{sql.NullInt32{n, true}}
	return ct
}

func (ct ColumnType)Int() ColumnType {
	ct.VarType = IntColumnVarType
	ct.Args = Valuers{sql.NullInt32{11, true}}
	return ct
}

func (ct ColumnType)IntN(n int32) ColumnType {
	ct.VarType = IntColumnVarType
	ct.Args = Valuers{sql.NullInt32{n, true}}
	return ct
}

func (ct ColumnType)Nvarchar(n int32) ColumnType {
	ct.VarType = NvarcharColumnVarType
	ct.Args = Valuers{sql.NullInt32{n, true}}
	return ct
}

func (ct ColumnType)Double() ColumnType {
	ct.VarType = DoubleColumnVarType
	ct.Args = Valuers{
		sql.NullInt32{16, true},
		sql.NullInt32{2, true},
	}
	return ct
}

