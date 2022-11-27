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
	q := sqlx.Q().
		Select("Column", "Column1").
		From("Table").
		Where("Column", sqlx.Gt, sqlx.Int(32)).
		And("Column1", sqlx.Lt, sqlx.Float(1.731))
	s, err := q.SqlRaw(db)
	fmt.Printf("%q, %q, %v\n", s, err, q.GetValues())

	q = sqlx.Q().RenameTable("Table", "NewName")
	s, err = q.SqlRaw(db)
	fmt.Printf("%q, %q, %v\n", s, err, q.GetValues())

	q = sqlx.Q().RenameColumn("Table", "OldName", "NewName")
	s, err = q.SqlRaw(db)
	fmt.Printf("%q, %q, %v\n", s, err, q.GetValues())
}

