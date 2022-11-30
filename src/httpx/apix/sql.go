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
	"log"
	"bytes"
	"net/http"
	"reflect"
	//"fmt"
)

type SqlResponseType int

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
	/*for _, v := range db.AMap {
		gob.Register(v)
	}*/
}

func Sql(
	db *sqlx.Db,
	pref string,
) muxx.HndlDef {
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

		err = SqlHandleQuery(db, q, a)
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
) error {
	_, err := q.SqlRaw(db)
	if err != nil {
		return err
	}

	switch q.Type {
	case sqlx.SelectQueryType :

		tname := q.GetTableName()

		cnames := q.GetColumnNames()
		if len(cnames) == 1 &&
				cnames[0] == "*" {
			cnames = db.TMap[tname].Columns.Names()
			q.ColumnNames = cnames
		}

		_, rs, err := db.Do(q)
		if err != nil {
			return err
		}
		defer rs.Close()

		values, err := db.ReadRowValues(
			rs,
			tname,
			cnames,
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

func SqlQuery(
	u string,
	q sqlx.Query,
	rc any,
) (sqlx.Result, chan any, error) {
	nilRes := sqlx.Result{}
	bts := bytes.NewBuffer([]byte{})
	enc := gob.NewEncoder(bts)

	err := enc.Encode(q)
	if err != nil {
		return nilRes, nil, err
	}
	resp, err := http.Post(
		u,
		"application/gob",
		bts)
	if err != nil {
		return nilRes, nil, err
	}

	dec := gob.NewDecoder(resp.Body)

	typ := NoSqlResponseType
	err = dec.Decode(&typ)
	if err != nil {
		println(err.Error())
		return nilRes, nil, err
	}

	switch typ {
	case ErrorSqlResponseType :
		var errbuf string
		err = dec.Decode(&errbuf)
		if err != nil {
			return nilRes, nil, err
		}
		err = errors.New(errbuf)
		return nilRes, nil, err
	case RowsSqlResponseType :
		chn := make(chan any)
		go func() {
			for {
				err = dec.Decode(rc)
				if err == io.EOF {
					break
				} else if err != nil {
					log.Println(err)
				}
				chn <- reflect.
					Indirect(
						reflect.
						ValueOf(rc),
					).Interface()
			}
			close(chn)
		}()
		return nilRes, chn, nil
	case ResultSqlResponseType :
		var buf sqlx.Result
		err = dec.Decode(&buf)
		if err != nil {
			return nilRes, nil, err
		}
		return buf, nil, nil
	}
	return nilRes, nil, nil
}
