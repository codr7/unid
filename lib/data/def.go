package data

type Def interface {
	Name() string
}

type RootDef interface {
	Def
	Create(cx *Cx) error
	Exists(cx *Cx) (bool, error)
	Drop(cx *Cx) error
}

type TableDef interface {
	Def
	Create(cx *Cx, table *Table) error
	Drop(cx *Cx, table *Table) error
}

type BasicDef struct {
	name string
}

func (self *BasicDef) Init(name string) {
	self.name = name
}

func (self *BasicDef) Name() string {
	return self.name
}
