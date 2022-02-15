package unid

import (
	"github.com/codr7/unid/lib/db"
)

type Pool struct {
	db.BasicRec
	Parent, Child db.Ref
}

func NewPool(parent, child *Rc) *Pool {
	cx := parent.Cx()
	self := new(Pool).Init(cx)
	self.Parent = parent
	self.Child = child
	return self
}

func (self *Pool) Init(cx *db.Cx) *Pool {
	self.BasicRec.Init(cx)
	return self 
}

func (self *Pool) Table() db.Table {
	return self.Cx().FindTable("Pools")
}
