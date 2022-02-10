package data

import (
	"context"
	"github.com/jackc/pgx/v4"
	"log"
)

type Cx struct {
	conn *pgx.Conn
	defs []RootDef
	tableLookup map[string]Table
}

func NewCx(conn *pgx.Conn) *Cx {
	return new(Cx).Init(conn)
}

func (self *Cx) Init(conn *pgx.Conn) *Cx {
	self.conn = conn
	self.tableLookup = make(map[string]Table)
	return self
}

func (self *Cx) NewTable(name string) Table {
	t := new(BasicTable).Init(self, name)
	self.tableLookup[name] = t
	self.defs = append(self.defs, t)
	return t
}

func (self *Cx) FindTable(name string) Table {
	return self.tableLookup[name]
}

func (self *Cx) ExecSQL(sql string, params...interface{}) error {
	log.Printf("%v\n%v\n", sql, params)
	_, err := self.conn.Exec(context.Background(), sql, params...)
	return err
}

func (self *Cx) QueryRow(sql string, params...interface{}) pgx.Row {
	log.Printf("%v\n%v\n", sql, params)
	return self.conn.QueryRow(context.Background(), sql, params...)
}

func (self *Cx) Query(sql string, params...interface{}) (pgx.Rows, error) {
	log.Printf("%v\n%v\n", sql, params)
	return self.conn.Query(context.Background(), sql, params...)
}

func (self *Cx) SyncAll() error {
	for _, d := range self.defs {
		if ok, err := d.Exists(); err != nil {
			return err
		}  else if !ok {
			if err := d.Create(); err != nil {
				return err
			}
		}
	}
	
	return nil
}

func (self *Cx) DropAll() error {
	for i := range self.defs {
		d := self.defs[len(self.defs)-i-1]
		
		if ok, err := d.Exists(); err != nil {
			return err
		}  else if ok {
			if err := d.Drop(); err != nil {
				return err
			}
		}
	}

	return nil
}
