package db

import (
	"fmt"
	"io"
)

type Join struct {
	left Rel
	right Rel
	conds []Cond
}

func NewJoin(left Rel, right Rel, conds...Cond) *Join {
	return new(Join).Init(left, right, conds...)
}

func (self *Join) Init(left Rel, right Rel, conds...Cond) *Join {
	self.left = left
	self.right = right
	self.conds = conds
	return self
}

func (self *Join) WriteRelSql(sql io.Writer) error {
	if err := self.left.WriteRelSql(sql); err != nil {
		return err
	}
	
	fmt.Fprintf(sql, " JOIN ")
	
	if err := self.right.WriteRelSql(sql); err != nil {
		return err
	}

	fmt.Fprintf(sql, " ON ")

	for i, c := range self.conds {
		if i > 0 {
			fmt.Fprintf(sql, " AND ")
		}
		
		if err := c.WriteCondSql(sql); err != nil {
			return err
		}
	}

	return nil
}

func (self *Join) RelParams() []interface{} {
	var out []interface{}

	out = append(out, self.left.RelParams()...)
	out = append(out, self.right.RelParams()...)
	
	for _, c := range self.conds {
		out = append(out, c.CondParams()...)

	}
	
	return out
}

