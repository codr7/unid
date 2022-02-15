package ui

import (
	"fmt"
	"github.com/codr7/unid/lib"
	"github.com/codr7/unid/lib/db"
	"github.com/codr7/unid/lib/dom"
	"log"
	"net/http"
	"strings"
)

func getRcs(cx *db.Cx) ([]*unid.Rc, error) {
	rcs := cx.FindTable("Rcs")
	q := rcs.Query().OrderBy(rcs.FindCol("Name"))
	
	if err := q.Run(); err != nil {
		return nil, err
	}

	var out []*unid.Rc
	
	for q.Next() {
		rc := unid.NewRc(cx)

		if err := db.Load(rc, q); err != nil {
			return nil, err
		}

		out = append(out, rc)
	}

	q.Close()
	return out, nil
}

func getPools(rc *unid.Rc) (string, error) {
	cx := rc.Cx()
	pools := cx.FindTable("Pools")
	parentName := pools.FindCol("ParentName")
	
	q := db.NewQuery(cx).
		Select(parentName).
		From(pools).
		Where(pools.FindCol("ChildName").Eq(rc.Name)).
		OrderBy(parentName)

	if err := q.Run(); err != nil {
		return "", err
	}
	
	var out strings.Builder
	defer q.Close()
	
	for q.Next() {
		var pn string
		q.Scan(&pn)

		if out.Len() > 0 {
			out.WriteString(", ")
		}
		
		out.WriteString(pn)
	}

	return out.String(), nil
}

func RcsView(w http.ResponseWriter, r *http.Request) {
	session := CurrentSession(w, r)
	
	if session == nil {
		return
	}
	
	cx := session.Cx()
	title := "Resources"
	
	d := dom.NewDoc(title)
	d.Style("styles/reset.css")
	d.Style("styles/site.css")

	fs := d.Body.FieldSet("frame")
	fs.H1(title)
	
	t := fs.Table("rcsTable")
	tr := t.Tr()
	tr.Th()
	tr.Th().Printf("Name")
	tr.Th().Printf("Cap")
	tr.Th().Printf("Pools")
	tr.Th().Printf("Created")
	tr.Th().Printf("by")
	
	rcs, err := getRcs(cx)

	if err != nil {
		http.Error(w,
			fmt.Sprintf("Failed loading resources: %v", err),
			http.StatusInternalServerError)
		return
	}
	
	for _, rc := range rcs {
		tr = t.Tr()
		tr.Td().A(fmt.Sprintf("rc.html?mode=show&name=%v", rc.Name), "...")
		tr.Td().Printf(rc.Name)
		tr.Td().Printf(strings.Title(rc.CapType))

		pools, err := getPools(rc)

		if err != nil {
			http.Error(w,
				fmt.Sprintf("Failed loading pools: %v", err),
				http.StatusInternalServerError)
			return
		}
		
		tr.Td().Printf(pools)
		tr.Td().Printf(rc.CreatedAt.Format(session.TimeFormat()))
		
		createdBy := rc.CreatedBy.(*db.RecProxy).KeyVals()[0].(string)
		tr.Td().A(fmt.Sprintf("user.html?mode=show?name=%v", createdBy), createdBy)
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
