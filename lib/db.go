package unid

import (
	"github.com/codr7/unid/lib/db"
	"log"
	"time"
)

func InitDb(cx *db.Cx) error {
	users := cx.NewTable("Users")
	users.NewStringCol("Name").SetPrimaryKey(true)
	users.NewTimeCol("CreatedAt")
	
	rcCapType := cx.NewEnum("RcCapType", "free", "pool", "unit")
	
	rcs := cx.NewTable("Rcs")
	rcs.NewForeignKey("CreatedBy", users)
	rcs.NewStringCol("Name").SetPrimaryKey(true)
	rcs.NewTimeCol("CreatedAt")
	rcs.NewEnumCol("CapType", rcCapType)
	
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

		adm := NewUser(cx)
		adm.Name = "adm"

		if err := db.Store(adm); err != nil {
			log.Fatal(err)
		}

		if err := rcCapType.Create(); err != nil {
			return err
		}

		if err := rcs.Create(); err != nil {
			return err
		}

		if err := caps.Create(); err != nil {
			return err
		}

		newRc := func(name string, capType string) *Rc {
			rc := NewRc(cx)
			rc.CreatedBy = adm
			rc.Name = name
			rc.CapType = capType
			
			if err := db.Store(rc); err != nil {
				log.Fatal(err)
			}

			return rc
		}
		
		newRc("breakfast", RcCapTypeFree)
		newRc("lodging", RcCapTypePool)
		newRc("cabins", RcCapTypePool)

		cabin1 := newRc("cabin1", RcCapTypeUnit)

		if err := cabin1.UpdateCaps(time.Now(), MaxTime(), 1, 0); err != nil {
			return err
		}

		newRc("rooms", RcCapTypePool)
		
		room1 := newRc("room1", RcCapTypeUnit)

		if err := room1.UpdateCaps(time.Now(), MaxTime(), 1, 0); err != nil {
			return err
		}
	}

	return nil
}
