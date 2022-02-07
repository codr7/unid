package data

import (
	"fmt"
)

type ForeignKey struct {
	BasicKey
	foreignTable Table
}

func (self *ForeignKey) Init(name string, foreignTable Table) *ForeignKey {
	var cols []Col
	
	for _, c := range foreignTable.PrimaryKey().Cols() {
		cols = append(cols, c.Clone(fmt.Sprintf("%v%v", name, c.Name())))
	}

	self.BasicKey.Init(name, cols...)
	self.foreignTable = foreignTable
	return self
}

func (self *ForeignKey) ForeignTable() Table {
	return self.foreignTable
}
