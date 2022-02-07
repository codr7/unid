package data

type Key interface {
	Def
	Rel
}

type BasicKey struct {
	BasicDef
	BasicRel
}

func NewKey(name string, cols...Col) *BasicKey {
	return new(BasicKey).Init(name, cols...)
}

func (self *BasicKey) Init(name string, cols...Col) *BasicKey {
	self.BasicDef.Init(name)
	self.AddCols(cols...)
	return self
}
