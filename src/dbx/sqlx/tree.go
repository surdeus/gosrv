package sqlx

import (
	"strings"
	"fmt"
)

const (
	noOp TreeOp = iota
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
	sumOp
	subOp
	mulOp
	divOp
	modOp
	bAndOp
	bOrOp
	bXorOp
)

var (
	// For the restx package.
	TreeOpStringMap = map[string] TreeOp {
		"eq" : eqOp,
		"ne" : neOp,
		"gt" : gtOp,
		"ge" : geOp,
		"lt" : ltOp,
		"le" : leOp,
	}
	TreeOpMap = map[TreeOp] Raw {
		eqOp: "=",
		neOp: "<>",
		gtOp: ">",
		geOp: ">=",
		ltOp: "<",
		leOp: "<=",
		inOp: "in",
		orOp: "or",
		andOp: "and",
		sumOp : "+",
		subOp : "-",
		mulOp : "*",
		divOp : "/",
		modOp : "%",
		bAndOp : "&",
		bOrOp : "|",
		bXorOp : "^",
	}
)

func (c TreeOp)SqlRaw(db *Db) (Raw, error) {
	ret, ok := TreeOpMap[c]
	if !ok {
		return "", UnknownTreeOpErr
	}

	return Raw(ret), nil
}

func T() Tree {
	return Tree{Pair: make([]Tree, 2)}
}


func (c Tree)And(
	cs ...Tree,
) Tree {
	return c.opMul(andOp, cs...)
}

/*func (c Tree) opMul(
	op TreeOp,
	cs ...Tree,
) Tree {
	if len(cs) < 2 {
		if len(cs) == 1 {
			return cs[0]
		} else {
			T()
		}
	}

	if len(cs) == 2 {
		c.Op = op
		c.Pair = []Tree{cs[0], cs[1]}
		return c
	}

	c.Op = op

	c.Pair[0] = cs[0]
	cs = cs[1:]

	cn := T().opMul(op, cs...)
	c.Pair[1] = cn

	return c
}*/

func (t Tree) opMul(
	op TreeOp,
	cs ...Tree,
) Tree {
	t.Op = op
	t.Pair = cs
	return t
}

func (c Tree)Or(
	cs ...Tree,
) Tree {
	c = c.opMul(orOp, cs...)
	return c
}

func (c Tree) Gt() Tree {
	c.Op = gtOp
	return c
}

func (c Tree) Lt() Tree {
	c.Op = ltOp
	return c
}

func (c Tree) Eq() Tree {
	c.Op = eqOp
	return c
}

func (c Tree) Le() Tree {
	c.Op = leOp
	return c
}

func (c Tree) Ge() Tree {
	c.Op = geOp
	return c
}

func (c Tree) In() Tree {
	c.Op = inOp
	return c
}

func (t Tree) Sum(
	cs ...Tree,
) Tree {
	t = t.opMul(sumOp, cs...)
	return t
}

func (t Tree) Mul(
	cs ...Tree,
) Tree {
	t = t.opMul(mulOp, cs...)
	return t
}

func (t Tree) Sub(
	cs ...Tree,
) Tree {
	t = t.opMul(subOp, cs...)
	return t
}

func (c Tree) C1(
	name ColumnName,
) Tree {
	c0 := T().C(name)
	c.Pair[0] = c0
	return c
}

func (c Tree) C2(
	name ColumnName,
) Tree {
	c.Pair[1] = T().C(name)
	return c
}

func (c Tree) V1(
	v ...Valuer,
) Tree {
	c.Pair[0] = T().V(v...)
	return c
}

func (c Tree) V2 (
	v ...Valuer,
) Tree {
	c.Pair[1] = T().V(v...)
	return c
}

func (c Tree) C(
	name ColumnName,
) Tree {
	c.Op = colOp
	c.Column = name
	c.Pair = []Tree{}
	return c
}

func (c Tree) V(
	v ...Valuer,
) Tree {
	c.Op = valOp
	c.Pair = []Tree{}
	c.Values = v
	return c
}

func (c Tree) values() Valuers {
	if c.Op == valOp {
		return c.Values
	}

	if c.Op == colOp {
		return Valuers{}
	}

	var ret Valuers
	if len(c.Pair) == 2 {
		ret = append(ret, c.Pair[0].values()...)
		ret = append(ret, c.Pair[1].values()...)
	}

	return ret
}

func (t Tree)SqlRaw(db *Db) (Raw, error) {

	switch t.Op {
	case valOp :
		return db.TupleBuf(t.Values), nil
	case colOp :
		return t.Column.SqlRaw(db)
	}

	if len(t.Pair) == 0 {
		return "", nil
	}

	var buf strings.Builder
	fmt.Fprint(&buf, "(")
	for i, p := range t.Pair {
		v, e := db.Rprintf("%s", p)
		if e != nil {
			return "", e
		}

		fmt.Fprint(&buf, v)


		if i < len(t.Pair) - 1 {
			v, e = db.Rprintf("%s", t.Op)
			if e != nil {
				return "", e
			}
			fmt.Fprint(&buf, " ")
			fmt.Fprint(&buf, v)
			fmt.Fprint(&buf, " ")
		}
	}
	fmt.Fprint(&buf, ")")

	ret := Raw(buf.String())

	return ret, nil
}
