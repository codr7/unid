package unid

import (
	"github.com/codr7/unid/lib/data"
	//"log"
	"time"
)

type Cap struct {
	data.BasicRec
	Rc data.Ref
	StartsAt, EndsAt time.Time
	Total, Used int
	ChangedAt time.Time
}

func NewCap(cx *data.Cx, rc *Rc, startsAt, endsAt time.Time, total, used int) *Cap {
	c := new(Cap).Init(cx, false)
	c.Rc = rc
	c.StartsAt = startsAt
	c.EndsAt = endsAt
	c.Total = total
	c.Used = used
	return c
}

func (self *Cap) Init(cx *data.Cx, exists bool) *Cap {
	self.BasicRec.Init(cx, exists)

	if !exists {
		self.ChangedAt = time.Now()
	}

	return self 
}

func (self *Cap) Table() *data.Table {
	return self.Cx().FindTable("Caps")
}

func (self *Cap) GetRc() (*Rc, error) {
	if p, ok := self.Rc.(*data.RecProxy); ok {
		out := new(Rc).Init(self.Cx(), true)
		
		if _, err := p.Load(out); err != nil {
			return nil, err
		}

		return out, nil
	}

	return self.Rc.(*Rc), nil
}