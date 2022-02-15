package dom

func (self *Node) Autofocus() *Node {
	return self.Set("autofocus", nil)
}

func (self *Node) Readonly() *Node {
	return self.Set("readonly", nil)
}
