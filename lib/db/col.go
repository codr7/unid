package db

import (
	"fmt"
	"io"
	//"log"
	"reflect"
	"time"
)

type Col interface {
	Field
	TableDef
	Val
	
	NewForeignCol(table Table, name string, key *ForeignKey) Col
	NewField() interface{}
	ValType() string
	IsPrimaryKey() bool
	SetPrimaryKey(bool)
	ForeignKey() *ForeignKey
}

type BasicCol struct {
	BasicDef
	BasicField
	isPrimaryKey bool
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
			return p.keyFields[i]
		}

		return self.foreignKey.foreignTable.PrimaryKey().cols[i].GetFieldAddr(ref)
	}

	return self.BasicField.GetFieldAddr(rec)
}

func (self *BasicCol) GetFieldValue(rec Ref) interface{} {
	if self.foreignKey != nil {
		ref := self.foreignKey.GetFieldValue(rec).(Ref)
		i := self.foreignKey.colIndices[self.name]
		
		if p, ok := ref.(*RecProxy); ok {
			return p.keyFields[i]
		}

		return self.foreignKey.foreignTable.PrimaryKey().cols[i].GetFieldValue(ref)
	}

	return self.BasicField.GetFieldValue(rec)
}

func (self *BasicCol) SetFieldValue(rec Ref, val interface{}) {
	if self.foreignKey != nil {
		ref := self.foreignKey.GetFieldValue(rec).(Ref)
		i := self.foreignKey.colIndices[self.name]
		
		if p, ok := ref.(*RecProxy); ok {
			p.keyFields[i] = reflect.ValueOf(val).Addr().Interface()
		} else {
			self.foreignKey.foreignTable.PrimaryKey().cols[i].SetFieldValue(ref, val)
		}
	} else {
		self.BasicField.SetFieldValue(rec, val)
	}
}

func (self *BasicCol) Create(table Table) error {
	return nil
}

func (self *BasicCol) Drop(table Table) error {
	return nil
}

func (self *BasicCol) IsPrimaryKey() bool {
	return self.isPrimaryKey
}

func (self *BasicCol) SetPrimaryKey(val bool) {
	self.isPrimaryKey = val
}

func (self *BasicCol) ForeignKey() *ForeignKey {
	return self.foreignKey
}

func (self *BasicCol) WriteValSql(out io.Writer) error {
	_, err :=fmt.Fprintf(out, "\"%v\"", self.name)
	return err
}

func (self *BasicCol) ValParams() []interface{} {
	return nil
}

type IntCol struct {
	BasicCol
}

func (self *IntCol) Init(name string) *IntCol {
	self.BasicCol.Init(name)
	return self
}

func (self *IntCol) NewForeignCol(table Table, name string, key *ForeignKey) Col {
	c := table.NewIntCol(name)
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

func (self *StringCol) Init(name string) *StringCol {
	self.BasicCol.Init(name)
	return self
}

func (self *StringCol) NewForeignCol(table Table, name string, key *ForeignKey) Col {
	c := table.NewStringCol(name)
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

func (self *TimeCol) Init(name string) *TimeCol {
	self.BasicCol.Init(name)
	return self
}

func (self *TimeCol) NewForeignCol(table Table, name string, key *ForeignKey) Col {
	c := table.NewTimeCol(name)
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