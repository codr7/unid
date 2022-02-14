package data

import (
	"github.com/jackc/pgx/v4"
	"strings"
)

type Query struct {
	cx *Cx
	vals []Val
	from []Rel
	conds []Cond
	rows pgx.Rows
	order []Val
}

func NewQuery(cx *Cx) *Query {
	return new(Query).Init(cx)
}

func (self *Query) Init(cx *Cx) *Query {
	self.cx = cx
	return self
}

func (self *Query) Select(in...Val) *Query {
	self.vals = append(self.vals, in...)
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

func (self *Query) OrderBy(in...Val) *Query {
	self.order = append(self.order, in...)
	return self
}

func (self *Query) Run() error {
	var sql strings.Builder
	var params []interface{}
	sql.WriteString("SELECT ")
	
	for i, v := range self.vals {
		if i > 0 {
			sql.WriteString(", ")
		}

		v.WriteValSql(&sql)
		params = append(params, v.ValParams()...)
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

	if self.order != nil {
		sql.WriteString(" ORDER BY ")

		for i, v := range self.order {
			if i > 0 {
				sql.WriteString(", ")
			}

			v.WriteValSql(&sql)
			params = append(params, v.ValParams()...)
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
