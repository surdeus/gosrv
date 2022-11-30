package main

import (
	"fmt"
	//"encoding/gob"
	//"net/http"
	"github.com/surdeus/gosrv/src/dbx/sqlx"
	"github.com/surdeus/gosrv/src/cmd/dbtest/structs"
	"github.com/surdeus/gosrv/src/httpx/apix"
)

func main() {
	apix.SqlGobRegister(
		structs.Structs,
	)
	/*q := sqlx.Q().
		Select("DickValue", "StringValue").
		From("Tests")
		//Where("DickValue", sqlx.Eq, sqlx.Int(5))

	_, rs, err := apix.SqlQuery(
		"http://localhost:8080/api/sql/",
		q,
		&structs.Test{},
	)
	if err != nil {
		panic(err)
	}

	for v := range rs {
		fmt.Println(v)
	}*/

	q := sqlx.Q().
		Insert("DickValue").
		Into("Tests").
		Values(sqlx.Int(1377))
	res, _, err := apix.SqlQuery(
		"http://localhost:8080/api/sql/",
		q,
		&structs.Test{},
	)
	if err != nil {
		panic(err)
	}
	fmt.Println(res)
}
