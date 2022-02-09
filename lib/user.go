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
	return new(User).Init(cx, false)
}

func (self *User) Init(cx *data.Cx, exists bool) *User {
	self.BasicRec.Init(cx, exists)

	if !exists {
		self.CreatedAt = time.Now()
	}

	return self 
}

func (self *User) Table() *data.Table {
	return self.Cx().FindTable("Users")
}
