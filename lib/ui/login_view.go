package ui

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/codr7/unid/lib"
	"github.com/codr7/unid/lib/db"
	"github.com/codr7/unid/lib/dom"
	"github.com/jackc/pgx/v4"
	"log"
	"net/http"
)

func LoginView(w http.ResponseWriter, r *http.Request) {
	title := "Login"
	
	d := dom.NewDoc(title)
	d.Style("styles/reset.css")
	d.Style("styles/site.css")
	d.Style("styles/login.css")
	d.Script("scripts/site.js")
	d.Script("scripts/login.js")
	
	fs := d.Body.FieldSet("frame")
	fs.H1(title)
	fs.Label("User")
	fs.Br().Input("user", "text").Autofocus()
	fs.Br().Label("Password")
	fs.Br().Input("password", "password")
	b := fs.Br().Div("buttons").Button("enterButton", "Enter")
	b.OnClick("login();")
	
	if err := d.Write(w); err != nil {
		log.Fatal(err)
	}
}

type LoginDb struct {
	User string
	Password string
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var in LoginDb
	
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		http.Error(w, fmt.Sprintf("Failed decoding json: %v", err), http.StatusBadRequest)
		return
	}

	url := "postgres://test:test@localhost:5432/test"
	dbc, err := pgx.Connect(context.Background(), url)
	
	if err != nil {
		http.Error(w,
			fmt.Sprintf("Failed connecting to dbbase: %v", err),
			http.StatusInternalServerError)

		return
	}
	
	cx := db.NewCx(dbc)
	
	if err := unid.InitDb(cx); err != nil {
		http.Error(w,
			fmt.Sprintf("Failed initializing dbbase: %v", err),
			http.StatusInternalServerError)

		return
	}
	
	if err := cx.SyncAll(); err != nil {
		http.Error(w,
			fmt.Sprintf("Failed syncing dbbase: %v", err),
			http.StatusInternalServerError)

		return
	}

	u := unid.NewUser(cx)
	u.Name = in.User

	if err := cx.FindTable("Users").Load(u); err != nil {
		http.Error(w, fmt.Sprintf("Failed loading user: %v", u.Name), http.StatusInternalServerError)
		return
	}	

	StartSession(u, w)
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	var session *Session
	
	if session = CurrentSession(w, r); session == nil {
		http.Error(w, "Not logged in!", http.StatusBadRequest)
	}

	session.End()
}
