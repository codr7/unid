package data

type Def interface {
	Name() string
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
