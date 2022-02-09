package data

import (
	"github.com/jackc/pgx/v4"
	"reflect"
)

type Ref interface {
	Exists() bool
	Table() *Table
}

type Rec interface {
	Ref
	AfterInsert() error
	DoInsert(rec Rec) error
	DoUpdate(rec Rec) error
}

type BasicRec struct {
	cx *Cx
	exists bool
}

func (self *BasicRec) Init(cx *Cx, exists bool) {
	self.cx = cx
	self.exists = exists
}

func (self *BasicRec) Cx() *Cx {
	return self.cx
}

func (self *BasicRec) Exists() bool {
	return self.exists
}

func (self *BasicRec) AfterInsert() error {
	return nil
}

func (self *BasicRec) DoInsert(rec Rec) error {
	if err := rec.Table().Insert(rec); err != nil {
		return err
	}

	self.exists = true
	return rec.AfterInsert()
}

func (self *BasicRec) DoUpdate(rec Rec) error {
	self.exists = false
	return rec.Table().Update(rec)
}

func Load(rec Rec, row pgx.Row) error {
	table := rec.Table()
	cols := table.Cols()
	fs := make([]interface{}, len(cols))

	for i, c := range cols {
		fs[i] = c.GetFieldAddr(rec)
	}
	
	return row.Scan(fs...)
}

func Store(rec Rec) error {
	if rec.Exists() {
		return rec.DoUpdate(rec)
	}

	return rec.DoInsert(rec)
}

type RecProxy struct {
	table *Table
	keyFields []interface{}
}

func NewRecProxy(table *Table) *RecProxy {
	return new(RecProxy).Init(table)
}

func (self *RecProxy) Init(table *Table) *RecProxy {
	self.table = table
	cs := table.PrimaryKey().cols
	self.keyFields = make([]interface{}, len(cs))
	
	for i, c := range cs {
		self.keyFields[i] = c.NewField()
	}
	
	return self
}

func (self *RecProxy) KeyVals() []interface{} {
	vals := make([]interface{}, len(self.keyFields))

	for i, v := range self.keyFields {
		vals[i] = reflect.ValueOf(v).Elem().Interface()
	}

	return vals
}

func (self *RecProxy) Exists() bool {
	return true
}

func (self *RecProxy) Table() *Table {
	return self.table
}

func (self *RecProxy) Load(rec Rec) (Rec, error) {
	vs := self.KeyVals()
	
	for i, c := range self.table.PrimaryKey().cols {
		c.SetFieldValue(rec, vs[i])
	}
	
	if err := self.table.Load(rec); err != nil {
		return nil, err
	}

	return rec, nil
}
