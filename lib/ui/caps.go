package ui

import (
	"fmt"
	"github.com/codr7/unid/lib"
	"github.com/codr7/unid/lib/dom"
	"log"
	"math"
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
	
	doc := dom.NewDoc(title)
	doc.Style("styles/reset.css")
	doc.Style("styles/site.css")
	doc.Style("styles/caps.css")

	fs := doc.Body.FieldSet("frame")
	fs.H1(title)
	

	now := time.Now()
	startTime := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	slotCount := 30
	slotMax := math.MaxInt32
	endTime := startTime.AddDate(0, 0, slotCount)
	
	capsTable := fs.Table("capsTable")
	tr := capsTable.Tr()
	tr.Th().Printf(startTime.Format(session.DateFormat()))
	t := startTime
	
	
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

	totals := make([]int, slotCount)

	for _, rc := range rcs {
		tr = capsTable.Tr()
		tr.Td().A(fmt.Sprintf("rc.html?mode=show&name=%v", url.QueryEscape(rc.Name)), rc.Name)

		slots := make([]int, slotCount)

		for i, _ := range slots {
			slots[i] = slotMax
		}
		
		q := rc.CapsQuery(startTime, endTime)
		caps, err := unid.LoadCaps(q)

		if err != nil {
			http.Error(w,
				fmt.Sprintf("Failed loading capacity: %v", err),
				http.StatusInternalServerError)
			return
		}

		for _, c := range caps {
			st := startTime

			if c.StartsAt.After(st) {
				st = c.StartsAt
			}

			et := endTime

			if c.EndsAt.Before(et) {
				et = c.EndsAt
			}

			for i := int(st.Sub(startTime).Hours()) / 24; i < int(et.Sub(startTime).Hours()) / 24; i++ {
				v := c.Total - c.Used
				fmt.Printf("%v v: %v %v\n", rc.Name, v, slots[i])
				
				if v < slots[i] {
					slots[i] = v
				}
			}
		}
		
		for i, s := range slots {
			td := tr.Td()

			if s < slotMax {
				td.Printf("%v", s)
				totals[i] += s
			}
		}
	}

	tr = capsTable.Tr()
	tr.Th().Printf("Total")
	
	for _, t := range totals {
		tr.Td().Printf("%v", t)
	}
	
	bs := fs.Br().Div("buttons")

	d := bs.Div("")
	d.Br().Input("changeTotal", "number").
		Set("value", "1").
		Set("size", "2")

	d = bs.Div("")
	d.Span().Set("class", "shortcut").Printf("Alt+C")
	b := d.Br().Button("changeButton", "Change Total")
	b.Set("accesskey", "c")
	b.OnClick("changeTotal();")

	d = bs.Div("")
	d.Span().Set("class", "shortcut").Printf("Alt+R")
	b = d.Br().Button("newRvnButton", "New Reservation")
	b.OnClick("window.location = 'rvn.html?mode=new';")
	b.Set("accesskey", "r")

	d = bs.Div("")
	d.Span().Set("class", "shortcut").Printf("Alt+I")
	b = d.Br().Button("newRvnItemButton", "New Item")
	b.OnClick("window.location = 'rvn.html?mode=edit';")
	b.Set("accesskey", "i")

	if err := doc.Write(w); err != nil {
		log.Fatal(err)
	}
}
