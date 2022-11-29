package main

import (
	"fmt"
	"encoding/gob"
	"net/http"
	"github.com/surdeus/gosrv/src/dbx/sqlx"
	"github.com/surdeus/gosrv/src/cmd/dbtest/structs"
	"github.com/surdeus/gosrv/src/httpx/apix"
	"bytes"
	"io"
	//"errors"
)

func main() {
	apix.SqlGobRegister()
	q := sqlx.Q().
		Select("DickValue", "StringValue").
		From("Tests").
		Where("DickValue", sqlx.Eq, sqlx.Int(5))

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

	dec := gob.NewDecoder(resp.Body)

	typ := apix.ErrorSqlResponseType
	err = dec.Decode(&typ)
	if err != nil {
		panic(err)
	}

	switch typ {
	case apix.ErrorSqlResponseType :
		var errbuf error
		err = dec.Decode(&errbuf)
		if err != nil {
			panic(err)
		}
		panic(errbuf)
	case apix.RowsSqlResponseType :
		var buf structs.Test
		for {
			err = dec.Decode(&buf)
			if err == io.EOF {
				break
			} else if err != nil {
				panic(err)
			}
			fmt.Println(buf)
		}
	case apix.ResultSqlResponseType :
		var buf sqlx.Result
		err = dec.Decode(&buf)
		if err != nil {
			panic(err)
		}
		fmt.Println(buf)
	}
}
