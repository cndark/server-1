package utils

import (
	"fw/src/game/app/gamedata"
	Err "fw/src/proto/errorcode"
)

// ============================================================================

// 查询物品基础价格
func ConfBasePriceQuery(id int32, ccy int32) (int32, float64) {
	conf := gamedata.ConfBasePrice.Query(id)
	if conf != nil {
		for _, v := range conf.Price {
			if v.Id == ccy {
				return Err.OK, v.N
			}
		}
	}

	return Err.Failed, 0
}
