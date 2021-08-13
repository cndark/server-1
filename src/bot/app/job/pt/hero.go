package pt

import (
	"fw/src/bot/app"
	"fw/src/bot/msg"
	"fw/src/core"
	"fw/src/core/evtmgr"
	"math/rand"
	"sync"
	"time"
)

// ============================================================================

type data_hero_t struct {
	heroes []*hero_t

	locker sync.Mutex
}

type hero_t struct {
	seq int64
	id  int32
}

// ============================================================================

func init() {
	evtmgr.On("userinfo", func(args ...interface{}) {
		bot := args[0].(*app.Bot)
		info := args[1].(*msg.GS_UserInfo)

		heroes := make([]*hero_t, 0, len(info.Bag.Heroes))
		for _, v := range info.Bag.Heroes {
			heroes = append(heroes, &hero_t{
				seq: v.Seq,
				id:  v.Id,
			})
		}

		bot.SetData("pt_hero", &data_hero_t{
			heroes: heroes,
		})

		bot.JobAdd("pt_hero_levelup", pt_hero_levelup)
		bot.JobAdd("pt_hero_starup", pt_hero_starup)
		bot.JobAdd("pt_hero_reset", pt_hero_reset)
		bot.JobAdd("pt_hero_decompose", pt_hero_decompose)
		bot.JobAdd("pt_hero_trinket", pt_hero_trinket)
	})

	evtmgr.On(app.MsgEvt(&msg.GS_BagUpdate{}), func(args ...interface{}) {
		bot := args[0].(*app.Bot)
		res := args[1].(*msg.GS_BagUpdate)

		d := bot.GetData("pt_hero").(*data_hero_t)

		d.locker.Lock()
		defer d.locker.Unlock()

		// add
		for _, v := range res.Heroes {
			d.heroes = append(d.heroes, &hero_t{
				seq: v.Seq,
				id:  v.Id,
			})
		}

		// del
		for _, a := range res.HeroesDel {
			for i, b := range d.heroes {
				if b.seq == a {
					L := len(d.heroes)
					d.heroes[i] = d.heroes[L-1]
					d.heroes = d.heroes[:L-1]
					break
				}
			}
		}
	})
}

func hero_rand_get(bot *app.Bot) *hero_t {
	d := bot.GetData("pt_hero").(*data_hero_t)

	d.locker.Lock()
	defer d.locker.Unlock()

	L := len(d.heroes)
	if L == 0 {
		return nil
	}

	return d.heroes[rand.Intn(L)]
}

func hero_rand_seqs(bot *app.Bot, n int) (ret []int64) {
	d := bot.GetData("pt_hero").(*data_hero_t)

	d.locker.Lock()
	defer d.locker.Unlock()

	L := len(d.heroes)
	if L == 0 {
		return
	}

	// shuffle
	rand.Shuffle(L, func(i, j int) {
		d.heroes[i], d.heroes[j] = d.heroes[j], d.heroes[i]
	})

	if n > L {
		n = L
	}

	ret = make([]int64, 0, n)
	for _, v := range d.heroes[:n] {
		ret = append(ret, v.seq)
	}

	return
}

func hero_get_count(bot *app.Bot) int {
	d := bot.GetData("pt_hero").(*data_hero_t)

	d.locker.Lock()
	defer d.locker.Unlock()

	L := len(d.heroes)
	return L
}

// ============================================================================

func pt_hero_levelup(bot *app.Bot) {
	hero := hero_rand_get(bot)
	if hero == nil {
		return
	}

	bot.SendMsg(&msg.C_HeroLevelUp{
		Seq: hero.seq,
		N:   core.RandInt32(1, 3),
	})
}

func pt_hero_starup(bot *app.Bot) {
	hero := hero_rand_get(bot)
	if hero == nil {
		return
	}

	bot.SendMsg(&msg.C_HeroStarUp{
		Seq:  hero.seq,
		Cost: nil,
	})
}

func pt_hero_reset(bot *app.Bot) {
	hero := hero_rand_get(bot)
	if hero == nil {
		return
	}

	bot.SendMsg(&msg.C_HeroReset{
		Seq: hero.seq,
	})
}

func pt_hero_decompose(bot *app.Bot) {
	cnt := hero_get_count(bot)
	if cnt <= 8 {
		return
	}

	hero := hero_rand_get(bot)
	if hero == nil {
		return
	}

	bot.SendMsg(&msg.C_HeroDecompose{
		Seqs: []int64{hero.seq},
	})
}

func pt_hero_trinket(bot *app.Bot) {
	hero := hero_rand_get(bot)
	if hero == nil {
		return
	}

	// do
	p := rand.Float32()

	if p < 0.2 { // unlock
		bot.SendMsg(&msg.C_HeroTrinketUnlock{
			Seq: hero.seq,
		})
	} else if p < 0.6 { // up
		b := false
		if rand.Float32() < 0.5 {
			b = true
		}

		bot.SendMsg(&msg.C_HeroTrinketUp{
			Seq:  hero.seq,
			Lock: b,
		})
	} else { // trans
		// gen
		for i := 0; i < 3; i++ {
			bot.SendMsg(&msg.C_HeroTrinketTransformGen{
				Seq: hero.seq,
			})

			time.Sleep(time.Millisecond * time.Duration(core.RandInt(500, 1000)))
		}

		// commit
		bot.SendMsg(&msg.C_HeroTrinketTransformCommit{
			Seq: hero.seq,
		})
	}
}
