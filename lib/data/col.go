package data

import (
	"time"
)

type Col interface {
	TableDef
	Clone(name string) Col
	NewField() interface{}
	ValType() string
}

type BasicCol struct {
	BasicDef
}

func (self *BasicCol) Create(table *Table) error {
	return nil
}

func (self *BasicCol) Drop(table *Table) error {
	return nil
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

func (self *IntCol) Clone(name string) Col {
	return NewIntCol(name)
}

func (self *IntCol) NewField() interface{} {
	return new(*int)
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

func (self *StringCol) Clone(name string) Col {
	return NewStringCol(name)
}

func (self *StringCol) NewField() interface{} {
	return new(*string)
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

func (self *TimeCol) Clone(name string) Col {
	return NewTimeCol(name)
}

func (self *TimeCol) NewField() interface{} {
	return new(*time.Time)
}

func (self *TimeCol) ValType() string {
	return "TIMESTAMP"
}
