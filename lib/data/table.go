package data

import (
	"fmt"
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
	self.AddCols(keyCols...)
	return self
}

func (self *Table) Cx() *Cx {
	return self.cx
}

func (self *Table) PrimaryKey() *Key {
	return self.primaryKey
}

func (self *Table) NewForeignKey(name string, foreignTable *Table) *ForeignKey {
	k := new(ForeignKey).Init(fmt.Sprintf("%v%vKey", self.name, name), foreignTable)
	self.foreignKeys = append(self.foreignKeys, k)
	self.lookup[name] = k
	return k
}

func (self *Table) Create() error {
	var sql strings.Builder
	fmt.Fprintf(&sql, "CREATE TABLE %v (", self.name)

	for i, c := range self.cols {
		if i > 0 {
			sql.WriteString(", ")
		}
		
		fmt.Fprintf(&sql, "%v %v NOT NULL", c.Name(), c.ValType())
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
	//TODO pg_tables
	return false, nil
}

func (self *Table) Drop() error {
	for _, k := range self.foreignKeys {
		if err := k.Drop(self); err != nil {
			return err
		}
	}

	sql := fmt.Sprintf("DROP TABLE IF EXISTS %v", self.name)

	if err := self.cx.ExecSQL(sql); err != nil {
		return err
	}

	return nil
}

func (self *Table) Insert(rec Rec) error {
	var sql strings.Builder
	fmt.Fprintf(&sql, "INSERT INTO %v (", self.name)
	var params []interface{}
	
	for i, c := range self.cols {
		if i > 0 {
			sql.WriteString(", ")
		}
		
		sql.WriteString(c.Name())
	}
	
	sql.WriteString(") VALUES (")

	for i, c := range self.cols {
		if i > 0 {
			sql.WriteString(", ")
		}
		
		fmt.Fprintf(&sql, "$%v", i+1)
		params = append(params, GetFieldValue(rec, c.Name()))
	}

	sql.WriteRune(')')

	if err := self.cx.ExecSQL(sql.String()); err != nil {
		return err
	}

	return nil
}

func (self *Table) Update(rec Rec) error {
	var sql strings.Builder
	fmt.Fprintf(&sql, "UPDATE %v SET ", self.name)
	var params []interface{}
	
	for i, c := range self.cols {
		if self.primaryKey.FindCol(c.Name()) != nil {
			continue
		}
		
		if i > 0 {
			sql.WriteString(", ")
		}

		cn := c.Name()
		fmt.Fprintf(&sql, "%v = $%v", cn, len(params)+1)
		params = append(params, GetFieldValue(rec, cn))
	}
	
	sql.WriteString(" WHERE ")

	for i, c := range self.primaryKey.Cols() {
		if i > 0 {
			sql.WriteString(" AND ")
		}
		
		fmt.Fprintf(&sql, "%v = $%v", c.Name(), len(params)+1)
		params = append(params, GetFieldValue(rec, c.Name()))
	}

	if err := self.cx.ExecSQL(sql.String()); err != nil {
		return err
	}

	return nil
}
