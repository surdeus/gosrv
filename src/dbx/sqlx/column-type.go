package sqlx

import (
	"errors"
	"fmt"
	//"database/sql"
)

type ColumnVarType int
type ColumnType struct {
	VarType ColumnVarType
	Args []int
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

	args := ""
	for i, a := range ct.Args {
		args += fmt.Sprintf("%d", a)
		if i != len(ct.Args)-1 {
			args += ","
		}
	}

	ret = fmt.Sprintf("%s(%s)", t, args)

	return Raw(ret), nil
}

func (ct ColumnType)Varchar(n int) ColumnType {
	ct.VarType = VarcharColumnVarType
	ct.Args = []int{n}
	return ct
}

func (ct ColumnType)Int() ColumnType {
	ct.VarType = IntColumnVarType
	ct.Args = []int{11}
	return ct
}

func (ct ColumnType)IntN(n int) ColumnType {
	ct.VarType = IntColumnVarType
	ct.Args = []int{n}
	return ct
}

func (ct ColumnType)Nvarchar(n int) ColumnType {
	ct.VarType = NvarcharColumnVarType
	ct.Args = []int{n}
	return ct
}

func (ct ColumnType)Double() ColumnType {
	ct.VarType = DoubleColumnVarType
	ct.Args = []int{16, 2}
	return ct
}

