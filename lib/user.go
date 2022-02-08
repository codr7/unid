package unid

import (
	"github.com/codr7/unid/lib/data"
)

type User struct {
	data.BasicRec
	Name string
}

func (self *User) Init(cx *data.Cx, exists bool) *User {
	self.BasicRec.Init(cx, exists)
	return self
}

func (self *User) Table() *data.Table {
	return self.Cx().FindTable("Users")
}
