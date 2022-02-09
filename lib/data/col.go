package data

import (
	//"log"
	"reflect"
	"time"
)

type Col interface {
	Field
	TableDef
	NewForeignCol(name string, key *ForeignKey) Col
	NewField() interface{}
	ValType() string
	ForeignKey() *ForeignKey
}

type BasicCol struct {
	BasicDef
	BasicField
	foreignKey *ForeignKey
}

func (self *BasicCol) Init(name string) {
	self.BasicDef.Init(name)
	self.BasicField.Init(name)
}

func (self *BasicCol) GetFieldAddr(rec Ref) interface{} {
	if self.foreignKey != nil {
		ref := self.foreignKey.GetFieldValue(rec).(Ref)
		i := self.foreignKey.colIndices[self.name]
		
		if p, ok := ref.(*RecProxy); ok {
			return p.key[i]
		}

		return self.foreignKey.foreignTable.primaryKey.cols[i].GetFieldAddr(ref)
	}

	return self.BasicField.GetFieldAddr(rec)
}

func (self *BasicCol) GetFieldValue(rec Ref) interface{} {
	if self.foreignKey != nil {
		ref := self.foreignKey.GetFieldValue(rec).(Ref)
		i := self.foreignKey.colIndices[self.name]
		
		if p, ok := ref.(*RecProxy); ok {
			return p.key[i]
		}

		return self.foreignKey.foreignTable.primaryKey.cols[i].GetFieldValue(ref)
	}

	return self.BasicField.GetFieldValue(rec)
}

func (self *BasicCol) SetFieldValue(rec Ref, val interface{}) {
	if self.foreignKey != nil {
		ref := self.foreignKey.GetFieldValue(rec).(Ref)
		i := self.foreignKey.colIndices[self.name]
		
		if p, ok := ref.(*RecProxy); ok {
			p.key[i] = reflect.ValueOf(val).Addr().Interface()
		} else {
			self.foreignKey.foreignTable.primaryKey.cols[i].SetFieldValue(ref, val)
		}
	} else {
		self.BasicField.SetFieldValue(rec, val)
	}
}

func (self *BasicCol) Create(table *Table) error {
	return nil
}

func (self *BasicCol) Drop(table *Table) error {
	return nil
}

func (self *BasicCol) ForeignKey() *ForeignKey {
	return self.foreignKey
}

type IntCol struct {
	BasicCol
}

func NewIntCol(name string) *IntCol {
	return new(IntCol).Init(name)
}

func (self *IntCol) Init(name string) *IntCol {
	self.BasicCol.Init(name)
	return self
}

func (self *IntCol) NewForeignCol(name string, key *ForeignKey) Col {
	c := NewIntCol(name)
	c.foreignKey = key
	return c
}

func (self *IntCol) NewField() interface{} {
	var v int
	return &v
}
func (self *IntCol) ValType() string {
	return "INTEGER"
}

type StringCol struct {
	BasicCol
}

func NewStringCol(name string) *StringCol {
	return new(StringCol).Init(name)
}

func (self *StringCol) Init(name string) *StringCol {
	self.BasicCol.Init(name)
	return self
}

func (self *StringCol) NewForeignCol(name string, key *ForeignKey) Col {
	c := NewStringCol(name)
	c.foreignKey = key
	return c
}

func (self *StringCol) NewField() interface{} {
	var v string
	return &v
}

func (self *StringCol) ValType() string {
	return "TEXT"
}

type TimeCol struct {
	BasicCol
}

func NewTimeCol(name string) *TimeCol {
	return new(TimeCol).Init(name)
}

func (self *TimeCol) Init(name string) *TimeCol {
	self.BasicCol.Init(name)
	return self
}

func (self *TimeCol) NewForeignCol(name string, key *ForeignKey) Col {
	c := NewTimeCol(name)
	c.foreignKey = key
	return c
}

func (self *TimeCol) NewField() interface{} {
	var v time.Time
	return &v
}

func (self *TimeCol) ValType() string {
	return "TIMESTAMP"
}
