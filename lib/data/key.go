package data

import (
	"fmt"
	"strings"
)

type Key struct {
	BasicDef
	BasicRel
}

func NewKey(name string, cols...Col) *Key {
	return new(Key).Init(name, cols...)
}

func (self *Key) Init(name string, cols...Col) *Key {
	self.BasicDef.Init(name)
	self.BasicRel.Init()
	self.AddCols(cols...)
	return self
}

func (self *Key) Create(table *Table) error {
	ct := "UNIQUE"

	if table.primaryKey == self {
		ct = "PRIMARY KEY"
	}
	
	var sql strings.Builder
	fmt.Fprintf(&sql, "ALTER TABLE %v ADD CONSTRAINT %v %v (", table.name, self.name, ct)

	for i, c := range self.cols {
		if i > 0 {
			sql.WriteString(", ")
		}

		sql.WriteString(c.Name())
	}
	
	sql.WriteRune(')')

	if err := table.Cx().ExecSQL(sql.String()); err != nil {
		return err
	}

	return nil
}

func (self *Key) Drop(table *Table) error {
	sql := fmt.Sprintf("ALTER TABLE %v DROP CONSTRAINT IF EXISTS %v", table.name, self.name)

	if err := table.Cx().ExecSQL(sql); err != nil {
		return err
	}

	return nil
}
