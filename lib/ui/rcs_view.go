package ui

import (
	"fmt"
	"github.com/codr7/unid/lib"
	"github.com/codr7/unid/lib/db"
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
	tr.Th()
	tr.Th().Printf("Name")
	tr.Th().Printf("Cap")
	tr.Th().Printf("Created")
	tr.Th().Printf("by")
	rcs := cx.FindTable("Rcs")
	q := rcs.Query().OrderBy(rcs.FindCol("Name"))
	
	if err := q.Run(); err != nil {
		http.Error(w,
			fmt.Sprintf("Failed querying resources: %v", err),
			http.StatusInternalServerError)

		return
	}

	defer q.Close()
	
	for q.Next() {
		rc := unid.NewRc(cx)

		if err := db.Load(rc, q); err != nil {
			http.Error(w,
				fmt.Sprintf("Failed loading resource: %v", err),
				http.StatusInternalServerError)
		}

		tr = t.Tr()
		tr.Td().A(fmt.Sprintf("rc.html?mode=view&name=%v", rc.Name), "...")
		tr.Td().Printf(rc.Name)
		tr.Td().Printf(rc.CapType)
		tr.Td().Printf(rc.CreatedAt.Format(session.TimeFormat()))
		createdBy := rc.CreatedBy.(*db.RecProxy).KeyVals()[0].(string)
		tr.Td().A(fmt.Sprintf("user.html?mode=view?name=%v", createdBy), createdBy)
	}


	bs := fs.Br().Div("buttons")
	bs.Span().Set("class", "shortcut").Printf("Alt+N")
	bs.Br()
	b := bs.Button("newButton", "New Resource")
	b.OnClick("window.location = 'rc.html?mode=new';")
	b.Set("accesskey", "n")

	if err := d.Write(w); err != nil {
		log.Fatal(err)
	}
}
