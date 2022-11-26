package structs

import (
	"github.com/surdeus/go-srv/src/dbx/sqlx"
)


type Test struct {
	Id int
	Value int
	StringValue string
}

func (t Test)Sql() *sqlx.TableSchema {
	return &sqlx.TableSchema {
		OldName: "NewTests",
		Name: "Tests",
		Columns: sqlx.Columns {
			{
				Name: "Id",
				Type: sqlx.CT().Int(),
				Nullable: false,
				Key: sqlx.PrimaryKey(),
				Extra: "auto_increment",
			},{
				//OldName: "SuckValue",
				//OldName: "NewValueName",
				OldName: "KillValue",
				Name: "DickValue",
				Type: sqlx.CT().Int(),
				Nullable: true,
				Default: sqlx.Int(5),
			},{
				Name: "StringValue",
				Type: sqlx.CT().Varchar(32),
				Nullable: true,
				Default: sqlx.String(
					"some русская' string"),
			},{
				Name: "NewValue",
				Type: sqlx.CT().Int(),
				Nullable: true,
				Default: sqlx.Int(0),
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
