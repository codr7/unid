package data

type Rel interface {
	Cols() []Col
	AddCols(cols...Col)
}

type BasicRel struct {
	cols []Col
	lookup map[string]Def
}

func (self *BasicRel) Init() {
	self.lookup = make(map[string]Def)
}

func (self *BasicRel) Cols() []Col {
	return self.cols
}

func (self *BasicRel) AddCols(cols...Col) {
	self.cols = append(self.cols, cols...)

	for _, c := range cols {
		self.lookup[c.Name()] = c
	}
}

func (self *BasicRel) FindCol(name string) Col {
	return self.lookup[name].(Col)
}

