package main

import (
	"fmt"
	//"encoding/gob"
	//"net/http"
	//"github.com/mojosa-software/gosrv/src/dbx/sqlx"
	//"github.com/mojosa-software/gosrv/src/cmd/dbtest/structs"
	//"github.com/mojosa-software/gosrv/src/httpx/apix"
	//"reflect"
	"github.com/mojosa-software/gosrv/src/httpx"
	"log"
)

func main() {
	resp, err := httpx.ApiQuery("http://localhost:8080/api/", 100)
	if err != nil {
		panic(err)
	}
	
	var v int
	fmt.Println("resp:", resp)
	for resp.Scan(&v) {
		fmt.Println("shit", v)
	}
	if resp.Err() != nil {
		log.Fatal(resp.Err())
	}
	
	/*apix.SqlGobRegister()
	q := sqlx.Q().
		Select("DickValue", "StringValue").
		From("Tests").
		Where(
			sqlx.T().Eq().
				C1("DickValue").
				V2(sqlx.Int(1377)),
		)

	v, err := q.SqlRaw(nil)
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println(v)
	}

	_, rs, err := apix.SqlQuery(
		"http://localhost:8080/api/sql/",
		q,
		reflect.TypeOf(structs.Test{}),
	)
	if err != nil {
		panic(err)
	}

	var buf structs.Test
	for v := range rs {
		buf = v.(structs.Test)
		fmt.Println(buf)
	}

	q = sqlx.Q().
		Insert("DickValue").
		Into("Tests").
		Values(sqlx.Int(1377))
	res, _, err := apix.SqlQuery(
		"http://localhost:8080/api/sql/",
		q,
		reflect.TypeOf(structs.Test{}),
	)
	if err != nil {
		panic(err)
	}
	fmt.Println(res)

	fmt.Println(
	sqlx.ValuerToString(sqlx.Null()),
	sqlx.ValuerToString(sqlx.Int(53))) */
}
