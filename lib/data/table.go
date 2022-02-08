package data

import (
	"fmt"
	"strings"
)

type Table struct {
	BasicDef
	BasicRel
	
	primaryKey *Key
	foreignKeys []*ForeignKey
	lookup map[string]Def
}

func (self *Table) Init(name string, keyCols...Col) *Table {
	self.BasicDef.Init(name)
	self.primaryKey = NewKey(fmt.Sprintf("%vPrimaryKey", name), keyCols...)
	self.lookup = make(map[string]Def)
	self.AddCols(keyCols...)
	return self
}

func (self *Table) AddCols(cols...Col) {
	self.BasicRel.AddCols(cols...)

	for _, c := range cols {
		self.lookup[c.Name()] = c
	}
}

func (self *Table) PrimaryKey() *Key {
	return self.primaryKey
}

func (self *Table) Create(cx *Cx) error {
	var sql strings.Builder
	fmt.Fprintf(&sql, "CREATE TABLE %v (", self.name)

	for i, c := range self.cols {
		if i > 0 {
			sql.WriteString(", ")
		}
		
		fmt.Fprintf(&sql, "%v %v NOT NULL", c.Name(), c.ValType())
	}
	
	sql.WriteRune(')')

	if err := cx.ExecSQL(sql.String()); err != nil {
		return err
	}
		
	if err := self.primaryKey.Create(cx, self); err != nil {
		return err
	}

	for _, k := range self.foreignKeys {
		if err := k.Create(cx, self); err != nil {
			return err
		}
	}

	return nil
}

func (self *Table) Exists(cx *Cx) (bool, error) {
	//TODO pg_tables
	return false, nil
}

func (self *Table) Drop(cx *Cx) error {
	for _, k := range self.foreignKeys {
		if err := k.Drop(cx, self); err != nil {
			return err
		}
	}

	sql := fmt.Sprintf("DROP TABLE IF EXISTS %v", self.name)

	if err := cx.ExecSQL(sql); err != nil {
		return err
	}

	return nil
}
