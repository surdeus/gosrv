package main

import(
	"fmt"
	"log"
	_ "github.com/go-sql-driver/mysql"
	"github.com/surdeus/ghost/src/dbs/sqls"
)

type Test struct {
	Id int `sql: "int not null primary key"`
	Value int `sql: "int"`
	StringValue string `sql`
}

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
		for _, f := range schema.Fields {
			fmt.Println(f)
			fmt.Println(db.FieldToSQL(f))
		}
	}

	fmt.Println(db.TableExists("Organizations"))
	fmt.Println(db.TableExists("SurelyDoesNot"))
}
