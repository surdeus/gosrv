package main

import(
	"fmt"
	"log"
	_ "github.com/go-sql-driver/mysql"
	"github.com/surdeus/ghost/src/db/sqlx"
)

type Test struct {
	Id int
	Value int
	StringValue string
}

func (t Test)Sql() sqlx.TableSchema {
	return sqlx.TableSchema {
		Name: "Tests",
		Fields: []sqlx.TableField {
			{
				Name: "Id",
				Type: "int(11)",
				Nullable: false,
				Key: "PRI",
				Extra: "auto_increment",
			},{
				Name: "Value",
				Type: "int(11)",
				Nullable: true,
				Default: "25",
			},{
				Name: "StringValue",
				Type: "varchar(64)",
				Nullable: true,
				Default: "'some русская string'",
			},

		},
	}
}

type AnotherTest struct {
	Id int
	Value int
	StringValue string
}

func (t AnotherTest)Sql() sqlx.TableSchema {
	return sqlx.TableSchema {
		OldName: "AnotherTests",
		Name: "BetterTests",
		Fields: []sqlx.TableField {
			{
				Name: "Id",
				Type: "int(11)",
				Nullable: false,
				Key: "PRI",
				Extra: "auto_increment",
			},{
				Name: "AnotherValue",
				Type: "int(11)",
				Nullable: true,
				Default: "25",
			},
		},
	}
}

func main(){
	db, err := sqlx.Open(
		"mysql",
		sqlx.ConnConfig{
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
			fmt.Printf("'%s'", db.FieldToSql(f))
			fmt.Println(f)
		}
	}

	fmt.Println(db.TableExists("Organizations"))
	fmt.Println(db.TableExists("SurelyDoesNot"))


	fmt.Println(db.TableCreationStringFor(Test{}))

	err = db.Migrate(
		[]sqlx.Sqler{
			Test{},
			AnotherTest{},
		},
	)
	if err != nil {
	    log.Println(err)
	}
}
