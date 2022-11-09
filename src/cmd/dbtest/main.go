package main

import(
	"fmt"
	"log"
	_ "github.com/go-sql-driver/mysql"
	"github.com/surdeus/ghost/src/dbs/sqls"
)

func main(){
	db, err := sqls.Open(
		"mysql",
		sqls.ConnConfig{
			Login: "test",
			Password: "hello",
			Host: "localhost",
			Port: 3306,
			Name: "test",
		},
	)
	defer db.Close()
	if err != nil {
		log.Fatal(err)
	}

	schemas, err := db.GetTableSchemas()
	if err != nil {
		log.Println(err)
	}

	for _, schema := range schemas {
		fmt.Println(schema)
	}
}
