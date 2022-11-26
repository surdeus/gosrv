package main

import(
	"fmt"
	"log"
	_ "github.com/go-sql-driver/mysql"
	"github.com/surdeus/go-srv/src/dbx/sqlx"
	"github.com/surdeus/go-srv/src/cmd/dbtest/structs"
)

func main(){
	db, err := sqlx.Open(
		sqlx.ConnConfig{
			Driver: "mysql",
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

	t, err := db.ParseColumnType("suck(11)")
	fmt.Printf("%v %s\n", t, err)

	fmt.Println(db.TableCreationStringFor(structs.Test{}))
	schemas, err := db.GetTableSchemas()
	if err != nil {
		log.Println(err)
	}

	for _, schema := range schemas {
		for _, f := range schema.Columns {
			sql, err := db.ColumnToSql(f)
			if err != nil {
				log.Println(err)
			} else {
				fmt.Printf("'%s'\n", sql)
				fmt.Printf("%v\n", f)
			}
		}
	}

	err = db.Migrate(
		[]sqlx.Sqler{
			structs.Test{},
			structs.AnotherTest{},
		},
	)
	if err != nil {
	    log.Println(err)
	}

	err = db.CreateTable(structs.Test{})
	if err != nil {
		log.Println(err)
	}

	fmt.Println(db.TableExists("Organizations"))
	fmt.Println(db.TableExists("SurelyDoesNot"))
	fmt.Println(db.ColumnExists("Tests", "DickValue"))
	fmt.Println(db.ColumnExists("Tests", "SurelyDoesNot"))

	return
	ts := structs.Test{}.Sql()
	i, f, err := ts.PrimaryKeyColumn()
	if err != nil {
		log.Println(err)
	} else {
		fmt.Println(i, f)
	}
}

