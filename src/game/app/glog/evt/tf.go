package evt

import (
	"fw/src/game/app"
	"fw/src/game/msg"
)

// ============================================================================

type log_chg_t struct {
	Id  int32 `json:"id"`
	Op  int64 `json:"op"`
	Cur int64 `json:"cur"`
}

type log_team_hero_t struct {
	Pos    int32         `json:"pos"`
	Id     int32         `json:"id"`
	Lv     int32         `json:"lv"`
	Star   int32         `json:"star"`
	TkLv   int32         `json:"tklv"`
	Armors []int32       `json:"armors"`
	Relic  []*log_id_n_t `json:"relic"`
}

type log_draw_hero_t struct {
	Id   int32 `json:"id"`
	Star int32 `json:"n"`
}

type log_id_n_t struct {
	Id int32 `json:"id"`
	N  int64 `json:"n"`
}

// ============================================================================

func team_log(plr *app.Player, T *msg.TeamFormation) []*log_team_hero_t {
	ret := []*log_team_hero_t{}
	if T == nil {
		return ret
	}

	for seq, pos := range T.Formation {
		hero := plr.GetBag().FindHero(seq)
		if hero == nil {
			continue
		}

		hero_info := hero.ToMsg_Detail()
		one := &log_team_hero_t{
			Pos:    pos,
			Id:     hero.Id,
			Lv:     hero.Lv,
			Star:   hero.Star,
			TkLv:   hero_info.Hero.Trinket.Lv,
			Armors: hero_info.Armors,
		}

		for id, star := range hero_info.Relic {
			one.Relic = append(one.Relic, &log_id_n_t{
				Id: id,
				N:  int64(star),
			})
		}

		ret = append(ret, one)
	}

	return ret
}
