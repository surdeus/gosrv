package sqlx

func E() ExtraColInfo {
	return ExtraColInfo{}
}

func (e ExtraColInfo) AutoInc(
	v bool,
) ExtraColInfo {
	e.AutoIncrement = v
	return e
}

