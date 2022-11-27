package sqlx

type ConditionOp int
type Condition struct {
	Column ColumnName
	Op ConditionOp
	Values Valuers
}
type Conditions []Condition

const (
	Eq ConditionOp = iota
	Gt
	Lt
	Ge
	Le
	Ne
	In
)

var (
	// For the restx package.
	ConditionOpStringMap = map[string] ConditionOp {
		"eq" : Eq,
		"ne" : Ne,
		"gt" : Gt,
		"ge" : Ge,
		"lt" : Lt,
		"le" : Le,
	}
	ConditionOpMap = map[ConditionOp] Raw {
		Eq: "=",
		Ne: "<>",
		Gt: ">",
		Ge: ">=",
		Lt: "<",
		Le: "<=",
		In: "in",
	}
)
func (w Conditions)SqlRaw(db *Db) (Raw, error) {
	if len(w) == 0 {
		return "", nil
	}

	ret := Raw("")
	prespace := Raw("")
	for i, c := range w {
		op, ok := ConditionOpMap[c.Op]
		if !ok {
			return "", UnknownConditionOpErr
		}

		cond, err := db.Rprintf(
			"%s%s %s (%s)",
			prespace,
			c.Column,
			op,
			db.MultiBuf(c.Values),
		)
		if err != nil {
			return "", err
		}
		prespace = " "
		ret += cond

		if i < len(w)-1 {
			ret += Raw(" and")
		}
	}
	return Raw(ret), nil
}
