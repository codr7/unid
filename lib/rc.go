package unid

import (
	"github.com/codr7/unid/lib/data"
	//"log"
	"time"
)

type Rc struct {
	data.BasicRec
	Name string
	CreatedAt time.Time
	CreatedBy data.Ref
}

func NewRc(cx *data.Cx) *Rc {
	return new(Rc).Init(cx, false)
}

func (self *Rc) Init(cx *data.Cx, exists bool) *Rc {
	self.BasicRec.Init(cx, exists)

	if !exists {
		self.CreatedAt = time.Now()
	}
	
	return self 
}

func (self *Rc) Table() *data.Table {
	return self.Cx().FindTable("Rcs")
}

func (self *Rc) GetCreatedBy() (*User, error) {
	if p, ok := self.CreatedBy.(*data.RecProxy); ok {
		out := NewUser(self.Cx())
		
		if _, err := p.Load(out); err != nil {
			return nil, err
		}

		return out, nil
	}

	return self.CreatedBy.(*User), nil
}
