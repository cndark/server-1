package utils

import (
	"fw/src/core"
	Err "fw/src/proto/errorcode"
)

// ============================================================================

var istatustab IStatusTab

// ============================================================================

type IStatusTab interface {
	IsStatusTabConform(body interface{}, id int32, val float64) int32
}

// ============================================================================

// 用法： utils.IsStatusTabConform(body,id,val) == Err.OK
func IsStatusTabConform(body interface{}, id int32, val float64) int32 {
	return istatustab.IsStatusTabConform(body, id, val)
}

// 用法： utils.IsStatusTabArrayConform(body,confXXX.FieldName, or) == Err.OK
// 须知： confXXX.FieldName是策划表里面的字段,字段名可以自定义,但类型必须和下面的断言类型一致.否则异常
// 须知： or: true表示满足任意, false或不传表示都要满足
func IsStatusTabArrayConform(body interface{}, confXXX_FieldName interface{}, or ...bool) int32 {
	status := confXXX_FieldName.([]*struct {
		StatusId int32   `json:"statusId"`
		Val      float64 `json:"val"`
	})

	ec := int32(Err.OK)
	if core.DefFalse(or) { // 或
		for _, v := range status {
			ec = istatustab.IsStatusTabConform(body, v.StatusId, v.Val)
			if ec == Err.OK {
				break
			}
		}
	} else { // 且
		for _, v := range status {
			ec = istatustab.IsStatusTabConform(body, v.StatusId, v.Val)
			if ec != Err.OK {
				break
			}
		}
	}

	return ec
}
