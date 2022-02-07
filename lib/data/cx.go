package data

import (
	"github.com/jackc/pgx/v4"
)

type Cx struct {
	conn *pgx.Conn
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

func (self *Cx) NewTable(name string, primaryCols...Col) Table {
	t := new(BasicTable).Init(name, primaryCols...)
	self.tableLookup[name] = t
	return t
}

func (self *Cx) FindTable(name string) Table {
	return self.tableLookup[name]
}


