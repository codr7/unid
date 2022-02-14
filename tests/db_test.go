package tests

import (
	"context"
	"github.com/codr7/unid/lib"
	"github.com/codr7/unid/lib/db"
	"github.com/jackc/pgx/v4"
	"testing"
)

func TestForeignKey(t *testing.T) {
	const TEST_NAME = "TestForeignKey"
	
	url := "postgres://test:test@localhost:5432/test"
	dbc, err := pgx.Connect(context.Background(), url)

	if err != nil {
		t.Fatal(err)
	}

	defer dbc.Close(context.Background())
	cx := db.NewCx(dbc)
	
	if err := unid.InitDb(cx); err != nil {
		t.Fatal(err)
	}
	
	if err := cx.SyncAll(); err != nil {
		t.Fatal(err)
	}

	u := unid.NewUser(cx)
	u.Name = TEST_NAME
	
	if err := db.Store(u); err != nil {
		t.Fatal(err)
	}
			
	rc := unid.NewRc(cx)
	rc.CreatedBy = u
	rc.Name = TEST_NAME
	db.Store(rc)

	rc.CreatedBy = nil

	if err := cx.FindTable("Rcs").Load(rc); err != nil {
		t.Fatal(err)
	}
		
	if createdBy, err := rc.GetCreatedBy(); err != nil {
		t.Fatal(err)
	} else if createdBy.Name != TEST_NAME {
		t.Fatalf("Wrong name: %v", createdBy.Name)
	}
}

	
