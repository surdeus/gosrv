package sqlx

// Simple package embedding value recievers
// for more friendly usage.

import (
	"database/sql/driver"
)

type Valuer = driver.Valuer
type Valuers []Valuer

