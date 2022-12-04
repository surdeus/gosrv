package main

import (
	"fmt"
	//"encoding/gob"
	//"net/http"
	"github.com/surdeus/gosrv/src/dbx/sqlx"
	"github.com/surdeus/gosrv/src/cmd/dbtest/structs"
	"github.com/surdeus/gosrv/src/httpx/apix"
	"reflect"
)

func main() {
	apix.SqlGobRegister()
	q := sqlx.Q().
		Select("DickValue", "StringValue").
		From("Tests")

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
	sqlx.ValuerToString(sqlx.Int(53)))
}
