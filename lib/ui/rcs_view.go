package ui

import (
	"fmt"
	"github.com/codr7/unid/lib"
	"github.com/codr7/unid/lib/data"
	"github.com/codr7/unid/lib/dom"
	"log"
	"net/http"
)

func RcsView(w http.ResponseWriter, r *http.Request) {
	session := CurrentSession(w, r)
	
	if session == nil {
		return
	}
	
	title := "Resources"
	
	d := dom.NewDoc(title)
	d.Style("styles/reset.css")
	d.Style("styles/site.css")

	fs := d.Body.FieldSet("frame")
	fs.H1(title)

	cx := session.Cx()
	t := fs.Table("rcsTable")
	tr := t.Tr()
	tr.Th().Printf("Name")
	tr.Th().Printf("Created @")
	rcs := cx.FindTable("Rcs")
	q := rcs.Query().OrderBy(rcs.FindCol("Name"))
	
	if err := q.Run(); err != nil {
		http.Error(w,
			fmt.Sprintf("Failed querying resources: %v", err),
			http.StatusInternalServerError)

		return
	}
	
	for q.Next() {
		rc := unid.NewRc(cx)

		if err := data.Load(rc, q); err != nil {
			http.Error(w,
				fmt.Sprintf("Failed loading resource: %v", err),
				http.StatusInternalServerError)
		}

		tr = t.Tr()
		tr.Td().Printf(rc.Name)
		tr.Td().Printf(rc.CreatedAt.Format(session.TimeFormat()))
	}
	
	if err := d.Write(w); err != nil {
		log.Fatal(err)
	}
}
