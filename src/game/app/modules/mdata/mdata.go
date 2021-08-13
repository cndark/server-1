package mdata

import (
	"fw/src/core"
	"fw/src/core/db"
	"fw/src/core/log"
	"fw/src/core/sched/async"
	"fw/src/core/sched/loop"
	"fw/src/game/app/dbmgr"
	"math/rand"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

// ============================================================================

var (
	m_data = make(map[string]*mdata_t) // [module-name]
)

// ============================================================================

type mdata_t struct {
	Data  interface{}
	utype bool // is data of user-type ?
}

// ============================================================================

func Open() {
	md_load_data()
	md_save_start()
}

func Close() {
	md_save()

	log.Info("module data is saved")
}

func md_load_data() {
	for _, reg := range m_registry {
		// load data
		var md mdata_t

		if reg.LoadData == nil {
			// system flow
			err := dbmgr.DBGame.GetObject(
				dbmgr.C_tabname_mdata,
				reg.Name,
				&md,
			)
			if err != nil && !db.IsNotFound(err) {
				core.Panic("loading module data failed:", reg.Name, err)
			}
		} else {
			// customized
			data := reg.LoadData()

			md.Data = data
			md.utype = true
		}

		// add
		m_data[reg.Name] = &md

		// data loaded
		if reg.DataLoaded != nil {
			reg.DataLoaded()
		}
	}
}

func Get(name string) interface{} {
	// get reg
	reg := m_registry[name]
	if reg == nil {
		return nil
	}

	// find
	md := m_data[name]
	if md == nil {
		md = &mdata_t{}
		m_data[name] = md
	}

	// check conversion
	if !md.utype {
		md.Data = md_to_utype(md.Data, reg)
		md.utype = true
	}

	// return
	return md.Data
}

func md_to_utype(in interface{}, reg *Reg) (out interface{}) {
	out = reg.NewModuleData()

	if in == nil {
		return
	}

	data, err := bson.Marshal(in)
	if err != nil {
		log.Error("marshal module data failed:", reg.Name, err)
		return
	}

	err = bson.Unmarshal(data, out)
	if err != nil {
		log.Error("unmarshal module data failed:", reg.Name, err)
		return
	}

	return
}

// ============================================================================
// save

func md_save_start() {
	loop.SetTimeout(time.Now().Add(time.Duration(1200+rand.Intn(1200))*time.Second), func() {
		md_save_async()
		md_save_start()
	})
}

func md_save_async() {
	now := time.Now()

	for _, reg := range m_registry {
		reg := reg

		md := m_data[reg.Name]
		if md == nil {
			continue
		}

		loop.SetTimeout(now.Add(time.Duration(rand.Intn(1000))*time.Second), func() {
			if reg.LoadData == nil { // only check 'LoadData' for consistency
				// system flow
				_id := reg.Name
				doc := core.CloneBsonObject(md)

				async.Push(func() {
					err := dbmgr.DBGame.Replace(
						dbmgr.C_tabname_mdata,
						_id,
						doc,
					)
					if err != nil {
						log.Error("saving module data failed:", _id, err)
					}
				})
			} else {
				// customized
				reg.SaveAsync()
			}
		})
	}
}

func md_save() {
	for _, reg := range m_registry {
		md := m_data[reg.Name]
		if md == nil {
			continue
		}

		if reg.LoadData == nil {
			// system flow
			_id := reg.Name
			doc := md

			err := dbmgr.DBGame.Replace(
				dbmgr.C_tabname_mdata,
				_id,
				doc,
			)
			if err != nil {
				log.Error("saving module data failed:", _id, err)
			}
		} else {
			// customized
			reg.Save()
		}
	}
}
