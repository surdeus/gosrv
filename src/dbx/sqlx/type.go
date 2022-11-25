package sqlx

import (
	"fmt"
	"errors"
)

var (
	UnknownColumnType = errors.New("unknown column type")
)

func CT() ColumnType {
	return ColumnType{}
}

func (ct ColumnType)Code() (Code, error) {
	ret := ""
	t, ok := MysqlColumnTypeMapString[ct.VarType]
	if !ok {
		return "", UnknownColumnType
	}

	args, err := ct.Args.SqlCodeTuple()
	if err != nil {
		return "", err
	}

	ret = fmt.Sprintf("%s%s", t, args)

	return Code(ret), nil
}

func (ct ColumnType)Varchar(n int) ColumnType {
	ct.VarType = VarcharColumnVarType
	ct.Args = RawValuers{Int(n)}
	return ct
}

func (ct ColumnType)Int(n int) ColumnType {
	ct.VarType = IntColumnVarType
	ct.Args = RawValuers{Int(n)}
	return ct
}

func (ct ColumnType)Nvarchar(n int) ColumnType {
	ct.VarType = NvarcharColumnVarType
	ct.Args = RawValuers{Int(n)}
	return ct
}


