package dom

import (
	"fmt"
)

func (self *Node) OnChange(spec string, args...interface{}) *Node {
	fmt.Fprintf(&self.script,
		"document.getElementById('%v').addEventListener('change', (event) => {\n%v\n});\n",
		self.Id(),
		fmt.Sprintf(spec, args...))

	return self
}

func (self *Node) OnClick(spec string, args...interface{}) *Node {
	fmt.Fprintf(&self.script,
		"document.getElementById('%v').addEventListener('click', (event) => {\n%v\n});\n",
		self.Id(),
		fmt.Sprintf(spec, args...))

	return self
}
