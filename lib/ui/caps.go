package ui

import (
	"fmt"
	"github.com/codr7/unid/lib/dom"
	"log"
	"net/http"
	"net/url"
	"time"
)

func CapsView(w http.ResponseWriter, r *http.Request) {
	session := CurrentSession(w, r)
	
	if session == nil {
		return
	}
	
	cx := session.Cx()
	title := "Capacity"
	
	d := dom.NewDoc(title)
	d.Style("styles/reset.css")
	d.Style("styles/site.css")

	fs := d.Body.FieldSet("frame")
	fs.H1(title)
	

	now := time.Now()
	startTime := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	
	capsTable := fs.Table("capsTable")
	tr := capsTable.Tr()
	tr.Th().Printf(startTime.Format(session.DateFormat()))
	t := startTime
	slotCount := 30
	
	for i := 0; i < slotCount; i++ {
		tr.Th().Printf("%v", t.Day())
		t = t.AddDate(0, 0, 1)
	}
	
	rcs, err := getRcs(cx)

	if err != nil {
		http.Error(w,
			fmt.Sprintf("Failed loading resources: %v", err),
			http.StatusInternalServerError)
		return
	}
	
	for _, rc := range rcs {
		tr = capsTable.Tr()
		tr.Td().A(fmt.Sprintf("rc.html?mode=show&name=%v", url.QueryEscape(rc.Name)), rc.Name)

		for i := 0; i < slotCount; i++ {
			tr.Td()
		}
	}

	bs := fs.Br().Div("buttons")
	bs.Span().Set("class", "shortcut").Printf("Alt+N")
	bs.Br()
	b := bs.Button("newButton", "New Reservation")
	b.OnClick("window.location = 'rvn.html?mode=new';")
	b.Set("accesskey", "n")

	if err := d.Write(w); err != nil {
		log.Fatal(err)
	}
}
