package sqlx

import (
	"reflect"
	"database/sql"
	"log"
)

func (db *Db)ReadRowValues(
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
