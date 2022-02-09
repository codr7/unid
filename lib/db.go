package unid

import (
	"github.com/codr7/unid/lib/data"
	"log"
)

func InitDb(cx *data.Cx) error {
	users := cx.NewTable("Users", data.NewStringCol("Name"))
	
	rcs := cx.NewTable("Rcs", data.NewStringCol("Name"))
	rcs.NewForeignKey("CreatedBy", users)
	rcs.AddCol(data.NewTimeCol("CreatedAt"))
	
	caps := cx.NewTable("Caps", data.NewStringCol("RcName"), data.NewTimeCol("StartsAt"))
	caps.NewForeignKey("Rc", rcs)
	caps.AddCol(data.NewTimeCol("EndsAt"), data.NewIntCol("Total"), data.NewIntCol("Used"))
	
	if err := cx.DropAll(); err != nil {
		return err
	}

	if ok, err := users.Exists(); err != nil {
		return err
	} else if !ok {
		users.Create()

		admin := NewUser(cx)
		admin.Name = "admin"
		if err := data.Store(admin); err != nil {
			log.Fatal(err)
		}

		rcs.Create()
		
		lodging := NewRc(cx)
		lodging.CreatedBy = admin
		lodging.Name = "lodging"
		data.Store(lodging)
	}

	return nil
}
