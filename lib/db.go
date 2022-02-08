package unid

import (
	"github.com/codr7/unid/lib/data"
)

func InitDb(cx *data.Cx) {
	users := cx.NewTable("Users", data.NewStringCol("Name"))
	
	rcs := cx.NewTable("Rcs", data.NewStringCol("Name"))
	rcs.NewForeignKey("CreatedBy", users)
	
	caps := cx.NewTable("Caps", data.NewStringCol("RcName"), data.NewTimeCol("StartsAt"))
	caps.NewForeignKey("Rc", rcs)
	caps.AddCols(data.NewTimeCol("EndsAt"), data.NewIntCol("Total"), data.NewIntCol("Used"))
}
