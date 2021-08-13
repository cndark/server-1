package battle

import (
	"encoding/json"
	"fw/src/core"
	"fw/src/core/log"
	"fw/src/core/sched/loop"
	"fw/src/game/app/gamedata"
	"fw/src/game/msg"
	"math/rand"
	"strings"
)

// ============================================================================

type monster_team_t struct {
	fts []*msg.BattleFighter
}

// ============================================================================

func NewMonsterTeam() *monster_team_t {
	return &monster_team_t{}
}

func (self *monster_team_t) AddMonster(id int32, lv int32, pos int32) *monster_team_t {
	// conf
	conf_mon := gamedata.ConfMonster.Query(id)
	if conf_mon == nil {
		return self
	}

	conf_power := gamedata.ConfMonsterPower.Query(lv)
	if conf_power == nil {
		return self
	}

	conf_bat := gamedata.ConfGlobalBattle.Query(1)
	if conf_bat == nil {
		return self
	}

	// check pos
	if pos < 0 {
		return self
	}

	// calc props
	props := make(map[int32]float32)

	for i, v := range conf_bat.MonsterBaseProps {
		props[v.Id] += v.Val * conf_mon.MonsterPropsRatio[i]
	}
	for i, v := range conf_power.BaseProps {
		props[v.Id] += v.Val * conf_mon.MonsterPropsRatio[i]
	}

	// delete zero-value props
	var to_del []int32
	for k, v := range props {
		if v == 0 {
			to_del = append(to_del, k)
		}
	}
	for _, k := range to_del {
		delete(props, k)
	}

	// add
	self.fts = append(self.fts, &msg.BattleFighter{
		Id:    id,
		Lv:    lv,
		Star:  conf_power.Star,
		Props: props,
		Pos:   pos,
	})

	return self
}

func (self *monster_team_t) ModifyProps(ratio float32) {
	for _, ft := range self.fts {
		for id, v := range ft.Props {
			ft.Props[id] = v * ratio
		}
	}

	return
}

func (self *monster_team_t) ToMsg_BattleTeam() *msg.BattleTeam {
	self.check()

	return &msg.BattleTeam{
		Player:   nil,
		Fighters: self.fts,
	}
}

// 策划小子说的: 如果怪物阵容有6号位, 就不要3,4号位的怪. 成全他
func (self *monster_team_t) check() {
	for _, v := range self.fts {
		if v.Pos == 6 {
			self.remove_fighter(3)
			self.remove_fighter(4)
			break
		}
	}
}

func (self *monster_team_t) remove_fighter(pos int32) {
	for i, v := range self.fts {
		if v.Pos == pos {
			L := len(self.fts)
			self.fts[i] = self.fts[L-1]
			self.fts = self.fts[:L-1]
			break
		}
	}
}

// ============================================================================

func Fight(input *msg.BattleInput, f func(r *msg.BattleResult)) {
	// fail func
	fail := func() {
		loop.Push(func() {
			f(nil)
		})
	}

	// choose a bat
	bat := choose_bat()
	if bat == nil {
		fail()
		return
	}

	// gen seed
	if input.Args == nil {
		input.Args = make(map[string]string)
	}

	input.Args["seed"] = core.I32toa(rand.Int31())

	// http to bat
	core.Go(func() {
		// prepare input
		data, err := json.Marshal(input)
		if err != nil {
			log.Warning("battle input marshal error:", err)
			fail()
			return
		}

		// send for calculation
		ret := core.HttpPostJson(bat.addr, string(data))
		if ret == "" {
			fail()
			return
		}

		// get result
		var r *msg.BattleResult
		err = json.Unmarshal([]byte(ret), &r)
		if err != nil {
			log.Warning("battle result unmarshal error:", err)
			fail()
			return
		}

		// callback
		loop.Push(func() {
			f(r)
		})
	})
}

// ============================================================================

func BattleResultHpLoss(args map[string]string) (map[int32]float64, map[int32]float64) {
	t1 := make(map[int32]float64)
	t2 := make(map[int32]float64)

	for s, v := range args {
		if strings.HasPrefix(s, "hp_loss") {
			arr := strings.Split(s, ".")
			if len(arr) != 3 {
				continue
			}

			if arr[1] == "1" {
				t1[core.Atoi32(arr[2])] = core.Atof64(v)
			} else if arr[1] == "2" {
				t2[core.Atoi32(arr[2])] = core.Atof64(v)
			}
		}
	}

	return t1, t2
}
