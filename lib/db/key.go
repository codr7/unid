package db

import (
	"fmt"
	"strings"
)

type Key struct {
	BasicCols
	BasicDef
}

func NewKey(name string, cols...Col) *Key {
	return new(Key).Init(name, cols...)
}

func (self *Key) Init(name string, cols...Col) *Key {
	self.BasicCols.Init()
	self.BasicDef.Init(name)
	self.AddCol(cols...)
	return self
}

func (self *Key) Create(table Table) error {
	ct := "UNIQUE"

	if table.PrimaryKey() == self {
		ct = "PRIMARY KEY"
	}
	
	var sql strings.Builder
	fmt.Fprintf(&sql, "ALTER TABLE \"%v\" ADD CONSTRAINT \"%v\" %v (", table.Name(), self.name, ct)

	for i, c := range self.cols {
		if i > 0 {
			sql.WriteString(", ")
		}

		fmt.Fprintf(&sql, "\"%v\"", c.Name())
	}
	
	sql.WriteRune(')')

	if err := table.Cx().ExecSQL(sql.String()); err != nil {
		return err
	}

	return nil
}

func (self *Key) Drop(table Table) error {
	sql := fmt.Sprintf("ALTER TABLE \"%v\" DROP CONSTRAINT IF EXISTS \"%v\"", table.Name(), self.name)

	if err := table.Cx().ExecSQL(sql); err != nil {
		return err
	}

	return nil
}
