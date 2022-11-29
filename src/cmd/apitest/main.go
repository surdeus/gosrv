package main

import (
	"fmt"
	"encoding/gob"
	"net/http"
	"github.com/surdeus/gosrv/src/dbx/sqlx"
	"github.com/surdeus/gosrv/src/cmd/dbtest/structs"
	"bytes"
	"io"
	"database/sql"
)

func main() {
	gob.Register(sql.NullInt32{})
	q := sqlx.Q().
		Select("Id", "StringValue", "DickValue").
		From("Tests").Where("DickValue", sqlx.Gt, sqlx.Int(5))
	bts := bytes.NewBuffer([]byte{})
	enc := gob.NewEncoder(bts)
	err := enc.Encode(q)
	if err != nil {
		fmt.Println(err)
	}
	resp, err := http.Post(
		"http://localhost:8080/api/sql/",
		"application/gob",
		bts)
	if err != nil {
		panic(err)
	}

	buf := structs.Test{}
	dec := gob.NewDecoder(resp.Body)
	for {
		err = dec.Decode(&buf)
		if err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}
		fmt.Println(buf)
	}
}
