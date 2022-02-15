package db

import (
	"fmt"
	"strings"
)

type Enum struct {
	BasicDef
	cx *Cx
	alts []string
}

func NewEnum(cx *Cx, name string, alts...string) *Enum {
	return new(Enum).Init(cx, name, alts...)
}

func (self *Enum) Init(cx *Cx, name string, alts...string) *Enum {
	self.BasicDef.Init(name)
	self.cx = cx
	self.alts = alts
	return self
}

func (self *Enum) Create() error {
	var sql strings.Builder
	fmt.Fprintf(&sql, "CREATE TYPE \"%v\" AS ENUM (", self.name)

	for i, a := range self.alts {
		if i > 0 {
			sql.WriteString(", ")
		}
		
		fmt.Fprintf(&sql, "'%v'", a)
	}
	
	sql.WriteRune(')')

	if err := self.cx.ExecSQL(sql.String()); err != nil {
		return err
	}
		
	return nil
}

func (self *Enum) Exists() (bool, error) {
	sql := "SELECT EXISTS (SELECT FROM pg_type WHERE typname  = $1)"
	row := self.cx.QueryRow(sql, self.name)
	var ok bool
	
	if err := row.Scan(&ok); err != nil {
		return false, err
	}
	
	return ok, nil
}

func (self *Enum) Drop() error {
	sql := fmt.Sprintf("DROP TYPE IF EXISTS \"%v\"", self.name)
	return self.cx.ExecSQL(sql)
}



