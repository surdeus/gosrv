package apix

// Implementing interface for

import (
	"github.com/surdeus/gosrv/src/httpx/muxx"
	"github.com/surdeus/gosrv/src/dbx/sqlx"
)

type SqlConfig struct {
	Db *sqlx.Db
}

func Sql(
	pref string,
	cfg SqlConfig,
) muxx.HndlDef {
	def := muxx.HndlDef {

	}
}
