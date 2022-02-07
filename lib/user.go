package unid

import (
	"github.com/codr7/unid/lib/data"
)

type User struct {
	BasicRec
	Name string
}

func (self *User) Table(cx *data.Cx) data.Table {
	return cx.FindTable("Users")
}
