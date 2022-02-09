package unid

import (
	"github.com/codr7/unid/lib/data"
)

func InitDb(cx *data.Cx) error {
	users := cx.NewTable("Users", data.NewStringCol("Name"))
	
	rcs := cx.NewTable("Rcs", data.NewStringCol("Name"))
	rcs.NewForeignKey("CreatedBy", users)
	
	caps := cx.NewTable("Caps", data.NewStringCol("RcName"), data.NewTimeCol("StartsAt"))
	caps.NewForeignKey("Rc", rcs)
	caps.AddCols(data.NewTimeCol("EndsAt"), data.NewIntCol("Total"), data.NewIntCol("Used"))
	
	if err := cx.DropAll(); err != nil {
		return err
	}

	if ok, err := users.Exists(); err != nil {
		return err
	} else if !ok {
		users.Create()

		admin := NewUser(cx)
		admin.Name = "admin"
		data.Store(admin)
	}

	return nil
}
