package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"github.com/jackc/pgx/v4"
	"github.com/codr7/unid/lib/data"
)

func main() {
	url := "postgres://unid:unid@localhost:5432/unid"
	dbc, err := pgx.Connect(context.Background(), url)

	if err != nil {
		log.Fatal(err)
	}

	defer dbc.Close(context.Background())
	cx := data.NewCx(dbc)
	
	users := cx.NewTable("Users", data.NewStringCol("Name"))

	rcs := cx.NewTable("Rcs", data.NewStringCol("Name"))
	rcs.NewForeignKey("CreatedBy", users)
	
	caps := cx.NewTable("Caps", data.NewStringCol("RcName"), data.NewTimeCol("StartsAt"))
	caps.NewForeignKey("Rc", rcs)
	caps.AddCols(data.NewTimeCol("EndsAt"), data.NewIntCol("Total"), data.NewIntCol("Used"))

	cx.DropAll()
	cx.CreateAll()
}
