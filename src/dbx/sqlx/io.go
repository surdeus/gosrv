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
	cMap, ok := db.TCMap[tname]
	if !ok {
		return nil, TableDoesNotExistErr
	}

	rc := db.AMap[tname]

	row := make([]any, len(cnames))
	t := reflect.TypeOf(rc)
	val := reflect.New(t)
	val = val.Elem()
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

