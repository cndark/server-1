package utils

import (
	"fw/src/core"
	"fw/src/game/app/dbmgr"
	"fw/src/game/app/gamedata"
	"math/rand"
)

// ============================================================================

const (
	code_book = "abcdefghijkmnpqrstuvwxyzABCDEFGHJKLMNPQRSTUVWXYZ1234567890"
)

// ============================================================================

func GenRandName(lang string, isSave ...bool) (name string) {
	conf := gamedata.ConfNameM.QueryItems(lang)
	if len(lang) == 0 {
		conf = gamedata.ConfNameM.QueryItems("cn")
	}

	L := len(conf)
	cnt := -1

	for {
		if L > 0 {
			name = conf[rand.Intn(L)].FirstName
			name += conf[rand.Intn(L)].LastName
		} else {
			name += GenRandCode(10)
		}

		cnt++
		if cnt < 5 {
		} else if cnt < 30 {
			name += GenRandCode(2)
		} else if cnt < 50 {
			name += GenRandCode(5)
		} else {
			name += GenRandCode(10)
		}

		if core.DefTrue(isSave) {
			if dbmgr.Center_InsertName(name) {
				return
			}
		} else {
			return
		}
	}

}

func GenRandCode(n int) string {
	L := len(code_book)

	b := make([]byte, n)
	for i := 0; i < n; i++ {
		b[i] = code_book[rand.Intn(L)]
	}
	return string(b)
}
