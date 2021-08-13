package gm

import (
	"errors"
	"net/http"
)

// ============================================================================

var (
	ErrNoKey    = errors.New("invalid key")
	ErrArgs     = errors.New("invalid params")
	ErrNoPlayer = errors.New("player not found")
	ErrNoHero   = errors.New("hero not found")
	ErrHeroFull = errors.New("hero full")
	ErrNoGuild  = errors.New("guild not found")
	ErrNoConf   = errors.New("conf file not found")
	ErrSdk      = errors.New("invalid sdk")
)

// ============================================================================

var handlers = map[string]func(*http.Request) (r string, err error){
	"dev": handle_dev,

	"plrinfo": handle_plrinfo,
	"gldinfo": handle_gldinfo,

	"w.res":      handle_res,
	"w.hero":     handle_hero,
	"w.plr":      handle_plr,
	"w.ban":      handle_ban,
	"w.fakebill": handle_fake_bill,

	"w.lamp":  handle_lamp,
	"w.gmail": handle_gmail,
	"w.pmail": handle_pmail,

	"w.conf": handle_conf,
}
