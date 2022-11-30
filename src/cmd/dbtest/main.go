package main

import(
	"fmt"
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


	/*err = db.Migrate(sqlers)
	if err != nil{
		log.Println(err)
	}*/

	fmt.Println(err)
	q := sqlx.Q().
		Select("Id", "DickValue", "StringValue").
		From("Tests").
		Where("StringValue", sqlx.Eq, sqlx.String("value"))
	_, rs, err := db.Do(q)
	if err != nil {
		panic(err)
	}
	defer rs.Close()
	for rs.Next() {
		var (
			id, dick int
			s string
		)
		rs.Scan(&id, &dick, &s)
		fmt.Println(id, dick, s)
	}

	q = sqlx.Q().DropPrimaryKey("Tests")
	_, _, err = db.Do(q)
	if err != nil {
		panic(err)
	}
}

