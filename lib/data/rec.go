package data

import (
)

type Rec interface {
	Exists() bool
	Table(cx *Cx) Table
}

type BasicRec struct {
	exists bool
}

func (self *BasicRec) Exists() bool {
	return self.exists
}
