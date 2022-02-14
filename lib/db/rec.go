package db

import (
	"reflect"
)

type Ref interface {
	Table() Table
}

type Rec interface {
	Ref
	AfterInsert() error
	DoInsert(rec Rec) error
	DoUpdate(rec Rec) error
}

type BasicRec struct {
	cx *Cx
}

func (self *BasicRec) Init(cx *Cx) {
	self.cx = cx
}

func (self *BasicRec) Cx() *Cx {
	return self.cx
}

func (self *BasicRec) AfterInsert() error {
	return nil
}

func (self *BasicRec) DoInsert(rec Rec) error {
	if err := rec.Table().Insert(rec); err != nil {
		return err
	}

	return rec.AfterInsert()
}

func (self *BasicRec) DoUpdate(rec Rec) error {
	return rec.Table().Update(rec)
}

func Load(rec Rec, src Source) error {
	return rec.Table().LoadFields(rec, src)
}

func Store(rec Rec) error {
	if srec := rec.Table().StoredRec(rec); srec != nil {
		return rec.DoUpdate(rec)
	}

	return rec.DoInsert(rec)
}

type RecProxy struct {
	table Table
	keyFields []interface{}
}

func NewRecProxy(table Table) *RecProxy {
	return new(RecProxy).Init(table)
}

func (self *RecProxy) Init(table Table) *RecProxy {
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

func (self *RecProxy) Table() Table {
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
