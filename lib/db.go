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
	rcs.NewStringCol("Name").SetPrimaryKey(true)
	rcs.NewTimeCol("CreatedAt")
	rcs.NewForeignKey("CreatedBy", users)
	
	caps := cx.NewTable("Caps")
	caps.NewTimeCol("StartsAt").SetPrimaryKey(true)
	caps.NewTimeCol("EndsAt")
	caps.NewIntCol("Total")
	caps.NewIntCol("Used")
	caps.NewForeignKey("Rc", rcs)
	caps.FindCol("RcName").SetPrimaryKey(true)
	
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
		
		lodging := NewRc(cx)
		lodging.CreatedBy = admin
		lodging.Name = "lodging"
		data.Store(lodging)
	}

	return nil
}
