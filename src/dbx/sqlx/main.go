package sqlx

import (
	"database/sql"
	"fmt"
)

type Db struct {
	*sql.DB
	ConnConfig ConnConfig
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


func Open(cfg ConnConfig) (*Db, error) {
	db, err := sql.Open(cfg.Driver, cfg.String())
	if err != nil {
		return nil, err
	}

	return &Db{db, cfg}, nil
}

func (db *Db)Do(q Query) (*sql.Rows, error) {
	qs, err := q.SqlRaw(db)
	if err != nil {
		return nil, err
	}

	return db.DB.Query(string(qs))
}

