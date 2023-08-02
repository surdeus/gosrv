package main

import (
	//"fmt"
	"log"
	"os"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/mojosa-software/gosrv/cmd/dbtest/structs"
	"github.com/mojosa-software/gosrv/src/dbx/sqlx"
)

func main() {
	sqlers := []sqlx.Sqler{
		structs.Test{},
		structs.AnotherTest{},
	}

	portStr := os.Getenv("MYSQL_PORT")
	port, err := strconv.Atoi(portStr)
	if err != nil {
		log.Fatal(err)
	}

	db, err := sqlx.Open(
		sqlx.ConnConfig{
			Driver:   "mysql",
			Login:    os.Getenv("MYSQL_USER"),
			Password: os.Getenv("MYSQL_PASSWORD"),
			Host:     os.Getenv("MYSQL_HOST"),
			Port:     port,
			Name:     os.Getenv("MYSQL_DB"),
		},
		sqlers,
	)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	db.Debug = true
	err = db.Migrate()
	if err != nil {
		log.Fatal(err)
	}

	_, _, err = db.Do(
		sqlx.Q().Insert("DickValue").Into("Tests").Values(sqlx.Int(5)),
	)
	if err != nil {
		log.Fatal(err)
	}

	_, rs, err := db.Do(
		sqlx.Q().Select("Id").From("Tests"),
	)
	if err != nil {
		log.Fatal(err)
	}

	for rs.Next() {
		var id int
		rs.Scan(&id)
		log.Println(id)
	}

	//db.Do(sqlx.Q().Insert("DickValue").Values())

	/*c2 := sqlx.T().Sum(
		sqlx.T().V(sqlx.Int(1377)),
		sqlx.T().V(sqlx.Int(1377)),
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
	log.Println("done")*/
}
