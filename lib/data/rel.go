package data

type Rel interface {
	Cols() []Col
	AddCol(cols...Col)
	FindCol(name string) Col
}

type BasicRel struct {
	cols []Col
	colIndices map[string]int
}

func (self *BasicRel) Init() {
	self.colIndices = make(map[string]int)
}

func (self *BasicRel) Cols() []Col {
	return self.cols
}

func (self *BasicRel) AddCol(cols...Col) {
	for _, c := range cols {
		self.colIndices[c.Name()] = len(self.cols)
		self.cols = append(self.cols, c)
	}
}

func (self *BasicRel) FindCol(name string) Col {
	if i, ok := self.colIndices[name]; ok {
		return self.cols[i]
	}		

	return nil
}

