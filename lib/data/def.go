package data

type Def interface {
	Name() string
}

type RootDef interface {
	Def
	Create() error
	Exists() (bool, error)
	Drop() error
}

type TableDef interface {
	Def
	Create(table Table) error
	Drop(table Table) error
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
