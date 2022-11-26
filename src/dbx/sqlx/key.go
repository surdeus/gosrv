package sqlx

func K() Key {
	return Key{}
}

func (k Key) Primary() Key {
	return PrimaryKey()
}

