package main

import(
	"fmt"
	"log"
	_ "github.com/go-sql-driver/mysql"
	"github.com/surdeus/gosrv/src/dbx/sqlx"
	//"github.com/surdeus/gosrv/src/cmd/dbtest/structs"
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
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	/*sqlers := []sqlx.Sqler{
		structs.Test{},
		structs.AnotherTest{},
	}*/

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

	q = sqlx.Q().Insert("NewValue", "Age").
		Into("Tests").
		Values(sqlx.Int(5), sqlx.Int(1337))
	_, _, err = db.Do(q)
	if err != nil {
		panic(err)
	}
}

