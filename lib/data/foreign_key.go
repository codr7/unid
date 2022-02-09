package data

import (
	"fmt"
	"strings"
)

type ForeignKey struct {
	BasicField
	Key
	foreignTable *Table
}

func (self *ForeignKey) Init(name string, fieldName string, foreignTable *Table) *ForeignKey {
	self.BasicField.Init(fieldName)
	self.Key.Init(name)
	self.foreignTable = foreignTable

	for _, c := range foreignTable.PrimaryKey().Cols() {
		self.AddCol(c.NewForeignCol(fmt.Sprintf("%v%v", fieldName, c.Name()), self))
	}

	return self
}

func (self *ForeignKey) ForeignTable() *Table {
	return self.foreignTable
}

func (self *ForeignKey) Create(table *Table) error {
	var sql strings.Builder
	fmt.Fprintf(&sql, "ALTER TABLE \"%v\" ADD CONSTRAINT \"%v\" FOREIGN KEY (", table.name, self.name)

	for i, c := range self.cols {
		if i > 0 {
			sql.WriteString(", ")
		}

		sql.WriteString(c.Name())
	}
	
	fmt.Fprintf(&sql, ") REFERENCES \"%v\" (", self.foreignTable.name)

	for i, c := range self.foreignTable.primaryKey.cols {
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

