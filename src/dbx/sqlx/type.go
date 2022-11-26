package sqlx

import (
	"fmt"
	"errors"
	"database/sql"
	"reflect"
	"log"
)

type RowValues = map[ColumnName] any

var (
	UnknownColumnType = errors.New("unknown column type")
)

func C() *Column {
	return &Column{}
}

func (c *Column) WOldName(n ColumnName) *Column {
	c.OldName = n
	return c
}

func (c *Column) WName(n ColumnName) *Column {
	c.Name = n
	return c
}

func (c *Column) WNullable() *Column {
	c.Nullable = true
	return c
}

func (c *Column) WType(t ColumnType) *Column {
	c.Type = t
	return c
}

func (c *Column) WDefault(d RawValuer) *Column {
	c.Default = d
	return c
}

func (c *Column) WKey(k Key) *Column {
	c.Key = k
	return c
}

func (c *Column) WExtra(e Code) *Column {
	c.Extra = e
	return c
}

func CT() ColumnType {
	return ColumnType{}
}

func (ct ColumnType)SqlCode(db *DB) (Code, error) {
	ret := ""
	t, ok := MysqlColumnTypeMapString[ct.VarType]
	if !ok {
		return "", UnknownColumnType
	}

	args, err := ct.Args.SqlCodeTuple(db)
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

func (ct ColumnType)Int() ColumnType {
	ct.VarType = IntColumnVarType
	ct.Args = RawValuers{Int(11)}
	return ct
}

func (ct ColumnType)IntN(n int) ColumnType {
	ct.VarType = IntColumnVarType
	ct.Args = RawValuers{Int(n)}
	return ct
}

func (ct ColumnType)Nvarchar(n int) ColumnType {
	ct.VarType = NvarcharColumnVarType
	ct.Args = RawValuers{Int(n)}
	return ct
}

func (db *DB)ReadRowValues(
	rs *sql.Rows,
	ts *TableSchema,
	cnames ColumnNames,
	tsMap map[ColumnName] *Column,
	rc any,
) (chan any, error) {
	row := make([]any, len(cnames))
	t := reflect.TypeOf(rc)
	val := reflect.New(t)
	val = val.Elem()
	for i, v := range cnames {
		f := val.FieldByName(string(v)).Addr()
		_, ok := tsMap[v]
		if !ok  || !f.IsValid() {
			return nil, ColumnDoesNotExistErr
		}
		row[i] = f.Interface()
	}

	ret := make(chan any)
	go func(){
		for rs.Next() {
			err := rs.Scan(row...)
			if err != nil{
				log.Println(err)
				return
			}
			ret <- val.Interface()
		}
		close(ret)
	}()

	return ret, nil
}

