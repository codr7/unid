package db

import (
	"fmt"
	"io"
	//"log"
	"reflect"
	"strings"
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
	EqCol(r Col) Cond
	Eq(r interface{}) Cond
	Lt(r interface{}) Cond
	Gt(r interface{}) Cond
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

func  (self *BasicCol) Op(op string, r interface{}) Cond {
	var sql strings.Builder
	self.WriteValSql(&sql)
	return NewCond(fmt.Sprintf("%v %v $$", sql.String(), op), r)
}

func (self *BasicCol) EqCol(r Col) Cond {
	var sql strings.Builder
	self.WriteValSql(&sql)
	sql.WriteString(" = ")
	r.WriteValSql(&sql)
	return NewCond(sql.String())
}

func (self *BasicCol) Eq(r interface{}) Cond {
	return self.Op("=", r)
}

func (self *BasicCol) Lt(r interface{}) Cond {
	return self.Op("<", r)
}

func (self *BasicCol) Gt(r interface{}) Cond {
	return self.Op(">", r)
}

type EnumCol struct {
	BasicCol
	enum *Enum
}

func (self *EnumCol) Init(name string, enum *Enum) *EnumCol {
	self.BasicCol.Init(name)
	self.enum = enum
	return self
}

func (self *EnumCol) NewForeignCol(table Table, name string, key *ForeignKey) Col {
	c := table.NewEnumCol(name, self.enum)
	c.foreignKey = key
	return c
}

func (self *EnumCol) NewField() interface{} {
	var v string
	return &v
}

func (self *EnumCol) ValType() string {
	return fmt.Sprintf("\"%v\"", self.enum.name)
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
