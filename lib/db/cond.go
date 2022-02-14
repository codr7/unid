package db

import (
	"fmt"
	"io"
)

type Cond interface {
	WriteCondSql(out io.Writer) error
	CondParams() []interface{}
}

type CustomCond struct {
	sql string
	params []interface{}
}

func NewCond(sql string, args []interface{}, params []interface{}) Cond {
	return new(CustomCond).Init(sql, args, params)
}

func (self *CustomCond) Init(sql string, args []interface{}, params []interface{}) *CustomCond {
	self.sql = fmt.Sprintf(sql, args...)
	self.params = params
	return self
}

func (self *CustomCond) WriteCondSql(out io.Writer) error {
	_, err :=fmt.Fprintf(out, self.sql)
	return err
}

func (self *CustomCond) CondParams() []interface{} {
	return self.params
}
