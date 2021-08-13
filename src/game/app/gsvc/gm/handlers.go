package gm

import (
	"fw/src/core"
	"fw/src/core/db"
	"fw/src/core/log"
	"fw/src/core/sched/async"
	"fw/src/game/app/dbmgr"
	"fw/src/game/app/gamedata"
	"fw/src/game/app/gconst"
	"fw/src/game/app/modules/lamp"
	"fw/src/game/app/modules/mail"
	"fw/src/game/app/modules/utils"
	"fw/src/game/msg"
	"fw/src/shared/config"
	"net/http"
	"strings"
	"time"
)

// ============================================================================

func handle_plrinfo(req *http.Request) (r string, err error) {
	module := get_string(req, "module")

	plr, err := get_player(req)
	if err != nil {
		return
	}

	switch module {
	case "base":
		r, err = get_plrinfo_base(plr)

	case "ccy":
		r, err = get_plrinfo_ccy(plr)

	case "item":
		r, err = get_plrinfo_item(plr)

	case "hero":
		r, err = get_plrinfo_hero(plr)

	case "guild":
		r, err = get_plrinfo_guild(plr)

	default:
		err = ErrArgs
	}

	return
}

func handle_gldinfo(req *http.Request) (r string, err error) {
	module := get_string(req, "module")

	gld, err := get_guild(req)
	if err != nil {
		return
	}

	switch module {
	case "base":
		r, err = get_gldinfo_base(gld)

	case "members":
		r, err = get_gldinfo_members(gld)

	default:
		err = ErrArgs
	}

	return
}

// ============================================================================

func handle_res(req *http.Request) (r string, err error) {
	plr, err := get_player(req)
	if err != nil {
		return
	}

	arr_key := get_string_arr(req, "res_k")
	arr_val := get_int64_arr(req, "res_v")

	op := plr.GetBag().NewOp(gconst.ObjFrom_GM)

	for i, k := range arr_key {
		p := strings.Split(k, " - ")
		if len(p) != 2 {
			continue
		}

		id := core.Atoi32(p[1])
		n := arr_val[i]

		if id == 0 || n == 0 {
			continue
		}

		op.Add(id, n, 0)
	}

	op.Apply()

	return
}

func handle_hero(req *http.Request) (r string, err error) {
	if !config.Common.DevMode {
		err = ErrNoKey
		return
	}

	plr, err := get_player(req)
	if err != nil {
		return
	}

	arr_id := get_string_arr(req, "hero_k")
	arr_lv := get_int32_arr(req, "lv")
	arr_star := get_int32_arr(req, "star")

	if len(arr_id) == 0 ||
		len(arr_lv) == 0 ||
		len(arr_star) == 0 {
		err = ErrArgs
		return
	}

	p := strings.Split(arr_id[0], " - ")
	if len(p) != 2 {
		return
	}

	id := core.Atoi32(p[1])
	lv := arr_lv[0]
	star := arr_star[0]

	// add new hero
	if gamedata.ConfMonster.Query(id) == nil {
		err = ErrNoHero
		return
	}

	op := plr.GetBag().NewOp(gconst.ObjFrom_GM)
	op.Inc(id, 1)
	rwds := op.Apply().ToMsg()

	if len(rwds.Heroes) == 0 {
		err = ErrHeroFull
		return
	}

	hero := plr.GetBag().FindHero(rwds.Heroes[0])
	if hero == nil {
		return
	}

	// check max
	conf_star := gamedata.ConfHeroStarUp.Query(star)
	if conf_star == nil {
		return
	}
	hero.SetStar(star)

	if lv > 1 {
		if lv > conf_star.MaxLv {
			lv = conf_star.MaxLv
		}

		hero.SetLevel(lv)
	}

	return
}

func handle_plr(req *http.Request) (r string, err error) {

	plr, err := get_player(req)
	if err != nil {
		return
	}

	lv := get_int32(req, "lv")
	if lv > gamedata.ConfLimitM.Query().MaxPlrLv {
		lv = gamedata.ConfLimitM.Query().MaxPlrLv
	}

	if lv > plr.GetLevel() {
		plr.SetLevel(lv)
	}

	return
}

func handle_ban(req *http.Request) (r string, err error) {
	plr, err := get_player(req)
	if err != nil {
		return
	}

	// > 0: 封号
	// < 0: 解封
	var ban_ts time.Time

	min := get_int32(req, "ban_acct")
	if min == 0 {
		return
	} else if min > 0 {
		ban_ts = time.Now().Add(time.Duration(min) * time.Minute)
	} else {
		ban_ts = time.Unix(0, 0)
	}

	// update userinfo
	async.Push(func() {
		err := dbmgr.DBShare.Update(
			dbmgr.C_tabname_userinfo,
			plr.GetId(),
			db.M{"$set": db.M{"ban_ts": ban_ts}},
		)
		if err != nil {
			log.Error("update ban-ts of userinfo failed:", err)
		}
	})

	// update user
	plr.User().BanTs = ban_ts

	// logout player
	if min > 0 {
		plr.Logout()
	}

	return
}

// ============================================================================

func handle_lamp(req *http.Request) (r string, err error) {
	// content, lampid
	content := get_string(req, "content")
	if len(content) == 0 || len(content) > 300 {
		err = ErrArgs
		return
	}
	switch content {
	case "@clearall":
		lamp.ClearAll()
	case "@clearsys":
		lamp.ClearSys()
	case "@clearuser":
		lamp.ClearUser()
	default:
		conf := gamedata.ConfLamp.Query(int32(1))
		if conf == nil {
			err = ErrNoConf
			return
		}

		if conf.Type != 0 {
			return
		}

		p := make(map[string]string)
		p["notice"] = content

		// send
		lamp.Add(conf.Type, conf.Id, p)
	}

	return
}

// ============================================================================

func handle_gmail(req *http.Request) (r string, err error) {

	// title & text
	title := get_string(req, "title")
	text := get_string(req, "text")

	if title == "" || text == "" {
		err = ErrArgs
		return
	}

	gmName := "GM"
	conf_g := gamedata.ConfGlobalPublic.Query(1)
	if conf_g != nil && conf_g.SysMailSender != "" {
		gmName = conf_g.SysMailSender
	}

	// cond
	cond := get_string(req, "cond")

	// new mail
	m := mail.New(nil).
		SetSender(gmName).
		SetTitle(title).
		SetText(text).
		SetCond(cond)

	// attachment
	arr_key := get_string_arr(req, "res_k")
	arr_val := get_int64_arr(req, "res_v")

	for i, k := range arr_key {
		p := strings.Split(k, " - ")
		if len(p) != 2 {
			continue
		}

		id := core.Atoi32(p[1])
		n := arr_val[i]

		if id == 0 || n == 0 {
			continue
		}

		m.AddAttachment(id, float64(n))
	}

	// send
	m.Send()

	return
}

func handle_pmail(req *http.Request) (r string, err error) {

	plr, err := get_player(req)
	if err != nil {
		return
	}

	// title & text
	title := get_string(req, "title")
	text := get_string(req, "text")

	if title == "" || text == "" {
		err = ErrArgs
		return
	}

	gmName := "GM"
	conf_g := gamedata.ConfGlobalPublic.Query(1)
	if conf_g != nil && conf_g.SysMailSender != "" {
		gmName = conf_g.SysMailSender
	}

	// new mail
	m := mail.New(plr).
		SetSender(gmName).
		SetTitle(title).
		SetText(text)

	// attachment
	arr_key := get_string_arr(req, "res_k")
	arr_val := get_int64_arr(req, "res_v")

	for i, k := range arr_key {
		p := strings.Split(k, " - ")
		if len(p) != 2 {
			continue
		}

		id := core.Atoi32(p[1])
		n := arr_val[i]

		if id == 0 || n == 0 {
			continue
		}

		m.AddAttachment(id, float64(n))
	}

	// send
	m.Send()

	return
}

func handle_fake_bill(req *http.Request) (r string, err error) {

	plr, err := get_player(req)
	if err != nil {
		return
	}

	prod_id := get_int32(req, "prod_id")
	csext := get_string(req, "csext")

	conf := gamedata.ConfBillProduct.Query(prod_id)
	if conf == nil {
		err = ErrNoConf
		return
	}

	if plr.GetSdk() != conf.Sdk {
		err = ErrSdk
		return
	}

	plr.GetBill().GiveItems(prod_id, csext, conf.Price, conf.Ccy)

	return
}

// ============================================================================

func handle_conf(req *http.Request) (r string, err error) {
	names := get_string_arr(req, "conf_k")

	// reload conf
	b := gamedata.ReloadConf(names)

	// notify all online players
	if b {
		utils.BroadcastPlayers(&msg.GS_GameDataReloaded{})
	}

	return
}
