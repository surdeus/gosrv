package sqlx

type ConditionOp int
type Condition struct {
	Op ConditionOp
	Column ColumnName
	Value Valuer
}
type Conditions []Condition

const (
	EqConditionOp ConditionOp = iota
	GtConditionOp
	LtConditionOp
	GeConditionOp
	LeConditionOp
	NeConditionOp
	InConditionOp
)

var (
	// For the restx package.
	ConditionOpStringMap = map[string] ConditionOp {
		"eq" : EqConditionOp,
		"ne" : NeConditionOp,
		"gt" : GtConditionOp,
		"ge" : GeConditionOp,
		"lt" : LtConditionOp,
		"le" : LeConditionOp,
	}
	ConditionOpMap = map[ConditionOp] Raw {
		EqConditionOp : "=",
		NeConditionOp : "<>",
		GtConditionOp : ">",
		GeConditionOp : ">=",
		LtConditionOp : "<",
		LeConditionOp : "<=",
		InConditionOp : "in",
	}
)
func (w Conditions)SqlRaw(db *Db) (Raw, error) {
	if len(w) == 0 {
		return "", nil
	}

	ret := Raw("")
	for i, c := range w {
		op, ok := ConditionOpMap[c.Op]
		if !ok {
			return "", UnknownConditionOpErr
		}

		cond, err := db.Rprintf(
			"%s%s %s %s",
			Raw(" "),
			c.Column,
			op,
			Raw("?"),
		)
		if err != nil {
			return "", err
		}
		ret += cond

		if i < len(w)-1 {
			ret += Raw(" and")
		}
	}
	return Raw(ret), nil
}
