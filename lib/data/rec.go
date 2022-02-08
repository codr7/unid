package data

import (
	"github.com/jackc/pgx/v4"
	"log"
	"reflect"
)

type Rec interface {
	Exists() bool
	Table() Table
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

func GetField(rec Rec, name string) reflect.Value {
	s := reflect.ValueOf(rec)

	if !s.IsValid() {
		log.Fatal("Invalid rec: ", rec)
	}

	f := reflect.Indirect(s).FieldByName(name)

	if !f.IsValid() {
		log.Fatal("Field '%v' not found in %v", name, rec)
	}

	return f
}

func GetFieldAddr(rec Rec, name string) interface{} {
	return GetField(rec, name).Addr().Interface()
}

func GetFieldValue(rec Rec, name string) interface{} {
	return GetField(rec, name).Interface()
}

func SetFieldValue(rec Rec, name string, val interface{}) {
	f := GetField(rec, name)
	
	if !f.CanSet() {
		log.Fatal("Field '%v' not settable in %v", name, rec)
	}

	v := reflect.ValueOf(val)

	if !v.IsValid() {
		log.Fatal("Failed reflecting field '%v' in %v", name, rec)
	}

	f.Set(v)
}

func LoadFields(rec Rec, in pgx.Row) error {
	table := rec.Table()
	cols := table.Cols()
	fields := make([]interface{}, len(cols))

	for i, c := range cols {
		fields[i] = GetFieldAddr(rec, c.Name())
	}
	
	return in.Scan(fields...)
}
