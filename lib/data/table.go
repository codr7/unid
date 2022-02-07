package data

import (
	"fmt"
	"github.com/jackc/pgx/v4"
)

type Table interface {
	Def
	Rel

	PrimaryKey() Key
}

type BasicTable struct {
	BasicDef
	BasicRel
	
	primaryKey Key
	otherKeys []Key
	lookup map[string]Def
}

func (self *BasicTable) Init(name string, keyCols...Col) *BasicTable {
	self.BasicDef.Init(name)
	self.primaryKey = NewKey(fmt.Sprintf("%vPrimaryKey", name), keyCols...)
	self.lookup = make(map[string]Def)
	return self
}

func (self *BasicTable) AddCols(cols...Col) {
	self.BasicRel.AddCols(cols...)

	for _, c := range cols {
		self.lookup[c.Name()] = c
	}
}

func (self *BasicTable) PrimaryKey() Key {
	return self.primaryKey
}

func (self *BasicTable) NewForeignKey(name string, foreignTable Table) *ForeignKey {
	k := new(ForeignKey).Init(name, foreignTable)
	self.otherKeys = append(self.otherKeys, k)
	self.lookup[name] = k
	return k
}

func (self *BasicTable) Scan(row pgx.Row) ([]interface{}, error) {
	out := make([]interface{}, len(self.Cols()))

	for i, c := range self.Cols() {
		out[i] = c.InitField()
	}
	
	return out, row.Scan(out...)
}
