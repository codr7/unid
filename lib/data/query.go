package data

import (
	"fmt"
	"github.com/jackc/pgx/v4"
	"strings"
)

type Query struct {
	BasicRel
	cx *Cx
	from []Rel
	conds []Cond
	rows pgx.Rows
}

func NewQuery(cx *Cx) *Query {
	return new(Query).Init(cx)
}

func (self *Query) Init(cx *Cx) *Query {
	self.cx = cx
	return self
}

func (self *Query) Select(in...Col) *Query {
	self.cols = append(self.cols, in...)
	return self
}

func (self *Query) From(rel Rel) *Query {
	self.from = append(self.from, rel)
	return self
}

func (self *Query) Where(in...Cond) *Query {
	self.conds = append(self.conds, in...)
	return self
}

func (self *Query) Run() error {
	var sql strings.Builder
	var params []interface{}
	sql.WriteString("SELECT ")
	
	for i, c := range self.cols {
		if i > 0 {
			sql.WriteString(", ")
		}

		fmt.Fprintf(&sql, "\"%v\"", c.Name())
	}

	sql.WriteString(" FROM ")
	
	for i, r := range self.from {
		if i > 0 {
			sql.WriteString(", ")
		}

		if err := r.WriteRelSql(&sql); err != nil {
			return err
		}
		
		params = append(params, r.RelParams()...)
	}

	if self.conds != nil {
		sql.WriteString(" WHERE ")
		
		for i, c := range self.conds {
			if i > 0 {
				sql.WriteString(" AND ")
			}
			
			if err := c.WriteCondSql(&sql); err != nil {
				return err
			}
			
			params = append(params, c.CondParams()...)
		}
	}

	var err error
	self.rows, err = self.cx.Query(sql.String(), params...)
	return err
}

func (self *Query) Next() bool {
	return self.rows.Next()
}

func (self *Query) Scan(dst...interface{}) error {
	return self.rows.Scan(dst...)
}
