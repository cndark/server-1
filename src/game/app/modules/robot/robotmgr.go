package robot

import (
	"fmt"
	"fw/src/game/app/gamedata"
	"fw/src/game/app/modules/mdata"
	"fw/src/game/app/modules/utils"
	"fw/src/game/msg"
	"fw/src/shared/config"
	"math/rand"
	"time"
)

// ============================================================================

var (
	rand_robot = rand.New(rand.NewSource(time.Now().Unix()))

	RobotMgr *robot_mgr_t
)

// ============================================================================

type robot_mgr_t struct {
	Robots map[int32][]*Robot // [lv]
	Index  map[string]*Robot  `bson:"-"`
}

// ============================================================================

func new_data() interface{} {
	return &robot_mgr_t{
		Robots: make(map[int32][]*Robot),
	}
}

func data_loaded() {
	RobotMgr = mdata.Get(NAME).(*robot_mgr_t)
	RobotMgr.init()
}

// ============================================================================

func (self *robot_mgr_t) init() {
	// create if none
	if len(self.Robots) == 0 {
		self.new_robot()
	}

	// shuffle
	for _, a := range self.Robots {
		L := len(a)
		if L > 10 {
			rand.Shuffle(len(a), func(i, j int) {
				a[i], a[j] = a[j], a[i]
			})
		}
	}

	// make index
	self.Index = make(map[string]*Robot)
	for _, a := range self.Robots {
		for _, v := range a {
			self.Index[v.Id] = v
		}
	}
}

func (self *robot_mgr_t) new_robot() {
	self.Robots = make(map[int32][]*Robot)

	f := func(m []int32) int64 {
		L := len(m)
		if L > 0 {
			return int64(m[rand_robot.Intn(L)])
		}
		return 0
	}

	for j, conf := range gamedata.ConfRobot.Items() {
		for i := int32(0); i < conf.Num; i++ {
			id := fmt.Sprintf("bot%d-%d", config.CurGame.Id, (j+1)*100000+i+1)
			if conf.Lv <= 0 {
				conf.Lv = 1
			}

			b := &Robot{
				Id:         id,
				Name:       utils.GenRandName("cn", false),
				Head:       conf.Head,
				HFrame:     conf.HFrame,
				Lv:         conf.Lv,
				AtkPwr:     conf.Power,
				ArenaScore: conf.ArenaInitialScore,
				HeroLv:     make(map[int32]int32),
			}

			b.Team = &msg.TeamFormation{
				Formation: make(map[int64]int32),
			}

			team := make(map[int64]int32)
			seq1, seq2, seq3, seq4, seq5, seq6 := f(conf.Monster1), f(conf.Monster2), f(conf.Monster3),
				f(conf.Monster4), f(conf.Monster5), f(conf.Monster6)

			if seq1 != 0 {
				team[seq1] = 0
				b.HeroLv[0] = self.rand_hero_lv(conf.Lv, conf.MaxLv)
			}
			if seq2 != 0 {
				team[seq2] = 1
				b.HeroLv[1] = self.rand_hero_lv(conf.Lv, conf.MaxLv)
			}
			if seq3 != 0 {
				team[seq3] = 2
				b.HeroLv[2] = self.rand_hero_lv(conf.Lv, conf.MaxLv)
			}
			if seq4 != 0 {
				team[seq4] = 3
				b.HeroLv[3] = self.rand_hero_lv(conf.Lv, conf.MaxLv)
			}
			if seq5 != 0 {
				team[seq5] = 4
				b.HeroLv[4] = self.rand_hero_lv(conf.Lv, conf.MaxLv)
			}
			if seq6 != 0 {
				team[seq6] = 5
				b.HeroLv[5] = self.rand_hero_lv(conf.Lv, conf.MaxLv)
			}

			b.Team.Formation = team

			self.add_robot(b)
		}
	}
}

func (self *robot_mgr_t) rand_hero_lv(lv, maxlv int32) int32 {
	if maxlv <= 0 || lv >= maxlv {
		return lv
	}

	return rand_robot.Int31n(maxlv-lv) + lv
}

func (self *robot_mgr_t) add_robot(b *Robot) {
	self.Robots[b.Lv] = append(self.Robots[b.Lv], b)
}

func (self *robot_mgr_t) ArenaEnemies(lv int32) (robots []string) {
	for _, v := range self.Index {
		if v.Lv >= lv-2 && v.Lv <= lv+2 {
			robots = append(robots, v.Id)
		}
	}

	return
}

func (self *robot_mgr_t) FindRobot(id string) *Robot {
	return self.Index[id]
}

func (self *robot_mgr_t) GetByLevel(lv int32, n int32) (ret []*Robot) {
	a := self.Robots[lv]
	if a == nil {
		return
	}

	L := int32(len(a))
	if n > L {
		n = L
	}

	return a[:n]
}
