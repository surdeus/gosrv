package structs

import (
	"github.com/surdeus/gosrv/src/dbx/sqlx"
	"database/sql"
	//"reflect"
)


type Test struct {
	Id sql.NullInt32
	DickValue sql.NullInt32
	StringValue sql.NullString
	NewValue sql.NullInt32
	AnotherValue sql.NullFloat64
}

func (t Test)Sql() *sqlx.TableSchema {
	return &sqlx.TableSchema {
		Name: "Tests",
		Columns: sqlx.Columns {
			{
				Name: "Id",
				Type: sqlx.CT().Int(),
				Key: sqlx.K().Primary(),
				Extra: "auto_increment",
			},{
				Name: "DickValue",
				Type: sqlx.CT().Int(),
				Nullable: true,
				Default: sqlx.Int(5),
			},{
				Name: "StringValue",
				Type: sqlx.CT().Varchar(64),
				Nullable: true,
				Default: sqlx.String(
					"some русская' string",
				),
			},{
				Name: "NewValue",
				Type: sqlx.CT().Int(),
				Nullable: true,
				Default: sqlx.Int(0),
			},{
				Name: "AnotherValue",
				Type: sqlx.CT().Double(),
				Nullable: true,
				Default: sqlx.Float(50),
			},
		},
	}
}

type AnotherTest struct {
	Id int
	AnotherValue int
}

func (t AnotherTest)Sql() *sqlx.TableSchema {
	return &sqlx.TableSchema {
		OldName: "BetterTests",
		Name: "AnotherTests",
		Columns: sqlx.Columns {
			{
				Name: "Id",
				Type: sqlx.CT().Int(),
				Nullable: false,
				Key: sqlx.K().Primary(),
				//Extra: "auto_increment",
			},{
				Name: "AnotherValue",
				Type: sqlx.CT().Int(),
				Nullable: true,
				Default: sqlx.Int(25),
			},
		},
	}
}

var (
	Structs = []any{Test{}, AnotherTest{}}
)
