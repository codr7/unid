package db

import (
	"log"
	"reflect"
)

type Field interface {
	FieldName() string
	GetField(rec Ref) reflect.Value
	GetFieldAddr(rec Ref) interface{}
	GetFieldValue(rec Ref) interface{}
	SetFieldValue(rec Ref, val interface{})
}

type BasicField struct {
	fieldName string
}

func (self *BasicField) Init(fieldName string) {
	self.fieldName = fieldName
}

func (self *BasicField) FieldName() string {
	return self.fieldName
}

func (self *BasicField) GetField(rec Ref) reflect.Value {
	s := reflect.ValueOf(rec)

	if !s.IsValid() {
		log.Fatalf("Invalid rec: ", rec)
	}

	f := reflect.Indirect(s).FieldByName(self.fieldName)

	if !f.IsValid() {
		log.Fatalf("Field '%v' not found in %v", self.fieldName, rec)
	}

	return f
}

func (self *BasicField) GetFieldAddr(rec Ref) interface{} {
	return self.GetField(rec).Addr().Interface()
}

func (self *BasicField) GetFieldValue(rec Ref) interface{} {
	return self.GetField(rec).Interface()
}

func (self *BasicField) SetFieldValue(rec Ref, val interface{}) {
	f := self.GetField(rec)
	
	if !f.CanSet() {
		log.Fatalf("Field '%v' not settable in %v", self.fieldName, rec)
	}

	v := reflect.ValueOf(val)

	if !v.IsValid() {
		log.Fatalf("Failed reflecting val: %v", val)
	}

	f.Set(v)
}


