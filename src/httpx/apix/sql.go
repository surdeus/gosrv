package apix

// Implementing interface for

import (
	"github.com/surdeus/gosrv/src/httpx/muxx"
	"github.com/surdeus/gosrv/src/dbx/sqlx"
	//"database/sql"
	"encoding/gob"
	"io"
)

type SqlConfig struct {
	Db *sqlx.Db
	Sqlers sqlx.Sqlers
}

func Sql(
	pref string,
	cfg SqlConfig,
) muxx.HndlDef {
	gob.Register(sqlx.Int(0))
	db := cfg.Db

	tMap := cfg.Sqlers.TableMap()
	tcMap := cfg.Sqlers.TableColumnMap()
	anyMap := cfg.Sqlers.AnyMap()

	postHndl := func(a muxx.HndlArg){
		dec := gob.NewDecoder(a.R.Body)
		q := sqlx.Query{}
		for {
			err := dec.Decode(&q)
			if err == io.EOF {
				break
			} else if err != nil {
				a.ServerError(err)
				return
			}
			err = SqlHandleQuery(
				db, q, a,
				tMap, tcMap, anyMap,
			)
			if err != nil {
				a.ServerError(err)
				return
			}
		}
	}

	def := muxx.HndlDef {
		pref,
		"^$",
		muxx.Handlers{
			"POST" : postHndl,
		},
	}

	return def
}

func SqlHandleQuery(
	db *sqlx.Db,
	q sqlx.Query,
	a muxx.HndlArg,
	tMap sqlx.TableMap,
	tcMap sqlx.TableColumnMap,
	anyMap sqlx.AnyMap,
) error {
	_, err := q.SqlRaw(db)
	if err != nil {
		return err
	}

	switch q.Type {
	case sqlx.SelectQueryType :
		tname := q.GetTableName()

		cMap := tcMap[tname]

		an, ok := anyMap[tname]
		if !ok {
			return sqlx.
				TableDoesNotExistErr
		}

		cnames := q.GetColumnNames()

		_, rs, err := db.Do(q)
		if err != nil {
			return err
		}

		values, err := db.ReadRowValues(
			rs,
			cnames,
			cMap,
			an,
		)
		if err != nil {
			return err
		}

		enc := gob.NewEncoder(a.W)
		for v := range values {
			err = enc.Encode(v)
			if err != nil {
				return err
			}
		}
	default :
		return sqlx.UnknownQueryTypeErr
	}

	return nil
}

