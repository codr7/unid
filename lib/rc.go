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
	rc := new(Rc).Init(cx)
	rc.CreatedAt = time.Now()
	return rc
}

func (self *Rc) Init(cx *data.Cx) *Rc {
	self.BasicRec.Init(cx)
	return self 
}

func (self *Rc) Table() *data.Table {
	return self.Cx().FindTable("Rcs")
}

func (self *Rc) AfterInsert() error {
	c := NewCap(self.Cx(), self, MinTime(), MaxTime(), 0, 0)
	return data.Store(c)
}

func (self *Rc) GetCreatedBy() (*User, error) {
	if p, ok := self.CreatedBy.(*data.RecProxy); ok {
		out := new(User).Init(self.Cx())
		
		if _, err := p.Load(out); err != nil {
			return nil, err
		}

		return out, nil
	}

	return self.CreatedBy.(*User), nil
}

func (self *Rc) GetCaps(startsAt, endsAt time.Time) ([]*Cap, error) {
	return nil, nil
}
