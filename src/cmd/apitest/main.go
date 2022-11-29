package main

import (
	"fmt"
	"encoding/gob"
	"net/http"
	"github.com/surdeus/gosrv/src/dbx/sqlx"
	"bytes"
)

func main() {
	q := sqlx.Q().
		Select("StringValue", "DickValue").
		From("Tests")
	bts := bytes.NewBuffer([]byte{})
	enc := gob.NewEncoder(bts)
	err := enc.Encode(q)
	if err != nil {
		fmt.Println(err)
	}
	_, err = http.Post(
		"http://localhost:8080/api/sql/",
		"application/gob",
		bts)
	if err != nil {
		fmt.Println(err)
	}
}
