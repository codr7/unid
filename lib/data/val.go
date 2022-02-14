package data

import (
	"fmt"
	"io"
)

type Val interface {
	WriteValSql(out io.Writer) error
	ValParams() []interface{}
}

type BasicVal struct {
	sql string
	params []interface{}
}

func NewVal(sql string, params...interface{}) Val {
	return new(BasicVal).Init(sql, params...)
}

func (self *BasicVal) Init(sql string, params...interface{}) *BasicVal {
	self.sql = sql
	self.params = params
	return self
}

func (self *BasicVal) WriteValSql(out io.Writer) error {
	_, err :=fmt.Fprintf(out, self.sql)
	return err
}

func (self *BasicVal) ValParams() []interface{} {
	return self.params
}
