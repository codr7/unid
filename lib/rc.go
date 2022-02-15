package unid

import (
	"github.com/codr7/unid/lib/db"
	//"log"
	"time"
)

const (
	RcCapTypeFree = "free"
	RcCapTypePool = "pool"
	RcCapTypeUnit = "unit"
)

type Rc struct {
	db.BasicRec
	Name string
	CreatedAt time.Time
	CreatedBy db.Ref
	CapType string
}

func NewRc(cx *db.Cx) *Rc {
	self := new(Rc).Init(cx)
	self.CreatedAt = time.Now()
	return self
}

func (self *Rc) Init(cx *db.Cx) *Rc {
	self.BasicRec.Init(cx)
	return self 
}

func (self *Rc) Table() db.Table {
	return self.Cx().FindTable("Rcs")
}

func (self *Rc) AfterInsert() error {
	p := NewPool(self, self)

	if err := db.Store(p); err != nil {
		return err
	}

	c := self.NewCap(MinTime(), MaxTime(), 0, 0)

	if err := db.Store(c); err != nil {
		return err
	}

	return nil
}

func (self *Rc) GetCreatedBy() (*User, error) {
	if p, ok := self.CreatedBy.(*db.RecProxy); ok {
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

func (self *Rc) CapsQuery(startsAt, endsAt time.Time) *db.Query {
	cx := self.Cx()
	caps := cx.FindTable("Caps")	
	pools := cx.FindTable("Pools")
	
	return caps.Query().
		Join2(pools, pools.FindForeignKey("Parent"), caps.FindForeignKey("Rc")).
		Where(pools.FindCol("ChildName").Eq(self.Name),
			caps.FindCol("StartsAt").Lt(endsAt),
			caps.FindCol("EndsAt").Gt(startsAt))
}

func (self *Rc) Caps(q *db.Query) ([]*Cap, error) {	
	if err := q.Run(); err != nil {
		return nil, err
	}

	defer q.Close()
	var c Cap
	c.Init(self.Cx())
	var out []*Cap
	
	for q.Next() {
		if err := db.Load(&c, q); err != nil {
			return nil, err
		}

		out = append(out, &c)
	}
	
	return out, nil
}

func (self *Rc) UpdateCaps(startsAt, endsAt time.Time, total, used int) error {
	cx := self.Cx()
	q := self.CapsQuery(startsAt, endsAt)
	
	if used == 0 {
		rcs := cx.FindTable("Rcs")
		
		q.Join(rcs, cx.FindTable("Caps").FindForeignKey("Rc")).
			Where(rcs.FindCol("CapType").Eq(RcCapTypePool))
			 
	}
	
	cs, err := self.Caps(q)

	if err != nil {
		return err
	}
	
	cs = UpdateCaps(cs, startsAt, endsAt, total, used)

	for _, c := range cs {
		if err = db.Store(c); err != nil {
			return err
		}
	}

	return nil
}

func (self *Rc) AddPool(parent *Rc) error {
	return db.Store(NewPool(parent, self))
}
