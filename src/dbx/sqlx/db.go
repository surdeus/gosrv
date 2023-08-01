package sqlx

import (
	"database/sql"
	"fmt"
	"log"
)

type Db struct {
	*sql.DB
	Debug bool
	ConnConfig ConnConfig
	Tables TableSchemas
	TMap TableMap
	//TCMap TableColumnMap
	TypeMap TypeMap
}

type ConnConfig struct {
	Driver string
	Login, Password, Host, Name string
	Port int
}

func (c ConnConfig)String() string {
	return fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s",
		c.Login,
		c.Password,
		c.Host,
		c.Port,
		c.Name,
	)
}

func Open(cfg ConnConfig, sqlers Sqlers) (*Db, error) {
	db, err := sql.Open(cfg.Driver, cfg.String())
	if err != nil {
		return nil, err
	}

	if len(sqlers) < 1 {
		return nil, NoTablesSpecifiedErr
	}

	tables := sqlers.Tables()
	tMap := tables.TableMap()
	typeMap := tables.TypeMap()

	return &Db{
		DB: db,
		ConnConfig: cfg,
		Tables: tables,
		TMap: tMap,
		TypeMap: typeMap,
	}, nil
}

func (db *Db)Do(
	q Query,
) (sql.Result, *sql.Rows, error) {
	//var val Sqler
	qs, err := q.SqlRaw(db)
	if err != nil {
		return nil, nil, err
	}
	
	if db.Debug {
		log.Printf("Handling the %q request...", string(qs))
	}
	
	switch q.Type {
	case SelectQueryType :
		rs, err := db.DB.Query(string(qs), q.GetValues()...)
		return nil, rs, err
	case InsertQueryType :
		v, err := db.ConstructInsertValue(q)
		if err != nil {
			return nil, nil, err
		}
		bef, ok := any(v).(interface{
			BeforeInsert(*Db) error
		})

		fmt.Println("do: in", ok)
		if ok {
			err := bef.BeforeInsert(db)
			if err != nil {
				return nil, nil, err
			}
		}

		res, err := db.DB.Exec(string(qs), q.GetValues()...)
		if err != nil {
			return nil, nil, err
		}

		aft, ok := any(v).(interface{
			AfterInsert(*Db)
		})
		if ok {
			aft.AfterInsert(db)
		}

		return res, nil, nil
	}

	res, err := db.DB.Exec(string(qs), q.GetValues()...)
	if err != nil {
		return nil, nil, err
	}

	return res, nil, nil
}

