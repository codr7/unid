package data

type Rel interface {
	Cols() []Col
}

type BasicRel struct {
	cols []Col
}

func (self *BasicRel) Cols() []Col {
	return self.cols
}

func (self *BasicRel) AddCols(cols...Col) {
	self.cols = append(self.cols, cols...)
}

