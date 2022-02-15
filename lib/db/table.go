package db

import (
	"fmt"
	"io"
	"log"
	"strings"
)

type Table interface {
	Cols
	RootDef
	Rel

	Cx() *Cx
	PrimaryKey() *Key

	NewForeignKey(name string, foreignTable Table) *ForeignKey
	FindForeignKey(name string) *ForeignKey
	
	NewEnumCol(name string, enum *Enum) *EnumCol
	NewIntCol(name string) *IntCol
	NewStringCol(name string) *StringCol
	NewTimeCol(name string) *TimeCol

	NewStoredRec() StoredRec
	StoredRec(rec Rec) StoredRec
	Insert(rec Rec) error
	Update(rec Rec) error
	LoadFields(rec Rec, src Source) error
	Load(rec Rec) error
	Query() *Query
}

type BasicTable struct {
	BasicCols
	BasicDef

	cx *Cx
	primaryKey *Key
	foreignKeys map[string]*ForeignKey
	storedRecs map[Rec]StoredRec
}

type StoredRec = []interface{}

func (self *BasicTable) Init(cx *Cx, name string) *BasicTable {
	self.BasicCols.Init()
	self.BasicDef.Init(name)
	self.cx = cx
	self.foreignKeys = make(map[string]*ForeignKey)
	self.storedRecs = make(map[Rec]StoredRec)
	return self
}

func (self *BasicTable) Cx() *Cx {
	return self.cx
}

func (self *BasicTable) PrimaryKey() *Key {
	if self.primaryKey == nil {
		var cs []Col
		
		for _, c := range self.cols {
			if c.IsPrimaryKey() {
				cs = append(cs, c)
			}
		}
		
		self.primaryKey = NewKey(fmt.Sprintf("%vPrimaryKey", self.name), cs...)
	}
	
	return self.primaryKey
}

func (self *BasicTable) NewEnumCol(name string, enum *Enum) *EnumCol {
	c := new(EnumCol).Init(self, name, enum)
	self.AddCol(c)
	return c
}

func (self *BasicTable) NewIntCol(name string) *IntCol {
	c := new(IntCol).Init(self, name)
	self.AddCol(c)
	return c
}

func (self *BasicTable) NewStringCol(name string) *StringCol {
	c := new(StringCol).Init(self, name)
	self.AddCol(c)
	return c
}

func (self *BasicTable) NewTimeCol(name string) *TimeCol {
	c := new(TimeCol).Init(self, name)
	self.AddCol(c)
	return c
}

func (self *BasicTable) NewForeignKey(name string, foreignTable Table) *ForeignKey {
	k := new(ForeignKey).Init(fmt.Sprintf("%v%vKey", self.name, name), name, foreignTable)

	for _, c := range foreignTable.PrimaryKey().Cols() {
		fn := fmt.Sprintf("%v%v", name, c.Name())
		
		if self.FindCol(fn) != nil {
			log.Fatalf("Duplicate column in %v: %v", self.name, fn)
		}
		
		k.AddCol(c.NewForeignCol(self, fn, k))
	}

	self.foreignKeys[name] = k
	return k
}

func (self *BasicTable) FindForeignKey(name string) *ForeignKey {
	return self.foreignKeys[name]
}

func (self *BasicTable) Create() error {
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
		
	if err := self.PrimaryKey().Create(self); err != nil {
		return err
	}

	for _, k := range self.foreignKeys {
		if err := k.Create(self); err != nil {
			return err
		}
	}

	return nil
}

func (self *BasicTable) Exists() (bool, error) {
	sql := "SELECT EXISTS (SELECT FROM pg_tables WHERE tablename  = $1)"
	row := self.cx.QueryRow(sql, self.name)
	var ok bool
	
	if err := row.Scan(&ok); err != nil {
		return false, err
	}
	
	return ok, nil
}

func (self *BasicTable) Drop() error {
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

func (self *BasicTable) NewStoredRec() StoredRec {
	return make(StoredRec, len(self.cols))
}

func (self *BasicTable) StoredRec(rec Rec) StoredRec {
	return self.storedRecs[rec]
}

func (self *BasicTable) Insert(rec Rec) error {
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

	srec := self.NewStoredRec()

	for i, c := range self.cols {
		srec[i] = c.GetFieldValue(rec)
	}

	self.storedRecs[rec] = srec
	return nil
}

func (self *BasicTable) Update(rec Rec) error {
	var sql strings.Builder
	fmt.Fprintf(&sql, "UPDATE \"%v\" SET ", self.name)
	var params []interface{}
	i := 0
	
	for _, c := range self.cols {
		if self.PrimaryKey().FindCol(c.Name()) != nil {
			continue
		}
		
		if i > 0 {
			sql.WriteString(", ")
		}

		params = append(params, c.GetFieldValue(rec))
		fmt.Fprintf(&sql, "\"%v\" = $%v", c.Name(), len(params))
		i++
	}
	
	sql.WriteString(" WHERE ")

	for i, c := range self.PrimaryKey().Cols() {
		if i > 0 {
			sql.WriteString(" AND ")
		}
		
		params = append(params, c.GetFieldValue(rec))
		fmt.Fprintf(&sql, "\"%v\" = $%v", c.Name(), len(params))
	}

	if err := self.cx.ExecSQL(sql.String(), params...); err != nil {
		return err
	}

	srec := self.storedRecs[rec]

	for i, c := range self.cols {
		srec[i] = c.GetFieldValue(rec)
	}

	return nil
}

func (self *BasicTable) LoadFields(rec Rec, src Source) error {
	fs := make([]interface{}, len(self.cols))

	for _, k := range self.foreignKeys {
		k.SetFieldValue(rec, NewRecProxy(k.foreignTable))
	}

	for i, c := range self.cols {
		fs[i] = c.GetFieldAddr(rec)
	}
	
	if err := src.Scan(fs...); err != nil {
		return err
	}

	srec := self.storedRecs[rec]

	if srec == nil {
		srec = self.NewStoredRec()
		self.storedRecs[rec] = srec
	}

	for i, c := range self.cols {
		srec[i] = c.GetFieldValue(rec)
	}

	return nil
}

func (self *BasicTable) Load(rec Rec) error {
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

	for i, c := range self.PrimaryKey().Cols() {
		if i > 0 {
			sql.WriteString(" AND ")
		}
		
		params = append(params, c.GetFieldValue(rec))
		fmt.Fprintf(&sql, "\"%v\" = $%v", c.Name(), len(params))
	}

	src := self.cx.QueryRow(sql.String(), params...)
	return self.LoadFields(rec, src)
}

func (self *BasicTable) Query() *Query {
	q := NewQuery(self.cx).From(self)

	for _, c := range self.cols {
		q.Select(c)
	}

	return q
}

func (self *BasicTable) WriteRelSql(out io.Writer) error {
	_, err := fmt.Fprintf(out, "\"%s\"", self.name)
	return err
}

func (self *BasicTable) RelParams() []interface{} {
	return nil
}

