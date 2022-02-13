package data

import (
	"fmt"
	"io"
)

type Cond interface {
	WriteCondSql(out io.Writer) error
	CondParams() []interface{}
}

type BasicCond struct {
	sql string
	params []interface{}
}

func NewCond(sql string, params...interface{}) Cond {
	return new(BasicCond).Init(sql, params...)
}

func (self *BasicCond) Init(sql string, params...interface{}) *BasicCond {
	self.sql = sql
	self.params = params
	return self
}

func (self *BasicCond) WriteCondSql(out io.Writer) error {
	_, err :=fmt.Fprintf(out, self.sql)
	return err
}

func (self *BasicCond) CondParams() []interface{} {
	return self.params
}
