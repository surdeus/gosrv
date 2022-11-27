package sqlx

type Column struct {
	OldName ColumnName
	Name ColumnName
	Type ColumnType
	Nullable bool
	Key Key
	Default Rawer
	Extra Raw
}

type Columns []*Column

func C() *Column {
	return &Column{}
}

func (c *Column) WOldName(n ColumnName) *Column {
	c.OldName = n
	return c
}

func (c *Column) WName(n ColumnName) *Column {
	c.Name = n
	return c
}

func (c *Column) WNullable() *Column {
	c.Nullable = true
	return c
}

func (c *Column) WType(t ColumnType) *Column {
	c.Type = t
	return c
}

func (c *Column) WDefault(d Raw) *Column {
	c.Default = d
	return c
}

func (c *Column) WKey(k Key) *Column {
	c.Key = k
	return c
}

func (c *Column) WExtra(e Raw) *Column {
	c.Extra = e
	return c
}
