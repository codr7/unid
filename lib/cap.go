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

func (self *Cap) Init(cx *data.Cx) *Cap {
	self.BasicRec.Init(cx)
	return self 
}

func (self *Cap) Table() data.Table {
	return self.Cx().FindTable("Caps")
}

func (self *Cap) GetRc() (*Rc, error) {
	if p, ok := self.Rc.(*data.RecProxy); ok {
		out := new(Rc).Init(self.Cx())
		
		if _, err := p.Load(out); err != nil {
			return nil, err
		}

		return out, nil
	}

	return self.Rc.(*Rc), nil
}

func UpdateCaps(in []*Cap, rc *Rc, startsAt, endsAt time.Time, total, used int) []*Cap {
	var out []*Cap
	
	for _, c := range in {
		if c.StartsAt.Before(startsAt) {
			prefix := c
			c = rc.NewCap(startsAt, c.EndsAt, c.Total, c.Used)
			prefix.EndsAt = startsAt
			out = append(out, prefix)
		} else {
			c.ChangedAt = time.Now()
		}

		c.Total += total
		c.Used += used
		out = append(out, c)

		if c.EndsAt.After(endsAt) {
			suffix := *c
			suffix.StartsAt = endsAt
			c.EndsAt = endsAt
			out = append(out, &suffix)
		}
	}

	return out
}
