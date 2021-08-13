package gm

import (
	"encoding/json"
	"fmt"
	"fw/src/core"
	"fw/src/game/app"
	"fw/src/game/app/gamedata"
	"fw/src/game/app/modules/guild"
	"math"
	"sort"
	"time"
)

// ============================================================================

func get_plrinfo_base(plr *app.Player) (r string, err error) {
	user := plr.User()

	var tab [][]interface{}

	tab = append(tab, []interface{}{"Num", "Field", "Value"})

	tab = append(tab, []interface{}{"*", "UserId", user.Id})
	tab = append(tab, []interface{}{"*", "AuthId", user.AuthId})
	tab = append(tab, []interface{}{"*", "Svr0", user.Svr0})
	tab = append(tab, []interface{}{"*", "Svr", user.Svr})
	tab = append(tab, []interface{}{"*", "Sdk", user.Sdk})
	tab = append(tab, []interface{}{"*", "Model", user.Model})
	tab = append(tab, []interface{}{"*", "DeviceId", user.DevId})
	tab = append(tab, []interface{}{"*", "Os", user.Os})
	tab = append(tab, []interface{}{"*", "OsVer", user.OsVer})
	tab = append(tab, []interface{}{"*", "Name", user.Name})
	tab = append(tab, []interface{}{"*", "Level", user.Lv})
	tab = append(tab, []interface{}{"*", "Vip", user.Vip.Lv})
	tab = append(tab, []interface{}{"*", "CreateTime", user.CreateTs.Format("2006-01-02T 15:04:05")})
	tab = append(tab, []interface{}{"*", "LoginTime", user.LoginTs.Format("2006-01-02T 15:04:05")})
	tab = append(tab, []interface{}{"*", "LoginIP", user.LoginIP})

	// 在线情况
	{
		var text string

		if plr.IsOnline() {
			text = "<b style='color:green'>Online</b>"
		} else {
			text = "Offline"
		}

		tab = append(tab, []interface{}{"*", "Online-Status", text})
	}

	// 累计在线时长
	{
		v := fmt.Sprintf("%d min", int32(math.Floor((time.Second * time.Duration(user.OnlineDur)).Minutes())))

		tab = append(tab, []interface{}{"*", "AccOnlineDur", v})
	}

	// 封号情况
	{
		var ban_text string
		if user.BanTs.After(time.Now()) {
			ban_text = "<b style='color:red'>Banned</b>"
		} else {
			ban_text = "Normal"
		}

		tab = append(tab, []interface{}{"*", "Ban-Status", ban_text})
		tab = append(tab, []interface{}{"*", "BanUtil", user.BanTs.Format("2006-01-02T 15:04:05")})
	}

	// marshal
	data, err := json.Marshal(&tab)
	if err == nil {
		r = string(data)
	}

	return
}

func get_plrinfo_ccy(plr *app.Player) (r string, err error) {
	var tab [][]interface{}

	tab = append(tab, []interface{}{"Name", "Amount"})

	for k, v := range plr.GetBag().Ccy {
		name := core.I32toa(k)
		conf := gamedata.ConfCurrency.Query(k)
		if conf != nil {
			name = conf.Txt_Name
		}

		tab = append(tab, []interface{}{name, v})
	}

	// marshal
	data, err := json.Marshal(&tab)
	if err == nil {
		r = string(data)
	}

	return
}

func get_plrinfo_item(plr *app.Player) (r string, err error) {
	var tab [][]interface{}

	tab = append(tab, []interface{}{"Name", "Amount"})

	for k, v := range plr.GetBag().Items {
		name := core.I32toa(k)
		conf := gamedata.ConfItem.Query(k)
		if conf != nil {
			name = conf.Txt_Name
		}

		tab = append(tab, []interface{}{name, v})
	}

	// marshal
	data, err := json.Marshal(&tab)
	if err == nil {
		r = string(data)
	}

	return
}

func get_plrinfo_hero(plr *app.Player) (r string, err error) {
	var tab [][]interface{}

	tab = append(tab, []interface{}{"Name", "Level", "Star", "TrinketLv"})

	for _, hero := range plr.GetBag().Heroes {
		name := core.I32toa(hero.Id)
		conf := gamedata.ConfMonster.Query(hero.Id)
		if conf != nil {
			name = conf.Name
		}

		tab = append(tab, []interface{}{name, hero.Lv, hero.Star, hero.Trinket.Lv})
	}

	// marshal
	data, err := json.Marshal(&tab)
	if err == nil {
		r = string(data)
	}

	return
}

func get_plrinfo_guild(plr *app.Player) (r string, err error) {
	var tab [][]interface{}

	tab = append(tab, []interface{}{"Num", "Field", "Value"})

	gld := plr.GetGuild()
	if gld != nil {
		tab = append(tab, []interface{}{"*", "Id", gld.GetId()})
		tab = append(tab, []interface{}{"*", "Name", gld.GetName()})
		tab = append(tab, []interface{}{"*", "Owner", gld.Owner().GetName()})
		tab = append(tab, []interface{}{"*", "Level", gld.GetLevel()})
		tab = append(tab, []interface{}{"*", "MemberCount", len(gld.Members)})
		tab = append(tab, []interface{}{"*", "Rank", plr.GetGuildRank()})
	}

	// marshal
	data, err := json.Marshal(&tab)
	if err == nil {
		r = string(data)
	}

	return
}

// ============================================================================

func get_gldinfo_base(gld *guild.Guild) (r string, err error) {
	var tab [][]interface{}

	tab = append(tab, []interface{}{"Num", "Field", "Value"})

	tab = append(tab, []interface{}{"*", "Id", gld.GetId()})
	tab = append(tab, []interface{}{"*", "Name", gld.GetName()})
	tab = append(tab, []interface{}{"*", "Owner", gld.Owner().GetName()})
	tab = append(tab, []interface{}{"*", "Level", gld.GetLevel()})
	tab = append(tab, []interface{}{"*", "MemberCount", len(gld.Members)})
	tab = append(tab, []interface{}{"*", "AtkPwr", gld.GetAtkPwr()})
	tab = append(tab, []interface{}{"*", "Rank", gld.GetRank()})

	// marshal
	data, err := json.Marshal(&tab)
	if err == nil {
		r = string(data)
	}

	return
}

func get_gldinfo_members(gld *guild.Guild) (r string, err error) {
	var tab [][]interface{}

	tab = append(tab, []interface{}{"Id", "Name", "Level", "AtkPwr", "Rank", "Online"})

	rk_dict := map[int32]string{
		1: "<b style='color:#ff6000'>Owner</b>",
		2: "VP",
		3: "Member",
	}

	ol_dict := map[bool]string{
		true:  "<b style='color:green'>Y</b>",
		false: "--",
	}

	for _, mb := range gld.Members {
		plr := app.PlayerMgr.LoadPlayer(mb.Id)
		if plr == nil {
			return
		}

		tab = append(tab, []interface{}{
			plr.GetId(),
			plr.GetName(),
			plr.GetLevel(),
			plr.GetAtkPwr(),
			rk_dict[mb.Rank],
			ol_dict[plr.IsOnline()],
		})
	}

	sort.Slice(tab[1:], func(i, j int) bool {
		return tab[i+1][3].(int32) > tab[j+1][3].(int32)
	})

	// marshal
	data, err := json.Marshal(&tab)
	if err == nil {
		r = string(data)
	}

	return
}
