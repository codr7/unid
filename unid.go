package main

import (
	"context"
	"log"
	"github.com/jackc/pgx/v4"
	"github.com/codr7/unid/lib"
	"github.com/codr7/unid/lib/data"
)

func main() {
	url := "postgres://test:test@localhost:5432/test"
	dbc, err := pgx.Connect(context.Background(), url)

	if err != nil {
		log.Fatal(err)
	}

	defer dbc.Close(context.Background())
	cx := data.NewCx(dbc)
	
	unid.InitDb(cx)
	
	if err := cx.DropAll(); err != nil {
		log.Fatal(err)
	}
	
	if err := cx.SyncAll(); err != nil {
		log.Fatal(err)
	}
}
