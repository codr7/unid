package db

import (
	"fmt"
	"io"
	"strings"
)

type Val interface {
	WriteValSql(out io.Writer) error
	ValParams() []interface{}
}

type CustomVal struct {
	sql string
	params []interface{}
}

func NewVal(sql string, params...interface{}) Val {
	return new(CustomVal).Init(sql, params...)
}

func (self *CustomVal) Init(sql string, params...interface{}) *CustomVal {
	self.sql = sql
	self.params = params
	return self
}

func (self *CustomVal) WriteValSql(out io.Writer) error {
	_, err :=fmt.Fprintf(out, self.sql)
	return err
}

func (self *CustomVal) ValParams() []interface{} {
	return self.params
}

func Op(op string, l Val, r interface{}) Cond {
	var sql strings.Builder
	l.WriteValSql(&sql)
	return NewCond("%v %v $$", []interface{}{sql.String(), op}, []interface{}{r})
}

func Eq(l Val, r interface{}) Cond {
	return Op("=", l, r)
}

func Lt(l Val, r interface{}) Cond {
	return Op("<", l, r)
}

func Gt(l Val, r interface{}) Cond {
	return Op(">", l, r)
}
