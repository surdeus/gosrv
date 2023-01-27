package sqlx

import (
	"reflect"
	"database/sql"
	"log"
)

func (db *Db)ReadRowValues(
	rs *sql.Rows,
	tname TableName,
	cnames ColumnNames,
) (chan any, error) {
	t, ok := db.TMap[tname]
	if !ok {
		return nil, TableDoesNotExistErr
	}

	cMap := t.ColMap

	rcType := db.TypeMap[tname]

	row := make([]any, len(cnames))
	val := reflect.New(rcType).Elem()
	//val = val.Elem()
	for i, v := range cnames {
		f := val.FieldByName(string(v)).Addr()
		_, ok := cMap[v]
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

