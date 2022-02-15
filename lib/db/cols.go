package db

type Cols interface {
	Cols() []Col
	AddCol(cols...Col)
	FindCol(name string) Col
	SetPrimaryKey(val bool)
}

type BasicCols struct {
	cols []Col
	colIndices map[string]int
}

func (self *BasicCols) Init() {
	self.colIndices = make(map[string]int)
}

func (self *BasicCols) Cols() []Col {
	return self.cols
}

func (self *BasicCols) AddCol(cols...Col) {
	for _, c := range cols {
		self.colIndices[c.Name()] = len(self.cols)
		self.cols = append(self.cols, c)
	}
}

func (self *BasicCols) FindCol(name string) Col {
	if i, ok := self.colIndices[name]; ok {
		return self.cols[i]
	}		

	return nil
}

func (self *BasicCols) SetPrimaryKey(val bool) {
	for _, c := range self.cols {
		c.SetPrimaryKey(true)
	}
}

