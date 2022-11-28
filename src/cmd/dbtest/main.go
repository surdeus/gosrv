package main

import(
	"fmt"
	"log"
	_ "github.com/go-sql-driver/mysql"
	"github.com/surdeus/gosrv/src/dbx/sqlx"
	"github.com/surdeus/gosrv/src/cmd/dbtest/structs"
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

	qs := `create table NewShit (
	        Id int(11) not null primary key auto_increment,
	                DickValue int(11) default ?,
	                        StringValue varchar(32) not null default ?,
	                                NewValue int(11) default ?,
	                                        AnotherValue double(16,2) default ?
	                                        ) ;`

	_, err = db.Query(qs, 60, "some string", 30, 35.5)
	fmt.Println(err)
	return 
	q := sqlx.Q().
		Select("Column", "Column1").
		From("Table").
		Where("Column", sqlx.Gt, sqlx.Int(32)).
		And("Column1", sqlx.Lt, sqlx.Float(1.731)).
		And("Column2", sqlx.In, sqlx.Int(1), sqlx.Int(2))
	s, err := q.SqlRaw(db)
	fmt.Printf("%q, %q, %v\n", s, err, q.GetValues())

	q = sqlx.Q().
		Insert("DickValue", "StringValue").
		Into("Tests").
		Values(sqlx.Int(25), sqlx.String("new"))
	s, err = q.SqlRaw(db)
	vals := []any{}
	for _, v := range q.GetValues() {
		vals = append(vals, any(v))
	}
	_, err = db.Query(string(s), vals...)
	fmt.Printf("%q, %q, %v\n", s, err, vals)

	q = sqlx.Q().RenameTable("Table", "NewName")
	s, err = q.SqlRaw(db)
	fmt.Printf("%q, %q, %v\n", s, err, q.GetValues())

	q = sqlx.Q().RenameColumn("Table", "OldName", "NewName")
	s, err = q.SqlRaw(db)
	fmt.Printf("%q, %q, %v\n", s, err, q.GetValues())

	ts := structs.Test{}.Sql()
	ts.Name = "NewShit"
	q = sqlx.Q().CreateTable(
		ts,
	)
	s, err = q.SqlRaw(db)
	fmt.Println(s, err, q.GetValues())
	_, err = db.Query(string(s), vals...)
	fmt.Println(err)

}

