package ui

import (
	"github.com/codr7/unid/lib/dom"
	"log"
	"net/http"
)

func RcsView(w http.ResponseWriter, r *http.Request) {
	/*if session := CurrentSession(w, r); session == nil {
		return
	}*/
	
	title := "Resources"
	
	d := dom.NewDoc(title)
	d.Style("css/reset.css")
	d.Style("css/site.css")

	t := d.Body.Table("resourceTable")
	t.Tr().Td().Printf("foo")
	t.Tr().Td().Printf("bar")
	t.Tr().Td().Printf("baz")
	
	if err := d.Write(w); err != nil {
		log.Fatal(err)
	}
}
