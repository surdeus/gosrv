package sqlx

import (
	"fmt"
)

// Substitute raw values with fmt.Sprintf
// and SqlRaw.
func (db *DB)Rprintf(
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

