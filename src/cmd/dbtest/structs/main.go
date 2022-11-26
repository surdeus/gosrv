package structs

import (
	"github.com/surdeus/go-srv/src/dbx/sqlx"
	"database/sql"
)


type Test struct {
	Id sql.NullInt32
	DickValue sql.NullInt32
	StringValue sql.NullString
	NewValue sql.NullInt32
	AnotherValue sql.NullInt32
}

func (t Test)Sql() *sqlx.TableSchema {
	return &sqlx.TableSchema {
		OldName: "NewTests",
		Name: "Tests",
		Columns: sqlx.Columns {
			{
				Name: "Id",
				Type: sqlx.CT().Int(),
				Key: sqlx.PrimaryKey(),
				Extra: "auto_increment",
			},{
				Name: "DickValue",
				Type: sqlx.CT().Int(),
				Nullable: true,
				Default: sqlx.Int(5),
			},{
				Name: "StringValue",
				Type: sqlx.CT().Varchar(32),
				Nullable: false,
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
				Type: sqlx.CT().Int(),
				Nullable: true,
				Default: sqlx.Int(100),
			},
		},
	}
}

type AnotherTest struct {
	Id int
	Value int
	StringValue string
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
				Key: sqlx.PrimaryKey(),
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
