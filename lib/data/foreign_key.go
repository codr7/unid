package data

import (
	"fmt"
	"strings"
)

type ForeignKey struct {
	Key
	foreignTable *Table
}

func (self *ForeignKey) Init(name string, foreignTable *Table) *ForeignKey {
	self.Key.Init(name)
	self.foreignTable = foreignTable

	for _, c := range foreignTable.PrimaryKey().Cols() {
		self.AddCols(c.Clone(fmt.Sprintf("%v%v", name, c.Name())))
	}

	return self
}

func (self *Table) NewForeignKey(name string, foreignTable *Table) *ForeignKey {
	k := new(ForeignKey).Init(fmt.Sprintf("%v%vKey", self.name, name), foreignTable)
	self.foreignKeys = append(self.foreignKeys, k)
	self.lookup[name] = k
	return k
}

func (self *ForeignKey) ForeignTable() *Table {
	return self.foreignTable
}

func (self *ForeignKey) Create(cx *Cx, table *Table) error {
	var sql strings.Builder
	fmt.Fprintf(&sql, "ALTER TABLE %v ADD CONSTRAINT %v FOREIGN KEY (", table.name, self.name)

	for i, c := range self.cols {
		if i > 0 {
			sql.WriteString(", ")
		}

		sql.WriteString(c.Name())
	}
	
	fmt.Fprintf(&sql, ") REFERENCES %v (", self.foreignTable.name)

	for i, c := range self.foreignTable.primaryKey.cols {
		if i > 0 {
			sql.WriteString(", ")
		}
		
		sql.WriteString(c.Name())
	}

	sql.WriteRune(')')

	if err := cx.ExecSQL(sql.String()); err != nil {
		return err
	}

	return nil
}

