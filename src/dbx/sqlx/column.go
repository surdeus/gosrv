package sqlx

type Column struct {
	OldName ColumnName
	Name ColumnName
	Type ColumnType
	Nullable bool
	Key Key
	Default Valuer
	Extra Raw
}

type Columns []*Column

func C() *Column {
	return &Column{}
}

