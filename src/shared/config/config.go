package config

import (
	"encoding/json"
	"fw/src/core"
	"fw/src/core/log"
	"io/ioutil"
	"regexp"
	"sort"
	"strings"
)

// ============================================================================

type common_t struct {
	Version     string            `json:"version"`
	Ver_Major   string            `json:"-"`
	Ver_Minor   string            `json:"-"`
	Ver_Build   string            `json:"-"`
	Area        *area_t           `json:"area"`
	DBShare     string            `json:"db_share"`
	DBCenter    string            `json:"db_center"`
	DBStats     string            `json:"db_stats"`
	DBLog       string            `json:"db_log"`
	DBBill      string            `json:"db_bill"`
	DBCross     string            `json:"db_cross"`
	DBUser      map[string]string `json:"db_user"`
	Kfk         *kfk_t            `json:"kfk"`
	Redis       *redis_t          `json:"redis"`
	BehindProxy bool              `json:"behind_proxy"`
	DevMode     bool              `json:"dev_mode"`
}

type area_t struct {
	Id   int32  `json:"id"`
	Name string `json:"name"`
}

type kfk_t struct {
	Urls []string `json:"urls"`
}

type redis_t struct {
	PoolSize int      `json:"pool_size"`
	Urls     []string `json:"urls"`
}

type auth_t struct {
	IP   string `json:"ip"`
	Port int    `json:"port"`
}

type bill_t struct {
	IP   string `json:"ip"`
	Port int    `json:"port"`
}

type switcher_t struct {
	IP    string `json:"ip"`
	Port  int    `json:"port"`
	Token string `json:"token"`
}

type reporter_t struct {
	IP    string `json:"ip"`
	Port  int    `json:"port"`
	Token string `json:"token"`
}

type agent_t struct {
	IP   string `json:"ip"`
	Port int    `json:"port"`
}

type admin_t struct {
	IP   string `json:"ip"`
	Port int    `json:"port"`
}

type router_t struct {
	Id     int32  `json:"-"`
	Name   string `json:"-"`
	Addr4C string `json:"addr4c"`
}

type bat_t struct {
	Id   int32  `json:"-"`
	Name string `json:"-"`
	IP   string `json:"ip"`
	Port int    `json:"port"`
}

type game_t struct {
	Id      int32  `json:"-"`
	Name    string `json:"-"`
	Addr4GW string `json:"addr4gw"`
	Svc     string `json:"svc"`
	DBGame  string `json:"db_game"`
}

type gate_t struct {
	Id     int32  `json:"-"`
	Name   string `json:"-"`
	IPWan  string `json:"ip_wan"`
	Port   int32  `json:"port"`
	WsPort int32  `json:"wsport"`
}

// ============================================================================

type config_t struct {
	Common   *common_t            `json:"common"`
	Auth     *auth_t              `json:"auth"`
	Bill     *bill_t              `json:"bill"`
	Switcher *switcher_t          `json:"switcher"`
	Reporter *reporter_t          `json:"reporter"`
	Agent    *agent_t             `json:"agent"`
	Admin    *admin_t             `json:"admin"`
	Routers  map[string]*router_t `json:"routers"`
	Bats     map[string]*bat_t    `json:"bats"`
	Games    map[string]*game_t   `json:"games"`
	Gates    map[string]*gate_t   `json:"gates"`
}

// ============================================================================

var (
	Common   *common_t
	Auth     *auth_t
	Bill     *bill_t
	Switcher *switcher_t
	Reporter *reporter_t
	Agent    *agent_t
	Admin    *admin_t
	Routers  map[string]*router_t
	Bats     map[string]*bat_t
	Games    map[string]*game_t
	Gates    map[string]*gate_t
)

var (
	CurRouter *router_t
	CurGame   *game_t
	CurGate   *gate_t

	GameIds   []int32
	GameIdMax int32
)

// ============================================================================

func Parse(fn string, server string) {
	defer func() {
		if err := recover(); err != nil {
			log.Warning(err)
		}
	}()

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

	parse_id_name(&conf)
	parse_version(&conf)

	// set variables
	Common = conf.Common
	Auth = conf.Auth
	Bill = conf.Bill
	Switcher = conf.Switcher
	Reporter = conf.Reporter
	Agent = conf.Agent
	Admin = conf.Admin
	Routers = conf.Routers
	Bats = conf.Bats
	Games = conf.Games
	Gates = conf.Gates

	// set current
	func() {
		for k, v := range Routers {
			if k == server {
				CurRouter = v
				return
			}
		}

		for k, v := range Games {
			if k == server {
				CurGame = v
				return
			}
		}

		for k, v := range Gates {
			if k == server {
				CurGate = v
				return
			}
		}
	}()

	// game ids
	{
		L := len(Games)
		ids := make([]int32, 0, L)

		for _, v := range Games {
			ids = append(ids, v.Id)
		}

		sort.Slice(ids, func(i, j int) bool {
			return ids[i] < ids[j]
		})

		GameIds = ids

		if L == 0 {
			GameIdMax = 0
		} else {
			GameIdMax = ids[L-1]
		}
	}
}

// ============================================================================

func parse_id_name(conf *config_t) {
	re := regexp.MustCompile(`^[a-z]+(\d+)$`)

	for k, v := range conf.Routers {
		arr := re.FindStringSubmatch(k)
		if arr == nil || len(arr) < 2 {
			core.Panic("invalid batr name:", k)
		}

		v.Id = core.Atoi32(arr[1])
		v.Name = k
	}

	for k, v := range conf.Bats {
		arr := re.FindStringSubmatch(k)
		if arr == nil || len(arr) < 2 {
			core.Panic("invalid bat name:", k)
		}

		v.Id = core.Atoi32(arr[1])
		v.Name = k
	}

	for k, v := range conf.Games {
		arr := re.FindStringSubmatch(k)
		if arr == nil || len(arr) < 2 {
			core.Panic("invalid game name:", k)
		}

		v.Id = core.Atoi32(arr[1])
		v.Name = k
	}

	for k, v := range conf.Gates {
		arr := re.FindStringSubmatch(k)
		if arr == nil || len(arr) < 2 {
			core.Panic("invalid gate name:", k)
		}

		v.Id = core.Atoi32(arr[1])
		v.Name = k
	}
}

func parse_version(conf *config_t) {
	arr := strings.Split(conf.Common.Version, ".")
	if len(arr) < 3 {
		core.Panic("invalid version:", conf.Common.Version)
	}

	conf.Common.Ver_Major = arr[0]
	conf.Common.Ver_Minor = arr[1]
	conf.Common.Ver_Build = arr[2]
}
