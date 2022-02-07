package data

import (
	"time"
)

type Col interface {
	Def
	InitField() interface{}
	ValType() string
}

type BasicCol struct {
	BasicDef
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

func (self *IntCol) InitField() interface{} {
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

func (self *StringCol) InitField() interface{} {
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

func (self *TimeCol) InitField() interface{} {
	return new(*time.Time)
}

func (self *TimeCol) ValType() string {
	return "TIMESTAMP"
}
