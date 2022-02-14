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

func (self *Rc) Table() data.Table {
	return self.Cx().FindTable("Rcs")
}

func (self *Rc) AfterInsert() error {
	c := self.NewCap(MinTime(), MaxTime(), 0, 0)
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

func (self *Rc) NewCap(startsAt, endsAt time.Time, total, used int) *Cap {
	c := new(Cap).Init(self.Cx())
	c.Rc = self
	c.StartsAt = startsAt
	c.EndsAt = endsAt
	c.Total = total
	c.Used = used
	c.ChangedAt = time.Now()
	return c
}

func (self *Rc) Caps(startsAt, endsAt time.Time) ([]*Cap, error) {
	cx := self.Cx()
	caps := cx.FindTable("Caps")
	q := caps.Query().
		Where(data.Eq(caps.FindCol("RcName"), self.Name),
			data.Lt(caps.FindCol("StartsAt"), endsAt),
			data.Gt(caps.FindCol("EndsAt"), startsAt))
		
	if err := q.Run(); err != nil {
		return nil, err
	}

	defer q.Close()
	var c Cap
	c.Init(self.Cx())
	var out []*Cap
	
	for q.Next() {
		if err := data.Load(&c, q); err != nil {
			return nil, err
		}

		out = append(out, &c)
	}
	
	return out, nil
}

func (self *Rc) UpdateCaps(startsAt, endsAt time.Time, total, used int) error {
	cs, err := self.Caps(startsAt, endsAt)

	if err != nil {
		return err
	}
	
	cs = UpdateCaps(cs, self, startsAt, endsAt, total, used)

	for _, c := range cs {
		if err = data.Store(c); err != nil {
			return err
		}
	}

	return nil
}
