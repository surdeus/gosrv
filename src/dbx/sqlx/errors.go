package sqlx

import (
	"errors"
)

var (
	NoTablesSpecifiedErr = errors.New("no table specified")
	NoColumnsSpecifiedErr = errors.New("no columns specified")
	WrongNumOfColumnsSpecifiedErr = errors.New(
		"wrong number of columns specified")
	WrongQueryInputFormatErr = errors.New(
		"wrong query input format",
	)
	WrongValuerFormatErr = errors.New("wrong valuer format")
	UnknownQueryTypeErr = errors.New("unknown query type")
	UnknownConditionOpErr = errors.New("unknown condition operator")
	NoDBSpecifiedErr = errors.New("no database specified")
	NoSchemaSpecifiedErr = errors.New("no schema specified")
	WrongRawFormatErr = errors.New(
		"wrong raw value format error",
	)
	MultiplePrimaryKeysErr = errors.New("multiple primary keys")
	NoPrimaryKeySpecifiedErr = errors.New("no primary key specified")
	UnknownKeyTypeErr = errors.New(
		"unknown key type",
	)
	UnknownColumnTypeErr = errors.New(
		"unknown column type",
	)
	WrongColumnTypeFormatErr = errors.New(
		"wrong column type format",
	)
	TableDoesNotExistErr = errors.New(
		"specified table does not exist",
	)
	TableAlreadyExistsErr = errors.New(
		"specified table already exists",
	)
	ColumnDoesNotExistErr = errors.New(
		"specified column does not exist",
	)
	ColumnAlreadyExistsErr = errors.New(
		"specified column already exists",
	)
	UnknownColumnType = errors.New("unknown column type")
)

