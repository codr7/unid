package dom

type Table struct {
	Node
}

type TableRow struct {
	Node
}

func (self *Node) Table(id string) *Table {
	t := new(Table)
	t.Printf("")
	self.Append(t.Init("table").Set("id", id))
	return t
}

func (self *Table) Tr() *TableRow {
	n := new(TableRow)
	self.Append(n.Init("tr"))
	return n
}

func (self *TableRow) Th() *Node {
	return self.NewNode("th")
}

func (self *TableRow) Td() *Node {
	return self.NewNode("td")
}
