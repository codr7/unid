package main

import (
	"context"
	"github.com/codr7/unid/lib"
	"github.com/codr7/unid/lib/data"
	"github.com/codr7/unid/lib/ui"
	"github.com/jackc/pgx/v4"
	"log"
	"net/http"
)

func StartWorker() {
	url := "postgres://test:test@localhost:5432/test"
	dbc, err := pgx.Connect(context.Background(), url)
	
	if err != nil {
		log.Fatal(err)
	}
	
	defer dbc.Close(context.Background())
	cx := data.NewCx(dbc)
	
	if err := unid.InitDb(cx); err != nil {
		log.Fatal(err)
	}
	
	if err := cx.SyncAll(); err != nil {
		log.Fatal(err)
	}
}

func main() {
	log.Printf("Welcome to Unid v%v", unid.Version())
	http.Handle("/", http.FileServer(http.Dir("www")))
	http.HandleFunc("/rcs.html", ui.RcsView)
	log.Fatal(http.ListenAndServe(":8080", nil))	
}
