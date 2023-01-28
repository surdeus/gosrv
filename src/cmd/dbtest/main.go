package main

import(
	//"fmt"
	"log"
	_ "github.com/go-sql-driver/mysql"
	"github.com/surdeus/gosrv/src/dbx/sqlx"
	"github.com/surdeus/gosrv/src/cmd/dbtest/structs"
)

func main(){
	sqlers := []sqlx.Sqler{
		structs.Test{},
		structs.AnotherTest{},
	}
	db, err := sqlx.Open(
		sqlx.ConnConfig{
			Driver: "mysql",
			Login: "test",
			Password: "hello",
			Host: "localhost",
			Port: 3306,
			Name: "test",
		},
		sqlers,
	)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	c2 := sqlx.C().And(
		sqlx.C().Eq().
			V1(sqlx.Int(1377)).
			C2("DickValue"),
	)

	log.Printf("%v\n", c2)
	r, err := c2.SqlRaw(db)
	if err != nil {
		log.Println(err)
	} else {
		log.Printf("%q\n", r)
	}

	q := sqlx.Q().
		Select("DickValue", "Id").
		From("Tests").
		Where(c2)

	_, rs, err := db.Do(q)
	if err != nil {
		log.Fatal(err)
	}

	for rs.Next() {
		var dick, id int
		rs.Scan(&dick, &id)
		log.Println(dick, id)
	}
	log.Println("done")
}

