package tests

import (
	"bytes"
	"github.com/codr7/unid/lib/dom"
	"strings"
	"testing"
)

func TestDomNode(t *testing.T) {
	d := dom.NewNode("html")
	h := d.NewNode("head")
	h.NewNode("title").Printf("Foo")
	b := d.NewNode("body")
	b.NewNode("a").Set("href", "http://foo.com").Printf("Home")

	var out bytes.Buffer
	d.Write(&out)

	actual := strings.ReplaceAll(out.String(), "\n", "")
	expected := "<html><head><title>Foo</title></head><body><a href=\"http://foo.com\">Home</a></body></html>"
	
	if actual != expected {
		t.Fatalf("Unexpected result:\n%v\n%v", actual, expected)
	}
}
