package unid

import (
	"github.com/codr7/unid/lib/db"
	//"log"
	"time"
)

type Cap struct {
	db.BasicRec
	Rc db.Ref
	StartsAt, EndsAt time.Time
	Total, Used int
	ChangedAt time.Time
}

func (self *Cap) Init(cx *db.Cx) *Cap {
	self.BasicRec.Init(cx)
	return self 
}

func (self *Cap) Table() db.Table {
	return self.Cx().FindTable("Caps")
}

func (self *Cap) GetRc() (*Rc, error) {
	if p, ok := self.Rc.(*db.RecProxy); ok {
		out := new(Rc).Init(self.Cx())
		
		if _, err := p.Load(out); err != nil {
			return nil, err
		}

		return out, nil
	}

	return self.Rc.(*Rc), nil
}

func UpdateCaps(in []*Cap, startsAt, endsAt time.Time, total, used int) []*Cap {
	var out []*Cap
	
	for _, c := range in {
		if c.StartsAt.Before(startsAt) {
			prefix := c
			c := new(Cap).Init(prefix.Cx())
			c.Rc = prefix.Rc
			c.StartsAt = startsAt
			c.EndsAt = prefix.EndsAt
			c.Total = prefix.Total
			c.Used = prefix.Used
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
