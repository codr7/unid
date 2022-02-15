package dom

var br Node

func init() {
	br.Init("br")
}

func (self *Node) A(href string, caption string) *Node {
	return self.NewNode("a").Set("href", href).Printf(caption)
}

func (self *Node) Br() *Node {
	self.Append(&br)
	return self
}

func (self *Node) Button(id string, caption string) *Node {
	return self.NewNode("button").Set("id", id).Printf(caption)
}

func (self *Node) Div(id string) *Node {
	return self.NewNode("div").Set("id", id)
}

func (self *Node) FieldSet(id string) *Node {
	return self.NewNode("fieldset").Set("id", id)
}

func (self *Node) H1(caption string) *Node {
	self.NewNode("h1").Printf(caption)
	return self
}

func (self *Node) H2(caption string) *Node {
	self.NewNode("h2").Printf(caption)
	return self
}

func (self *Node) Input(id string, inputType string) *Node {
	return self.NewNode("input").
		Set("id", id).
		Set("type", inputType)
}

func (self *Node) Label(caption string) *Node {
	return self.NewNode("label").Printf(caption)
}

func (self *Node) Span() *Node {
	return self.NewNode("span")
}

func (self *Node) Ul(id string) *Node {
	return self.NewNode("ul").Set("id", id).Printf("")
}
