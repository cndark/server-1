package botconf

import (
	"encoding/json"
	"fw/src/core"
	"io/ioutil"
	"strings"
)

// ============================================================================

type config_t struct {
	Version   string    `json:"version"`
	Ver_Major string    `json:"-"`
	Ver_Minor string    `json:"-"`
	Ver_Build string    `json:"-"`
	Pprof     *pprof_t  `json:"pprof"`
	Password  string    `json:"pwd"`
	Bots      []*bot_t  `json:"bots"`
	JobCtl    *jobctl_t `json:"jobctl"`
}

type pprof_t struct {
	Enabled bool
	Port    int
}

type bot_t struct {
	From      int32    `json:"from"`
	To        int32    `json:"to"`
	Addr      string   `json:"addr"`
	Svr       string   `json:"svr"`
	Sdk       string   `json:"sdk"`
	Model     string   `json:"model"`
	JobPrefix []string `json:"job_prefix"`
	JobItv    []int32  `json:"job_itv"`
	Args      string   `json:"args"`
	Grp       int32    `json:"grp"`
}

type jobctl_t struct {
	GrpBinding   map[string][]int32 `json:"grp_binding"`
	SelectWeight map[string]int32   `json:"sel_weight"`
	TriggerProb  map[string]float32 `json:"tri_prob"`
}

// ============================================================================

var (
	Ver_Major string
	Ver_Minor string
	Ver_Build string

	Pprof    *pprof_t
	Password string
	Bots     []*bot_t
	JobCtl   *jobctl_t
)

// ============================================================================

func Parse(fn string) {
	var conf config_t

	// read file
	d, err := ioutil.ReadFile(fn)
	if err != nil {
		core.Panic("open config file failed:", err)
	}

	// parse
	err = json.Unmarshal(d, &conf)
	if err != nil {
		core.Panic("parse config file failed:", err)
	}

	parse_version(&conf)

	// set variables
	Ver_Major = conf.Ver_Major
	Ver_Minor = conf.Ver_Minor
	Ver_Build = conf.Ver_Build

	Pprof = conf.Pprof
	Password = conf.Password
	Bots = conf.Bots
	JobCtl = conf.JobCtl
}

// ============================================================================

func parse_version(conf *config_t) {
	arr := strings.Split(conf.Version, ".")
	if len(arr) < 3 {
		core.Panic("invalid version:", conf.Version)
	}

	conf.Ver_Major = arr[0]
	conf.Ver_Minor = arr[1]
	conf.Ver_Build = arr[2]
}
