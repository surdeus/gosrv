package apix

// Implementing interface for

import (
	"github.com/surdeus/gosrv/src/httpx/muxx"
	"github.com/surdeus/gosrv/src/dbx/sqlx"
	//"database/sql"
	"encoding/gob"
	"io"
	"errors"
	"time"
)

type SqlResponseType int
type SqlConfig struct {
	Db *sqlx.Db
	Sqlers sqlx.Sqlers
}

const (
	NoSqlResponseType SqlResponseType = iota
	ErrorSqlResponseType
	RowsSqlResponseType
	ResultSqlResponseType
)

func SqlGobRegister() {
	gob.Register(sqlx.Byte(0))
	gob.Register(sqlx.Int16(0))
	gob.Register(sqlx.Int32(0))
	gob.Register(sqlx.Int64(0))
	gob.Register(sqlx.Float64(0))
	gob.Register(sqlx.String(""))
	gob.Register(sqlx.Time(time.Now()))
	gob.Register(errors.New(""))
	gob.Register(ErrorSqlResponseType)
}

func Sql(
	pref string,
	cfg SqlConfig,
) muxx.HndlDef {
	db := cfg.Db

	tMap := cfg.Sqlers.TableMap()
	tcMap := cfg.Sqlers.TableColumnMap()
	anyMap := cfg.Sqlers.AnyMap()

	for _, an := range anyMap {
		gob.Register(an)
	}

	postHndl := func(a muxx.HndlArg){
		dec := gob.NewDecoder(a.R.Body)
		q := sqlx.Query{}
		err := dec.Decode(&q)
		if err == io.EOF {
			return
		} else if err != nil {
			enc := gob.NewEncoder(a.W)
			enc.Encode(ErrorSqlResponseType)
			enc.Encode(err.Error())

			return
		}

		err = SqlHandleQuery(
			db, q, a,
			tMap, tcMap, anyMap,
		)
		if err != nil {
			enc := gob.NewEncoder(a.W)
			enc.Encode(ErrorSqlResponseType)
			enc.Encode(err.Error())

			return
		}
	}

	def := muxx.HndlDef {
		pref,
		"^$",
		muxx.Handlers{
			"POST" : postHndl,
		},
	}

	SqlGobRegister()

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

		cMap, ok := tcMap[tname]
		if !ok {
			return sqlx.
				TableDoesNotExistErr
		}

		an, ok := anyMap[tname]
		if !ok {
			return sqlx.
				TableDoesNotExistErr
		}

		cnames := q.GetColumnNames()
		if len(cnames) == 1 &&
				cnames[0] == "*" {
			cnames = tMap[tname].Columns.Names()
			q.ColumnNames = cnames
		}

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
		enc.Encode(RowsSqlResponseType)
		for v := range values {
			err = enc.Encode(v)
			if err != nil {
				return err
			}
		}
	default :

		result, _, err := db.Do(q)
		if err != nil {
			return err
		}

		r := sqlx.Result{}
		r.LastInsertId, _ = result.LastInsertId()
		r.RowsAffected, _ = result.RowsAffected()

		enc := gob.NewEncoder(a.W)
		enc.Encode(ResultSqlResponseType)

		err = enc.Encode(r)
		if err != nil {
			return err
		}
	}

	return nil
}

