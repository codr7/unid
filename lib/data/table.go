package data

import (
	"fmt"
	//"log"
	"strings"
)

type Table struct {
	BasicDef
	BasicRel

	cx *Cx
	primaryKey *Key
	foreignKeys []*ForeignKey
}

func (self *Table) Init(cx *Cx, name string, keyCols...Col) *Table {
	self.BasicDef.Init(name)
	self.BasicRel.Init()
	self.cx = cx
	self.primaryKey = NewKey(fmt.Sprintf("%vPrimaryKey", name), keyCols...)
	self.AddCol(keyCols...)
	return self
}

func (self *Table) Cx() *Cx {
	return self.cx
}

func (self *Table) PrimaryKey() *Key {
	return self.primaryKey
}

func (self *Table) NewForeignKey(name string, foreignTable *Table) *ForeignKey {
	k := new(ForeignKey).Init(fmt.Sprintf("%v%vKey", self.name, name), name, foreignTable)
	self.foreignKeys = append(self.foreignKeys, k)
	self.AddCol(k.cols...)
	return k
}

func (self *Table) Create() error {
	var sql strings.Builder
	fmt.Fprintf(&sql, "CREATE TABLE \"%v\" (", self.name)

	for i, c := range self.cols {
		if i > 0 {
			sql.WriteString(", ")
		}
		
		fmt.Fprintf(&sql, "\"%v\" %v NOT NULL", c.Name(), c.ValType())
	}
	
	sql.WriteRune(')')

	if err := self.cx.ExecSQL(sql.String()); err != nil {
		return err
	}
		
	if err := self.primaryKey.Create(self); err != nil {
		return err
	}

	for _, k := range self.foreignKeys {
		if err := k.Create(self); err != nil {
			return err
		}
	}

	return nil
}

func (self *Table) Exists() (bool, error) {
	sql := "SELECT EXISTS (SELECT FROM pg_tables WHERE tablename  = $1)"
	row := self.cx.QueryRow(sql, self.name)
	var ok bool
	
	if err := row.Scan(&ok); err != nil {
		return false, err
	}
	
	return ok, nil
}

func (self *Table) Drop() error {
	for _, k := range self.foreignKeys {
		if err := k.Drop(self); err != nil {
			return err
		}
	}

	sql := fmt.Sprintf("DROP TABLE IF EXISTS \"%v\"", self.name)

	if err := self.cx.ExecSQL(sql); err != nil {
		return err
	}

	return nil
}

func (self *Table) Insert(rec Rec) error {
	var sql strings.Builder
	fmt.Fprintf(&sql, "INSERT INTO \"%v\" (", self.name)
	var params []interface{}
	
	for i, c := range self.cols {
		if i > 0 {
			sql.WriteString(", ")
		}

		fmt.Fprintf(&sql, "\"%v\"", c.Name())
	}
	
	sql.WriteString(") VALUES (")

	for i, c := range self.cols {
		if i > 0 {
			sql.WriteString(", ")
		}
		
		fmt.Fprintf(&sql, "$%v", i+1)
		params = append(params, c.GetFieldValue(rec))
	}

	sql.WriteRune(')')

	if err := self.cx.ExecSQL(sql.String(), params...); err != nil {
		return err
	}

	return nil
}

func (self *Table) Update(rec Rec) error {
	var sql strings.Builder
	fmt.Fprintf(&sql, "UPDATE \"%v\" SET ", self.name)
	var params []interface{}
	
	for i, c := range self.cols {
		if self.primaryKey.FindCol(c.Name()) != nil {
			continue
		}
		
		if i > 0 {
			sql.WriteString(", ")
		}

		params = append(params, c.GetFieldValue(rec))
		fmt.Fprintf(&sql, "\"%v\" = $%v", c.Name(), len(params))
	}
	
	sql.WriteString(" WHERE ")

	for i, c := range self.primaryKey.Cols() {
		if i > 0 {
			sql.WriteString(" AND ")
		}
		
		params = append(params, c.GetFieldValue(rec))
		fmt.Fprintf(&sql, "\"%v\" = $%v", c.Name(), len(params))
	}

	if err := self.cx.ExecSQL(sql.String(), params...); err != nil {
		return err
	}

	return nil
}

func (self *Table) Load(rec Rec) error {
	var sql strings.Builder
	sql.WriteString("SELECT ")
	
	for i, c := range self.cols {
		if i > 0 {
			sql.WriteString(", ")
		}

		fmt.Fprintf(&sql, "\"%v\"", c.Name())
	}
	
	fmt.Fprintf(&sql, " FROM \"%v\" WHERE ", self.name)
	var params []interface{}

	for i, c := range self.primaryKey.Cols() {
		if i > 0 {
			sql.WriteString(" AND ")
		}
		
		params = append(params, c.GetFieldValue(rec))
		fmt.Fprintf(&sql, "\"%v\" = $%v", c.Name(), len(params))
	}

	row := self.cx.QueryRow(sql.String(), params...)
	var fs []interface{}

	for _, k := range self.foreignKeys {
		if v := k.GetFieldValue(rec); v == nil {
			k.SetFieldValue(rec, NewRecProxy(k.foreignTable))
		}
	}
	
	for _, c := range self.cols {
		fs = append(fs, c.GetFieldAddr(rec))
	}

	return row.Scan(fs...)
}
