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

	pools := cx.NewTable("Pools")
	pools.NewForeignKey("Parent", rcs).SetPrimaryKey(true)
	pools.NewForeignKey("Child", rcs).SetPrimaryKey(true)

	caps := cx.NewTable("Caps")
	caps.NewForeignKey("Rc", rcs).SetPrimaryKey(true)
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

		if err := pools.Create(); err != nil {
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
		lodging := newRc("lodging", RcCapTypePool)
		cabins := newRc("cabins", RcCapTypePool)
		
		if err := cabins.AddPool(lodging); err != nil {
			log.Fatal(err)
		}

		cabin1 := newRc("cabin1", RcCapTypeUnit)

		if err := cabin1.AddPool(lodging); err != nil {
			log.Fatal(err)
		}

		if err := cabin1.AddPool(cabins); err != nil {
			log.Fatal(err)
		}

		if err := cabin1.UpdateCaps(time.Now(), MaxTime(), 1, 0); err != nil {
			return err
		}

		rooms := newRc("rooms", RcCapTypePool)

		if err := rooms.AddPool(lodging); err != nil {
			log.Fatal(err)
		}

		newRoom := func(name string) *Rc {
			rc := newRc(name, RcCapTypeUnit)

			if err := rc.AddPool(lodging); err != nil {
				log.Fatal(err)
			}
			
			if err := rc.AddPool(rooms); err != nil {
				log.Fatal(err)
			}
			
			if err := rc.UpdateCaps(time.Now(), MaxTime(), 1, 0); err != nil {
				log.Fatal(err)
			}

			return rc
		}

		newRoom("room1")
		newRoom("room2")
	}

	return nil
}
