package structs

import (
	"github.com/mojosa-software/gosrv/src/dbx/sqlx"
	//"reflect"
	"errors"
	"fmt"
)

type NewTest struct {
	Id             int
	NullableString *string
	NotNullableInt int
}

type Test struct {
	Id           sqlx.NullInt32
	DickValue    sqlx.NullInt32
	StringValue  sqlx.NullString
	NewValue     sqlx.NullInt32
	AnotherValue sqlx.NullFloat64
}

func (t Test) Sql() *sqlx.TableSchema {
	return &sqlx.TableSchema{
		Name: "Tests",
		Columns: sqlx.Columns{
			{
				Name:  "Id",
				Type:  sqlx.CT().Int(),
				Key:   sqlx.K().Primary(),
				Extra: sqlx.E().AutoInc(true),
			}, {
				Name:     "DickValue",
				Type:     sqlx.CT().Int(),
				Nullable: true,
				Default:  sqlx.Int(5),
			}, {
				Name:     "StringValue",
				Type:     sqlx.CT().Varchar(64),
				Nullable: true,
				Default: sqlx.String(
					"some русская' string",
				),
			}, {
				Name:     "NewValue",
				Type:     sqlx.CT().Int(),
				Nullable: true,
				Default:  sqlx.Int(0),
			}, {
				Name:     "AnotherValue",
				Type:     sqlx.CT().Double(),
				Nullable: true,
				Default:  sqlx.Float(50),
			},
		},
	}
}

type AnotherTest struct {
	Id           sqlx.NullInt
	AnotherValue sqlx.NullInt
}

func (t AnotherTest) Sql() *sqlx.TableSchema {
	return &sqlx.TableSchema{
		OldName: "BetterTests",
		Name:    "AnotherTests",
		Columns: sqlx.Columns{
			{
				Name:     "Id",
				Type:     sqlx.CT().Int(),
				Nullable: false,
				Key:      sqlx.K().Primary(),
				Extra:    sqlx.E().AutoInc(true),
			}, {
				Name:     "AnotherValue",
				Type:     sqlx.CT().Int(),
				Nullable: true,
				Default:  sqlx.Int(25),
			},
		},
	}
}

func (t AnotherTest) BeforeInsert(
	db *sqlx.Db,
) error {
	fmt.Println("in this shit")
	if t.AnotherValue.Int32 > 5 {
		return errors.New("suck it")
	}

	return nil
}

func (t AnotherTest) AfterInsert(db *sqlx.Db) {
	fmt.Println("it must work")
}

func (t Test) AfterInsert(db *sqlx.Db) {
	fmt.Println("inserted test")
	fmt.Println(t)
}

var (
	Structs = []any{Test{}, AnotherTest{}}
)
