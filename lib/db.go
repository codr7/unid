package unid

import (
	"github.com/codr7/unid/lib/data"
	"log"
)

func InitDb(cx *data.Cx) error {
	users := cx.NewTable("Users")
	users.NewStringCol("Name").SetPrimaryKey(true)
	users.NewTimeCol("CreatedAt")
	
	rcs := cx.NewTable("Rcs")
	rcs.NewForeignKey("CreatedBy", users)
	rcs.NewStringCol("Name").SetPrimaryKey(true)
	rcs.NewTimeCol("CreatedAt")
	
	caps := cx.NewTable("Caps")
	caps.NewForeignKey("Rc", rcs)
	caps.FindCol("RcName").SetPrimaryKey(true)
	caps.NewTimeCol("StartsAt").SetPrimaryKey(true)
	caps.NewTimeCol("EndsAt")
	caps.NewIntCol("Total")
	caps.NewIntCol("Used")
	caps.NewTimeCol("ChangedAt")
	
	if err := cx.DropAll(); err != nil {
		return err
	}

	if ok, err := users.Exists(); err != nil {
		return err
	} else if !ok {
		if err := users.Create(); err != nil {
			return err
		}

		admin := NewUser(cx)
		admin.Name = "admin"

		if err := data.Store(admin); err != nil {
			log.Fatal(err)
		}

		if err := rcs.Create(); err != nil {
			return err
		}

		if err := caps.Create(); err != nil {
			return err
		}

		newRc := func(name string) *Rc {
			rc := NewRc(cx)
			rc.CreatedBy = admin
			rc.Name = name

			if err := data.Store(rc); err != nil {
				log.Fatal(err)
			}

			return rc
		}
		
		newRc("lodging")
		newRc("cabins")
		newRc("rooms")
	}

	return nil
}
