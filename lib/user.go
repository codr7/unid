package unid

import (
	"github.com/codr7/unid/lib/db"
	"time"
)

type User struct {
	db.BasicRec
	Name string
	CreatedAt time.Time
}

func NewUser(cx *db.Cx) *User {
	u := new(User).Init(cx)
	u.CreatedAt = time.Now()
	return u
}

func (self *User) Init(cx *db.Cx) *User {
	self.BasicRec.Init(cx)
	return self 
}

func (self *User) Table() db.Table {
	return self.Cx().FindTable("Users")
}
