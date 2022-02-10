package dom

import (
	"bytes"
	"io"
)

type Doc struct {
	Node
	Head Node
	Title Node
	Body Node
}

func NewDoc(title string) *Doc {
	return new(Doc).Init(title)
}

func (self *Doc) Init(title string) *Doc {
	self.Node.Init("html")
	self.Append(self.Head.Init("head"))
	self.Head.Append(self.Title.Init("title"))
	self.Title.Printf(title)
	self.Append(self.Body.Init("body"))
	return self
}

func (self *Doc) Style(href string) *Node {
	return self.Head.NewNode("link").Set("rel", "stylesheet").Set("href", href)
}

func (self *Doc) Script(src string) *Node {
	return self.Head.NewNode("script").Set("src", src).Printf("")
}

func (self *Doc) Write(out io.Writer) error {
	var js bytes.Buffer
	self.WriteScript(&js)
	
	if js.Len() > 0 {
		self.Head.NewNode("script").Printf(
			"document.addEventListener('DOMContentLoaded', (event) => {\n%v\n});",
			js.String())
	}
	
	io.WriteString(out, "<!DOCTYPE html>\n")
	self.Node.Write(out)
	return nil
}
