package sqlx

import (
	"fmt"
	"strings"
)

// Substitute raw values with fmt.Sprintf
// and SqlRaw.
func (db *Db)Rprintf(
	format string,
	rawers ...Rawer,
) (Raw, error) {
	raws := []any{}
	for _, v := range rawers {
		raw, err := v.SqlRaw(db)
		if err != nil {
			return "", err
		}
		raws = append(raws, any(raw))
	}

	ret := fmt.Sprintf(format, raws...)

	return Raw(ret), nil
}

// Return raw string with buffer for Valuer insertion
// in SQL queries.
func (db *Db)MultiBuf(vs Valuers) Raw {
	if len(vs) == 0 {
		return ""
	}
	buf := make([]string, len(vs))
	for i := range buf {
		buf[i] = "?"
	}

	return Raw(strings.Join(buf, ","))
}

func (db *Db)TupleBuf(vs Valuers) Raw {
	ret := db.MultiBuf(vs)
	if ret == "" {
		return ""
	}

	return Raw(fmt.Sprintf("(%s)", ret))
}

func (v TableName)SqlRaw(db *Db) (Raw, error) {
	if v == "" {
		return "", WrongValuerFormatErr
	}
	return Raw(v), nil
}

func (v ColumnName)SqlRaw(db *Db) (Raw, error) {
	if v == "" {
		return "", WrongValuerFormatErr
	}
	return Raw(v), nil
}

func (v Raw)SqlRaw(db *Db) (Raw, error) {
	return v, nil
}


func (tn TableNames)SqlRaw(db *Db) (Raw, error) {
	if len(tn) == 0 {
		return Raw(""), NoTablesSpecifiedErr
	}

	buf := make([]string, 0)
	for _, t := range tn {
		v, err := t.SqlRaw(db)
		if err != nil {
			return Raw(""), err
		}
		buf = append(buf, string(v))
	}

	ret := strings.Join(buf, ", ")
	return Raw(ret), nil
}

func (cn ColumnNames)SqlRaw(db *Db) (Raw, error) {
	if len(cn) == 0 {
		return Raw(""), NoColumnsSpecifiedErr
	}

	buf := make([]string, 0)
	for _, c := range cn {
		v, err := c.SqlRaw(db)
		if err != nil {
			return Raw(""), err
		}
		buf = append(buf, string(v))
	}

	ret := strings.Join(buf, ", ")
	return Raw(ret), nil
}

// Return raw values separated by comma for
// column and table names and also values.
func (rvs Rawers) SqlMultival(db *Db) (Raw, error) {
	var ret Raw
	for i, v := range rvs {
		raw, err := v.SqlRaw(db)
		if err != nil {
			return "", err
		}

		ret += raw

		if i != len(rvs) - 1 {
			ret += ","
		}
	}

	return ret, nil
}

// Return multivalue embraced with () .
func (rvs Rawers) SqlTuple(db *Db) (Raw, error) {
	mval, err := rvs.SqlMultival(db)
	if err != nil {
		return Raw(""), err
	}

	if mval == "" {
		return "", nil
	}

	return Raw(fmt.Sprintf("(%s)", mval)), nil
}

func (db *Db)RawersEq(
	v1, v2 Rawer,
) (bool, error) {

	if v1 == nil && v2 == nil {
		return true, nil
	}

	if v1 == nil || v2 == nil {
		fmt.Println("in")
		return false, nil
	}

	raw1, err := v1.SqlRaw(db)
	if err != nil {
		return false, err
	}

	raw2, err := v2.SqlRaw(db)
	if err != nil {
		return false, err
	}
	return raw1 == raw2, nil
}

