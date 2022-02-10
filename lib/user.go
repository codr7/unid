package unid

import (
	"github.com/codr7/unid/lib/data"
	"time"
)

type User struct {
	data.BasicRec
	Name string
	CreatedAt time.Time
}

func NewUser(cx *data.Cx) *User {
	u := new(User).Init(cx)
	u.CreatedAt = time.Now()
	return u
}

func (self *User) Init(cx *data.Cx) *User {
	self.BasicRec.Init(cx)
	return self 
}

func (self *User) Table() data.Table {
	return self.Cx().FindTable("Users")
}
