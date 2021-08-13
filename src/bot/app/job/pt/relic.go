package pt

import (
	"fw/src/bot/app"
	"fw/src/bot/msg"
	"fw/src/core"
	"fw/src/core/evtmgr"
	"math/rand"
	"sync"
)

// ============================================================================

type data_relic_t struct {
	relics []*relic_t

	locker sync.Mutex
}

type relic_t struct {
	seq int64
	id  int32
}

// ============================================================================

func init() {
	evtmgr.On("userinfo", func(args ...interface{}) {
		bot := args[0].(*app.Bot)
		info := args[1].(*msg.GS_UserInfo)

		relics := make([]*relic_t, 0, len(info.Bag.Relics))
		for _, v := range info.Bag.Relics {
			relics = append(relics, &relic_t{
				seq: v.Seq,
				id:  v.Id,
			})
		}

		bot.SetData("pt_relic", &data_relic_t{
			relics: relics,
		})

		bot.JobAdd("pt_relic_equip", pt_relic_equip)
		bot.JobAdd("pt_relic_eat", pt_relic_eat)
	})

	evtmgr.On(app.MsgEvt(&msg.GS_BagUpdate{}), func(args ...interface{}) {
		bot := args[0].(*app.Bot)
		res := args[1].(*msg.GS_BagUpdate)

		d := bot.GetData("pt_relic").(*data_relic_t)

		d.locker.Lock()
		defer d.locker.Unlock()

		// add
		for _, v := range res.Relics {
			d.relics = append(d.relics, &relic_t{
				seq: v.Seq,
				id:  v.Id,
			})
		}

		// del
		for _, a := range res.RelicsDel {
			for i, b := range d.relics {
				if b.seq == a {
					L := len(d.relics)
					d.relics[i] = d.relics[L-1]
					d.relics = d.relics[:L-1]
					break
				}
			}
		}
	})
}

func relic_rand_get(bot *app.Bot) *relic_t {
	d := bot.GetData("pt_relic").(*data_relic_t)

	d.locker.Lock()
	defer d.locker.Unlock()

	L := len(d.relics)
	if L == 0 {
		return nil
	}

	return d.relics[rand.Intn(L)]
}

func relic_rand_seqs(bot *app.Bot, n int) (ret []int64) {
	d := bot.GetData("pt_relic").(*data_relic_t)

	d.locker.Lock()
	defer d.locker.Unlock()

	L := len(d.relics)
	if L == 0 {
		return
	}

	// shuffle
	rand.Shuffle(L, func(i, j int) {
		d.relics[i], d.relics[j] = d.relics[j], d.relics[i]
	})

	if n > L {
		n = L
	}

	ret = make([]int64, 0, n)
	for _, v := range d.relics[:n] {
		ret = append(ret, v.seq)
	}

	return
}

// ============================================================================

func pt_relic_equip(bot *app.Bot) {
	hero := hero_rand_get(bot)
	if hero == nil {
		return
	}

	if rand.Float32() < 0.5 {
		// equip
		rlc := relic_rand_get(bot)
		if rlc == nil {
			return
		}

		bot.SendMsg(&msg.C_RelicEquip{
			HeroSeq: hero.seq,
			Seq:     rlc.seq,
		})
	} else {
		// unequip
		bot.SendMsg(&msg.C_RelicUnequip{
			HeroSeq: hero.seq,
		})
	}
}

func pt_relic_eat(bot *app.Bot) {
	rlc := relic_rand_get(bot)
	if rlc == nil {
		return
	}

	bot.SendMsg(&msg.C_RelicEat{
		Seq:     rlc.seq,
		EatSeqs: relic_rand_seqs(bot, core.RandInt(1, 3)),
	})
}
