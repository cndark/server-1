package gamedata

import (
	"encoding/json"
	"fw/src/core"
	"fw/src/core/evtmgr"
	"fw/src/core/log"
	"fw/src/core/wordsfilter"
	"fw/src/game/app/gconst"
	"io/ioutil"
	"reflect"
	"strings"
)

// ============================================================================

const (
	C_path_data   = "gamedata/data/"
	C_path_filter = "gamedata/filter/"
)

// ============================================================================

var (
	confs_by_name map[string]IConf
)

// ============================================================================

type IConf interface {
	Load()
}

// ============================================================================

func Load() {
	// load gamedata
	log.Info("Loading game-data ...")

	for _, conf := range confs_all {
		conf.Load()
	}
	for _, conf := range confs_manual {
		conf.Load()
	}

	// load words filter
	log.Info("Loading words-filter ...")
	err := wordsfilter.Load(C_path_filter + "filter.txt")
	if err != nil {
		log.Error("loading words-filter failed:", err)
	}

	// build conf index: name -> iconf
	confs_by_name = make(map[string]IConf)
	for _, conf := range confs_all {
		t := reflect.TypeOf(conf).Elem()
		name := strings.TrimSuffix(t.Name(), "Table")

		confs_by_name[name] = conf
	}

	evtmgr.Fire(gconst.Evt_GameDataLoaded)
}

func ReloadConf(names []string) bool {
	// reload specified confs
	b := false
	for _, name := range names {
		conf := confs_by_name[name]
		if conf != nil {
			conf.Load()
			b = true
		}
	}

	// reload all manual confs
	if b {
		for _, conf := range confs_manual {
			conf.Load()
		}

		// fire
		evtmgr.Fire(gconst.Evt_GameDataReload, names)
	}

	return b
}

// ============================================================================

func load_json(fn string, obj interface{}) bool {
	log.Infof("    %-24s loaded", fn)

	data, err := ioutil.ReadFile(C_path_data + fn)
	if err != nil {
		core.Panic("reading gamedata failed:", fn, err)
		return false
	}

	err = json.Unmarshal(data, &obj)
	if err != nil {
		core.Panic("parsing gamedata failed:", fn, err)
		return false
	}

	return true
}
