package sqlx

const (
	noOp ConditionOp = iota
	eqOp
	gtOp
	ltOp
	geOp
	leOp
	neOp
	inOp
	orOp
	andOp
	valOp
	colOp
)

var (
	// For the restx package.
	ConditionOpStringMap = map[string] ConditionOp {
		"eq" : eqOp,
		"ne" : neOp,
		"gt" : gtOp,
		"ge" : geOp,
		"lt" : ltOp,
		"le" : leOp,
	}
	ConditionOpMap = map[ConditionOp] Raw {
		eqOp: "=",
		neOp: "<>",
		gtOp: ">",
		geOp: ">=",
		ltOp: "<",
		leOp: "<=",
		inOp: "in",
		orOp: "or",
		andOp: "and",
	}
)

func (c ConditionOp)SqlRaw(db *Db) (Raw, error) {
	ret, ok := ConditionOpMap[c]
	if !ok {
		return "", UnknownConditionOpErr
	}

	return Raw(ret), nil
}

func C() Condition {
	return Condition{}
}


func (c Condition)And(
	cs ...Condition,
) Condition {
	return c.opMul(andOp, cs...)
}

func (c Condition) opMul(
	op ConditionOp,
	cs ...Condition,
) Condition {
	if len(cs) < 2 {
		if len(cs) == 1 {
			return cs[0]
		} else {
			C()
		}
	}

	if len(cs) == 2 {
		c.Op = op
		c.Pair[0] = &cs[0]
		c.Pair[1] = &cs[1]
		return c
	}

	c.Op = op

	c.Pair[0] = &cs[0]
	cs = cs[1:]

	cn := C().opMul(op, cs...)
	c.Pair[1] = &cn

	return c
}

func (c Condition)Or(
	cs ...Condition,
) Condition {
	c = c.opMul(orOp, cs...)
	return c
}

func (c Condition) Gt() Condition {
	c.Op = gtOp
	return c
}

func (c Condition) Lt() Condition {
	c.Op = ltOp
	return c
}

func (c Condition) Eq() Condition {
	c.Op = eqOp
	return c
}

func (c Condition) Le() Condition {
	c.Op = leOp
	return c
}

func (c Condition) Ge() Condition {
	c.Op = geOp
	return c
}

func (c Condition) In() Condition {
	c.Op = inOp
	return c
}

func (c Condition) C1(
	name ColumnName,
) Condition {
	c0 := C().C(name)
	c.Pair[0] = &c0
	return c
}

func (c Condition) C2(
	name ColumnName,
) Condition {
	c1 := C().C(name)
	c.Pair[1] = &c1
	return c
}

func (c Condition) V1(
	v ...Valuer,
) Condition {
	v0 := C().V(v...)
	c.Pair[0] = &v0
	return c
}

func (c Condition) V2 (
	v ...Valuer,
) Condition {
	v1 := C().V(v...)
	c.Pair[1] = &v1
	return c
}

func (c Condition) S (
	c1, c2 Condition,
) Condition {
	c.Pair = [2]*Condition{&c1, &c2}

	return c
}

func (c Condition) C(
	name ColumnName,
) Condition {
	c.Op = colOp
	c.Column = name
	return c
}

func (c Condition) V(
	v ...Valuer,
) Condition {
	c.Op = valOp
	c.Values = v
	return c
}

func (c Condition) values() Valuers {
	if c.Op == valOp {
		return c.Values
	}

	if c.Op == colOp {
		return Valuers{}
	}

	ret := c.Pair[0].values()
	ret = append(ret, c.Pair[1].values()...)

	return ret
}

func (c Condition)SqlRaw(db *Db) (Raw, error) {

	switch c.Op {
	case valOp :
		return db.TupleBuf(c.Values), nil
	case colOp :
		return c.Column.SqlRaw(db)
	}

	return db.Rprintf(
		"(%s %s %s)",
		*c.Pair[0],
		c.Op,
		*c.Pair[1],
	)
}
